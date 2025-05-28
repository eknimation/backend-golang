package controllers

type UserDTO struct {
	Name     string `json:"name" validate:"required,min=2,max=50" example:"Jimmy"`
	Email    string `json:"email" validate:"required,emailFormat" example:"jimmy@example.com"`
	Password string `json:"password" validate:"required,isComplexPassword" example:"pas$worD123"`
}

type LoginDTO struct {
	Email    string `json:"email" validate:"required,email" example:"jimmy@example.com"`
	Password string `json:"password" validate:"required" example:"pas$worD123"`
}

type LoginResponseDTO struct {
	Token string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

type UserResponseDTO struct {
	ID        string `json:"id" example:"507f1f77bcf86cd799439011"`
	Name      string `json:"name" example:"Jimmy"`
	Email     string `json:"email" example:"jimmy@example.com"`
	CreatedAt string `json:"createdAt" example:"2023-01-01T00:00:00Z"`
}

type PaginationRequest struct {
	Page  int `json:"page" query:"page" validate:"min=1" example:"1"`
	Limit int `json:"limit" query:"limit" validate:"min=1,max=100" example:"10"`
}

type PaginationResponse struct {
	Page       int `json:"page" example:"1"`
	Limit      int `json:"limit" example:"10"`
	TotalItems int `json:"totalItems" example:"25"`
	TotalPages int `json:"totalPages" example:"3"`
}

type UserListResponseDTO struct {
	Users      []UserResponseDTO  `json:"users"`
	Pagination PaginationResponse `json:"pagination"`
}
