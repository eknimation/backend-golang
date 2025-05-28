package usecase

import (
	"errors"
	"time"

	"backend-service/config"
	"backend-service/internal/domain"
	"backend-service/pkg/utilities/jwt"
	"backend-service/pkg/utilities/password"
)

func (uc *Usecase) CreateUser(user domain.User) error {
	hashedPassword, err := password.HashPassword(user.Password)
	if err != nil {
		return err
	}

	user.Password = hashedPassword

	err = uc.userRepo.CreateUser(user)
	if err != nil {
		return err
	}
	return nil
}

func (uc *Usecase) AuthenticateUser(user domain.User) (string, error) {
	// Get user by email
	storedUser, err := uc.userRepo.GetUserByEmail(user.Email)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	// Verify password
	err = password.VerifyPassword(storedUser.Password, user.Password)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	// Get JWT configuration
	appConfig := config.GetAppConfig()
	if appConfig.JWTSecret == "" {
		return "", errors.New("JWT secret not configured")
	}

	// Create JWT manager with 24 hour token duration
	jwtManager := jwt.NewJWTManager(appConfig.JWTSecret, 24*time.Hour)

	// Generate JWT token
	token, err := jwtManager.GenerateToken(storedUser.ID, storedUser.Email)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (uc *Usecase) GetUserByID(id string) (*domain.User, error) {
	user, err := uc.userRepo.GetUserByID(id)
	if err != nil {
		return nil, err
	}
	return user, nil
}
