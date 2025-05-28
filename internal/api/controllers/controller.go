package controllers

import (
	"net/http"

	"backend-service/internal/application/usecase"
	"backend-service/pkg/utilities/middlewares"
	"backend-service/pkg/utilities/responses"

	"github.com/labstack/echo/v4"
)

type Controller struct {
	uc *usecase.Usecase
}

// @Summary      Health Check
// @Description  Check if the API is running and healthy
// @Tags         Health
// @Accept       json
// @Produce      json
// @Success      200  {object}  responses.Response
// @Router       /health [get]
func (ctrl *Controller) HealthCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, responses.Ok(http.StatusOK, "API is healthy", map[string]string{
		"status":  "ok",
		"service": "backend-golang-api",
	}))
}

func InitController(e *echo.Echo, usecase *usecase.Usecase) {
	controller := &Controller{uc: usecase}

	// Health check endpoint
	e.GET("/health", controller.HealthCheck)

	apiV1 := e.Group("/v1")

	// Public user routes (no auth required)
	apiV1.POST("/users", controller.CreateUser)
	apiV1.POST("/users/login", controller.AuthenticateUser)

	// Protected user routes (JWT required)
	usersGroup := apiV1.Group("/users")
	usersGroup.Use(middlewares.JWTAuth())
	usersGroup.GET("", controller.GetAllUsers)
	usersGroup.GET("/:id", controller.GetUserByID)
	usersGroup.PUT("/:id", controller.UpdateUser)
	usersGroup.DELETE("/:id", controller.DeleteUser)
}
