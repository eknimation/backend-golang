package controllers

import (
	"net/http"
	"strings"

	"backend-service/internal/domain"
	"backend-service/pkg/utilities/responses"
	"backend-service/pkg/utilities/validator"

	"github.com/labstack/echo/v4"
)

// @Summary      Create User
// @Description  Create a new user account with name, email and password
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        user  body      UserDTO  true  "User registration details"
// @Success      201      {object}  responses.Response
// @Failure      400      {object}  responses.ErrorResponse  "Bad Request - Validation failed or email already exists"
// @Failure      500      {object}  responses.ErrorResponse
// @Router       /v1/user [post]
// @Security     X-API-Key
func (ctrl *Controller) CreateUser(c echo.Context) error {
	var userDTO UserDTO
	if err := c.Bind(&userDTO); err != nil {
		return c.JSON(http.StatusBadRequest, responses.Error(http.StatusBadRequest, err.Error()))
	}

	// Validate the userDTO struct
	if err := validator.Validate(&userDTO); err != nil {
		validationErr := validator.FormatValidationErrors(err)
		return c.JSON(http.StatusBadRequest, validationErr)
	}

	user := domain.User{
		Name:     userDTO.Name,
		Email:    userDTO.Email,
		Password: userDTO.Password,
	}

	err := ctrl.uc.CreateUser(user)
	if err != nil {
		if strings.Contains(err.Error(), "email already exists") {
			return c.JSON(http.StatusBadRequest, responses.Error(http.StatusBadRequest, "Email already exists"))
		}

		return c.JSON(http.StatusInternalServerError, responses.Error(http.StatusInternalServerError, err.Error()))
	}

	return c.JSON(http.StatusCreated, responses.Ok(http.StatusCreated, "User created successfully", nil))
}
