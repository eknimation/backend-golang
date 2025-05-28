package services

import (
	"context"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
)

// PeriodicTaskService represents a service that runs periodic tasks
type PeriodicTaskService struct {
	name     string
	logger   *logrus.Logger
	ticker   *time.Ticker
	done     chan struct{}
	taskFunc func()
}

// NewPeriodicTaskService creates a new periodic task service
func NewPeriodicTaskService(name string, interval time.Duration, logger *logrus.Logger, taskFunc func()) *PeriodicTaskService {
	return &PeriodicTaskService{
		name:     name,
		logger:   logger,
		ticker:   time.NewTicker(interval),
		done:     make(chan struct{}),
		taskFunc: taskFunc,
	}
}

// Start starts the periodic task service
func (pts *PeriodicTaskService) Start(ctx context.Context) {
	pts.logger.Info(fmt.Sprintf("Starting periodic task service: %s", pts.name))

	go func() {
		for {
			select {
			case <-pts.done:
				pts.logger.Info(fmt.Sprintf("Stopping periodic task service: %s", pts.name))
				return
			case <-ctx.Done():
				pts.logger.Info(fmt.Sprintf("Context cancelled, stopping periodic task service: %s", pts.name))
				return
			case <-pts.ticker.C:
				if pts.taskFunc != nil {
					pts.taskFunc()
				}
			}
		}
	}()
}

// Shutdown gracefully shuts down the periodic task service
func (pts *PeriodicTaskService) Shutdown(ctx context.Context) error {
	pts.logger.Info(fmt.Sprintf("Shutting down periodic task service: %s", pts.name))

	// Stop the ticker
	pts.ticker.Stop()

	// Signal the goroutine to stop
	select {
	case pts.done <- struct{}{}:
	case <-ctx.Done():
		return ctx.Err()
	default:
		// Channel might be closed already
	}

	pts.logger.Info(fmt.Sprintf("Periodic task service %s shutdown completed", pts.name))
	return nil
}

// Name returns the name of the service
func (pts *PeriodicTaskService) Name() string {
	return pts.name
}
