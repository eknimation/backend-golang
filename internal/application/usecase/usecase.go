package usecase

import (
	"backend-service/internal/domain"
	"backend-service/internal/infrastructure/repositories"
)

type Usecase struct {
	userRepo domain.UserInterface
}

func New(repo *repositories.Repository) *Usecase {
	return &Usecase{
		userRepo: repo.UserRepo,
	}
}
