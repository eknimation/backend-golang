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
// @Router       /v1/users [post]
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

// @Summary      Authenticate User
// @Description  Authenticate a user with email and password, returns JWT token
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        credentials  body      LoginDTO  true  "User login credentials"
// @Success      200      {object}  responses.Response{data=LoginResponseDTO}
// @Failure      400      {object}  responses.ErrorResponse  "Bad Request - Validation failed"
// @Failure      401      {object}  responses.ErrorResponse  "Unauthorized - Invalid credentials"
// @Failure      500      {object}  responses.ErrorResponse
// @Router       /v1/users/login [post]
func (ctrl *Controller) AuthenticateUser(c echo.Context) error {
	var loginDTO LoginDTO
	if err := c.Bind(&loginDTO); err != nil {
		return c.JSON(http.StatusBadRequest, responses.Error(http.StatusBadRequest, err.Error()))
	}

	// Validate the loginDTO struct
	if err := validator.Validate(&loginDTO); err != nil {
		validationErr := validator.FormatValidationErrors(err)
		return c.JSON(http.StatusBadRequest, validationErr)
	}

	user := domain.User{
		Email:    loginDTO.Email,
		Password: loginDTO.Password,
	}

	// Authenticate user and get JWT token
	token, err := ctrl.uc.AuthenticateUser(user)
	if err != nil {
		if strings.Contains(err.Error(), "invalid credentials") {
			return c.JSON(http.StatusUnauthorized, responses.Error(http.StatusUnauthorized, "Invalid credentials"))
		}
		return c.JSON(http.StatusInternalServerError, responses.Error(http.StatusInternalServerError, err.Error()))
	}

	// Create response
	responseData := LoginResponseDTO{
		Token: token,
	}

	return c.JSON(http.StatusOK, responses.Ok(http.StatusOK, "Authentication successful", responseData))
}

// @Summary      Get User by ID
// @Description  Retrieve a user's details by their ID
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "User ID"
// @Success      200  {object}  responses.Response{data=UserResponseDTO}
// @Failure      400  {object}  responses.ErrorResponse  "Bad Request - Invalid user ID format"
// @Failure      404  {object}  responses.ErrorResponse  "Not Found - User not found"
// @Failure      500  {object}  responses.ErrorResponse
// @Router       /v1/users/{id} [get]
// @Security     BearerAuth
func (ctrl *Controller) GetUserByID(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, responses.Error(http.StatusBadRequest, "User ID is required"))
	}

	user, err := ctrl.uc.GetUserByID(id)
	if err != nil {
		if strings.Contains(err.Error(), "invalid user ID format") {
			return c.JSON(http.StatusBadRequest, responses.Error(http.StatusBadRequest, "Invalid user ID format"))
		}
		if strings.Contains(err.Error(), "user not found") {
			return c.JSON(http.StatusNotFound, responses.Error(http.StatusNotFound, "User not found"))
		}
		return c.JSON(http.StatusInternalServerError, responses.Error(http.StatusInternalServerError, err.Error()))
	}

	// Create response DTO
	responseData := UserResponseDTO{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	return c.JSON(http.StatusOK, responses.Ok(http.StatusOK, "User retrieved successfully", responseData))
}
