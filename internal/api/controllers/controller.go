package controllers

import (
	"backend-service/internal/application/usecase"

	"github.com/labstack/echo/v4"
)

type Controller struct {
	uc *usecase.Usecase
}

func InitController(e *echo.Echo, usecase *usecase.Usecase) {
	controller := &Controller{uc: usecase}
	group := e.Group("/v1")
	group.POST("/user", controller.CreateUser)
}
