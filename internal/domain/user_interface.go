package domain

type UserInterface interface {
	CreateUser(user User) error
	GetUserByEmail(email string) (*User, error)
	GetUserByID(id string) (*User, error)
	GetUsersWithPagination(page, limit int) ([]*User, int, error)
	UpdateUser(id string, user User) error
	DeleteUser(id string) error
}
