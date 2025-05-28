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

// @Summary      Get All Users
// @Description  Retrieve all users in the system with pagination
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        page   query     int  false  "Page number (default: 1)"   minimum(1)
// @Param        limit  query     int  false  "Items per page (default: 10, max: 100)"  minimum(1)  maximum(100)
// @Success      200  {object}  responses.Response{data=UserListResponseDTO}
// @Failure      400  {object}  responses.ErrorResponse  "Bad Request - Invalid pagination parameters"
// @Failure      500  {object}  responses.ErrorResponse
// @Router       /v1/users [get]
// @Security     BearerAuth
func (ctrl *Controller) GetAllUsers(c echo.Context) error {
	// Parse and validate pagination parameters
	pagination, err := ctrl.parsePaginationParams(c)
	if err != nil {
		return err
	}

	// Get users with pagination from use case
	users, totalCount, err := ctrl.uc.GetUsersWithPagination(pagination.Page, pagination.Limit)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.Error(http.StatusInternalServerError, err.Error()))
	}

	// Build response
	responseData := ctrl.buildUserListResponse(users, pagination, totalCount)

	return c.JSON(http.StatusOK, responses.Ok(http.StatusOK, "Users retrieved successfully", responseData))
}

// parsePaginationParams extracts and validates pagination parameters from request
func (ctrl *Controller) parsePaginationParams(c echo.Context) (PaginationRequest, error) {
	pagination := PaginationRequest{
		Page:  1,  // default
		Limit: 10, // default
	}

	if err := c.Bind(&pagination); err != nil {
		return pagination, c.JSON(http.StatusBadRequest, responses.Error(http.StatusBadRequest, "Invalid pagination parameters"))
	}

	if err := validator.Validate(&pagination); err != nil {
		validationErr := validator.FormatValidationErrors(err)
		return pagination, c.JSON(http.StatusBadRequest, validationErr)
	}

	return pagination, nil
}

// buildUserListResponse creates the paginated user list response
func (ctrl *Controller) buildUserListResponse(users []*domain.User, pagination PaginationRequest, totalCount int) UserListResponseDTO {
	userDTOs := ctrl.convertUsersToResponseDTOs(users)

	totalPages := (totalCount + pagination.Limit - 1) / pagination.Limit

	return UserListResponseDTO{
		Users: userDTOs,
		Pagination: PaginationResponse{
			Page:       pagination.Page,
			Limit:      pagination.Limit,
			TotalItems: totalCount,
			TotalPages: totalPages,
		},
	}
}

// convertUsersToResponseDTOs converts domain users to response DTOs
func (ctrl *Controller) convertUsersToResponseDTOs(users []*domain.User) []UserResponseDTO {
	userDTOs := make([]UserResponseDTO, 0, len(users))

	for _, user := range users {
		userDTO := UserResponseDTO{
			ID:        user.ID,
			Name:      user.Name,
			Email:     user.Email,
			CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
		userDTOs = append(userDTOs, userDTO)
	}

	return userDTOs
}

// @Summary      Update User
// @Description  Update user information (name and/or email)
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        id   path      string        true  "User ID"
// @Param        user body      UpdateUserDTO true  "User update details"
// @Success      200  {object}  responses.Response
// @Failure      400  {object}  responses.ErrorResponse  "Bad Request - Validation failed or email already exists"
// @Failure      404  {object}  responses.ErrorResponse  "User not found"
// @Failure      500  {object}  responses.ErrorResponse
// @Router       /v1/users/{id} [put]
// @Security     BearerAuth
func (ctrl *Controller) UpdateUser(c echo.Context) error {
	// Get user ID from path parameter
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, responses.Error(http.StatusBadRequest, "User ID is required"))
	}

	var updateUserDTO UpdateUserDTO
	if err := c.Bind(&updateUserDTO); err != nil {
		return c.JSON(http.StatusBadRequest, responses.Error(http.StatusBadRequest, err.Error()))
	}

	// Validate the updateUserDTO struct
	if err := validator.Validate(&updateUserDTO); err != nil {
		validationErr := validator.FormatValidationErrors(err)
		return c.JSON(http.StatusBadRequest, validationErr)
	}

	// Check if at least one field is provided for update
	if updateUserDTO.Name == "" && updateUserDTO.Email == "" {
		return c.JSON(http.StatusBadRequest, responses.Error(http.StatusBadRequest, "At least one field (name or email) must be provided for update"))
	}

	user := domain.User{
		Name:  updateUserDTO.Name,
		Email: updateUserDTO.Email,
	}

	err := ctrl.uc.UpdateUser(id, user)
	if err != nil {
		if strings.Contains(err.Error(), "user not found") {
			return c.JSON(http.StatusNotFound, responses.Error(http.StatusNotFound, "User not found"))
		}
		if strings.Contains(err.Error(), "email already exists") {
			return c.JSON(http.StatusBadRequest, responses.Error(http.StatusBadRequest, "Email already exists"))
		}
		if strings.Contains(err.Error(), "invalid user ID format") {
			return c.JSON(http.StatusBadRequest, responses.Error(http.StatusBadRequest, "Invalid user ID format"))
		}
		return c.JSON(http.StatusInternalServerError, responses.Error(http.StatusInternalServerError, err.Error()))
	}

	return c.JSON(http.StatusOK, responses.Ok(http.StatusOK, "User updated successfully", nil))
}
