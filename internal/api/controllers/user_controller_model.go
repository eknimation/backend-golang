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
