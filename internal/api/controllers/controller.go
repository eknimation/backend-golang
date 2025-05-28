package controllers

import (
	"backend-service/internal/application/usecase"
	"backend-service/pkg/utilities/middlewares"

	"github.com/labstack/echo/v4"
)

type Controller struct {
	uc *usecase.Usecase
}

func InitController(e *echo.Echo, usecase *usecase.Usecase) {
	controller := &Controller{uc: usecase}
	apiV1 := e.Group("/v1")

	// Public user routes (no auth required)
	apiV1.POST("/users", controller.CreateUser) // Create user account
	apiV1.POST("/users/login", controller.AuthenticateUser)

	// Protected user routes (JWT required)
	usersGroup := apiV1.Group("/users")
	usersGroup.Use(middlewares.JWTAuth())
	usersGroup.GET("/:id", controller.GetUserByID)
	// usersGroup.GET("/me", controller.GetCurrentUser)
}
