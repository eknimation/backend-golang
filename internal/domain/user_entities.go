package domain

import (
	"time"
)

type User struct {
	CreatedAt *time.Time `json:"createdAt"`
	Name      string     `json:"name"`
	Email     string     `json:"email"`
	Password  string     `json:"-"` // Exclude password from JSON responses
}
