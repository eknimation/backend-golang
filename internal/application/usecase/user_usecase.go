package usecase

import (
	"errors"
	"strings"
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
	storedUser, err := uc.userRepo.GetUserByEmail(user.Email)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	err = password.VerifyPassword(storedUser.Password, user.Password)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	appConfig := config.GetAppConfig()
	if appConfig.JWTSecret == "" {
		return "", errors.New("JWT secret not configured")
	}

	jwtManager := jwt.NewJWTManager(appConfig.JWTSecret, 24*time.Hour)

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

func (uc *Usecase) GetUsersWithPagination(page, limit int) ([]*domain.User, int, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	users, totalCount, err := uc.userRepo.GetUsersWithPagination(page, limit)
	if err != nil {
		return nil, 0, err
	}
	return users, totalCount, nil
}

func (uc *Usecase) UpdateUser(id string, user domain.User) error {
	existingUser, err := uc.userRepo.GetUserByID(id)
	if err != nil {
		return err
	}

	if user.Email != "" && user.Email != existingUser.Email {
		_, err := uc.userRepo.GetUserByEmail(user.Email)
		if err == nil {
			return errors.New("email already exists")
		}

		if !strings.Contains(err.Error(), "user not found") {
			return err
		}
	}

	err = uc.userRepo.UpdateUser(id, user)
	if err != nil {
		return err
	}

	return nil
}

func (uc *Usecase) DeleteUser(id string) error {
	_, err := uc.userRepo.GetUserByID(id)
	if err != nil {
		return err
	}

	err = uc.userRepo.DeleteUser(id)
	if err != nil {
		return err
	}

	return nil
}
