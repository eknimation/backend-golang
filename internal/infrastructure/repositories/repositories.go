package repositories

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository struct {
	UserRepo *Repo
}

type Repo struct {
	db *mongo.Database
}

func New(client *mongo.Client, databaseName string) *Repository {
	db := client.Database(databaseName)
	return &Repository{
		UserRepo: &Repo{
			db: db,
		},
	}
}
