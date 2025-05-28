package domain

type UserInterface interface {
	CreateUser(user User) error
	GetUserByEmail(email string) (*User, error)
}
