package shutdown

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

// Service represents a service that can be shutdown gracefully
type Service interface {
	Shutdown(ctx context.Context) error
	Name() string
}

// GracefulShutdown handles the graceful shutdown of the application
type GracefulShutdown struct {
	logger   *logrus.Logger
	server   *echo.Echo
	dbClient *mongo.Client
	services []Service
	timeout  time.Duration
	mutex    sync.RWMutex
}

// NewGracefulShutdown creates a new graceful shutdown handler
func NewGracefulShutdown(logger *logrus.Logger, server *echo.Echo, dbClient *mongo.Client, timeout time.Duration) *GracefulShutdown {
	return &GracefulShutdown{
		logger:   logger,
		server:   server,
		dbClient: dbClient,
		services: make([]Service, 0),
		timeout:  timeout,
	}
}

// AddService adds a service to be shutdown gracefully
func (gs *GracefulShutdown) AddService(service Service) {
	gs.mutex.Lock()
	defer gs.mutex.Unlock()
	gs.services = append(gs.services, service)
}

// WaitForShutdown waits for shutdown signals and performs graceful shutdown
func (gs *GracefulShutdown) WaitForShutdown(ctx context.Context, cancel context.CancelFunc) {
	// Set up signal handling for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)

	// Block until we receive a signal
	<-quit
	gs.logger.Info("Received shutdown signal, starting graceful shutdown...")

	// Cancel the context to stop all goroutines
	cancel()

	// Perform shutdown
	if err := gs.Shutdown(context.Background()); err != nil {
		gs.logger.Error(fmt.Sprintf("Error during shutdown: %v", err))
		os.Exit(1)
	}
}

// SetupSignalHandler sets up signal handling without blocking
func (gs *GracefulShutdown) SetupSignalHandler() <-chan os.Signal {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	return quit
}

// Shutdown performs the shutdown operations with the given context
func (gs *GracefulShutdown) Shutdown(ctx context.Context) error {
	gs.logger.Info("Starting graceful shutdown...")

	// Create a context with timeout for shutdown operations
	shutdownCtx, shutdownCancel := context.WithTimeout(ctx, gs.timeout)
	defer shutdownCancel()

	// Create a WaitGroup to wait for all shutdown operations
	var wg sync.WaitGroup
	errChan := make(chan error, 2+len(gs.services))

	// Shutdown additional services first
	gs.mutex.RLock()
	for _, service := range gs.services {
		wg.Add(1)
		go func(svc Service) {
			defer wg.Done()
			gs.logger.Info(fmt.Sprintf("Shutting down service: %s", svc.Name()))
			if err := svc.Shutdown(shutdownCtx); err != nil {
				gs.logger.Error(fmt.Sprintf("Failed to shutdown service %s: %v", svc.Name(), err))
				errChan <- err
			} else {
				gs.logger.Info(fmt.Sprintf("Service %s shutdown completed", svc.Name()))
			}
		}(service)
	}
	gs.mutex.RUnlock()

	// Shutdown the Echo server
	wg.Add(1)
	go func() {
		defer wg.Done()
		gs.logger.Info("Shutting down HTTP server...")
		if err := gs.server.Shutdown(shutdownCtx); err != nil {
			gs.logger.Error(fmt.Sprintf("Failed to gracefully shutdown HTTP server: %v", err))
			errChan <- err
		} else {
			gs.logger.Info("HTTP server shutdown completed")
		}
	}()

	// Disconnect from MongoDB
	wg.Add(1)
	go func() {
		defer wg.Done()
		gs.logger.Info("Disconnecting from MongoDB...")
		if err := gs.dbClient.Disconnect(shutdownCtx); err != nil {
			gs.logger.Error(fmt.Sprintf("Failed to disconnect from MongoDB: %v", err))
			errChan <- err
		} else {
			gs.logger.Info("MongoDB disconnection completed")
		}
	}()

	// Wait for all shutdown operations to complete
	wg.Wait()
	close(errChan)

	// Check for any errors
	var shutdownErrors []error
	for err := range errChan {
		shutdownErrors = append(shutdownErrors, err)
	}

	if len(shutdownErrors) > 0 {
		gs.logger.Error(fmt.Sprintf("Shutdown completed with %d errors", len(shutdownErrors)))
		return shutdownErrors[0] // Return the first error
	}

	gs.logger.Info("Graceful shutdown completed successfully")
	return nil
}
