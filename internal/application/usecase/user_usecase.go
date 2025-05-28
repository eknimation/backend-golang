package usecase

import (
	"backend-service/internal/domain"
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
