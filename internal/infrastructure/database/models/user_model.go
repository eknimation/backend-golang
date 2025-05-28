package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserModel struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Name      string             `bson:"name"`
	Email     string             `bson:"email"`
	Password  string             `bson:"password"`
	CreatedAt time.Time          `bson:"createdAt"`
}
