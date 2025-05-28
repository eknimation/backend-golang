package domain

type UserInterface interface {
	CreateUser(user User) error
}
