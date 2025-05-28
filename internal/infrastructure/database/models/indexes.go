package models

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// SetupAllIndexes creates all necessary database indexes for all models
func SetupAllIndexes(client *mongo.Client, databaseName string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	db := client.Database(databaseName)

	// Setup user model indexes
	if err := setupUserIndexes(ctx, db); err != nil {
		return fmt.Errorf("failed to setup user indexes: %v", err)
	}

	// Add other model index setups here as needed
	// if err := setupOtherModelIndexes(ctx, db); err != nil {
	//     return fmt.Errorf("failed to setup other model indexes: %v", err)
	// }

	return nil
}

// setupUserIndexes creates indexes specific to the UserModel
func setupUserIndexes(ctx context.Context, db *mongo.Database) error {
	usersCollection := db.Collection("users")

	// Create unique index on email field
	emailIndexModel := mongo.IndexModel{
		Keys: bson.D{
			{Key: "email", Value: 1}, // 1 for ascending order
		},
		Options: options.Index().SetUnique(true),
	}

	indexName, err := usersCollection.Indexes().CreateOne(ctx, emailIndexModel)
	if err != nil {
		return fmt.Errorf("failed to create unique email index: %v", err)
	}

	fmt.Printf("Created unique index on email field: %s\n", indexName)

	// Add more user-specific indexes here if needed
	// Example: Index on createdAt for sorting
	// createdAtIndexModel := mongo.IndexModel{
	//     Keys: bson.D{{Key: "createdAt", Value: -1}}, // -1 for descending order
	// }
	// _, err = usersCollection.Indexes().CreateOne(ctx, createdAtIndexModel)
	// if err != nil {
	//     return fmt.Errorf("failed to create createdAt index: %v", err)
	// }

	return nil
}
