package main

import (
	"fmt"
	"log"

	"backend-service/config"
	"backend-service/docs"
	"backend-service/internal/api/controllers"
	"backend-service/internal/application/usecase"
	"backend-service/internal/infrastructure/database"
	"backend-service/internal/infrastructure/database/models"
	"backend-service/internal/infrastructure/repositories"
	"backend-service/pkg/utilities/logger"
	"backend-service/pkg/utilities/middlewares"

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
//
// @securityDefinitions.apikey X-API-Key
// @in header
// @name X-API-Key
// @description API key for authentication
func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	appConfig := config.GetAppConfig()
	logger := logger.GetLogger()

	docs.SwaggerInfo.BasePath = "/" + appConfig.BasePath

	e := echo.New()

	e.Use(middlewares.RequestResponseLogger(logger))

	dbClient := database.ConnectDB(appConfig.MongoDB)

	// Set up database indexes
	if err := models.SetupAllIndexes(dbClient, appConfig.MongoDB.DatabaseName); err != nil {
		log.Fatalf("Failed to setup database indexes: %v", err)
	}

	repo := repositories.New(dbClient, appConfig.MongoDB.DatabaseName)
	usecase := usecase.New(repo)

	if appConfig.Env != "production" {
		e.GET("/swagger/*", echoSwagger.WrapHandler)
	}

	controllers.InitController(e, usecase)

	// Log the routes
	for _, route := range e.Routes() {
		logger.Info(fmt.Sprintf("Method: %s, Path: %s", route.Method, route.Path))
	}

	e.Logger.Fatal(e.Start(":" + appConfig.Port))
}
