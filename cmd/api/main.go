package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"backend-service/config"
	"backend-service/docs"
	"backend-service/internal/api/controllers"
	"backend-service/internal/application/usecase"
	"backend-service/internal/infrastructure/database"
	"backend-service/internal/infrastructure/database/models"
	"backend-service/internal/infrastructure/repositories"
	"backend-service/pkg/utilities/logger"
	"backend-service/pkg/utilities/middlewares"
	"backend-service/pkg/utilities/services"
	"backend-service/pkg/utilities/shutdown"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// @contact.name			Akekapon S.
// @contact.email			eknimation@gmail.com
//
// @title					Backend Go API
// @version					v1.0.0
// @description				Backend service API for user management and authentication
// @host					localhost:5555
// @BasePath				/

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Bearer token for JWT authentication. Format: "Bearer {token}"
func main() {
	// Load .env file if it exists (optional for Docker environments)
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found or could not be loaded: %v", err)
		log.Println("Continuing with environment variables from system/docker-compose...")
	}

	appConfig := config.GetAppConfig()
	logger := logger.GetLogger()

	docs.SwaggerInfo.BasePath = "/" + appConfig.BasePath

	e := echo.New()

	e.Use(middlewares.RequestResponseLogger(logger))
	e.Use(middlewares.RequestContext(30 * time.Second)) // 30 second timeout for requests

	dbClient := database.ConnectDB(appConfig.MongoDB)

	// Set up database indexes
	if err := models.SetupAllIndexes(dbClient, appConfig.MongoDB.DatabaseName); err != nil {
		log.Fatalf("Failed to setup database indexes: %v", err)
	}

	repo := repositories.New(dbClient, appConfig.MongoDB.DatabaseName)
	usecase := usecase.New(repo)

	// Create a context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create periodic task service for user count monitoring
	userCountTask := services.NewPeriodicTaskService(
		"user-count-monitor",
		10*time.Second,
		logger,
		func() {
			count, err := usecase.GetUserCount()
			if err != nil {
				logger.Error(fmt.Sprintf("Failed to get user count: %v", err))
			} else {
				logger.Info(fmt.Sprintf("Total users in database: %d", count))
			}
		},
	)

	// Start the periodic task service
	userCountTask.Start(ctx)

	if appConfig.Env != "production" {
		e.GET("/swagger/*", echoSwagger.WrapHandler)
	}

	controllers.InitController(e, usecase)

	// Log the routes
	for _, route := range e.Routes() {
		logger.Info(fmt.Sprintf("Method: %s, Path: %s", route.Method, route.Path))
	}

	// Create graceful shutdown handler
	gracefulShutdown := shutdown.NewGracefulShutdown(logger, e, dbClient, 30*time.Second)

	// Add the periodic task service to graceful shutdown
	gracefulShutdown.AddService(userCountTask)

	// Start server in a goroutine
	go func() {
		logger.Info(fmt.Sprintf("Starting server on port %s", appConfig.Port))
		if err := e.Start(":" + appConfig.Port); err != nil && err != http.ErrServerClosed {
			logger.Error(fmt.Sprintf("Failed to start server: %v", err))
		}
	}()

	// Wait for shutdown signal and perform graceful shutdown
	gracefulShutdown.WaitForShutdown(ctx, cancel)
}
