package repositories

import (
	"context"
	"errors"
	"time"

	"backend-service/internal/domain"
	"backend-service/internal/infrastructure/database/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (r *Repo) CreateUser(user domain.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	dbModel := models.UserModel{
		Name:      user.Name,
		Email:     user.Email,
		Password:  user.Password,
		CreatedAt: time.Now(),
	}

	collection := r.db.Collection("users")

	_, err := collection.InsertOne(ctx, dbModel)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return errors.New("email already exists")
		}
		return err
	}

	return nil
}

func (r *Repo) GetUserByEmail(email string) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := r.db.Collection("users")

	var userModel models.UserModel
	err := collection.FindOne(ctx, bson.M{"email": email}).Decode(&userModel)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	user := r.convertUserModelToEntityWithPassword(userModel)

	return user, nil
}

func (r *Repo) GetUserByID(id string) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid user ID format")
	}

	collection := r.db.Collection("users")

	var userModel models.UserModel
	err = collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&userModel)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	user := r.convertUserModelToEntity(userModel)

	return user, nil
}

func (r *Repo) GetUsersWithPagination(page, limit int) ([]*domain.User, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := r.db.Collection("users")
	filter := bson.M{}

	totalCount, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	skip := (page - 1) * limit
	skipValue := int64(skip)
	limitValue := int64(limit)

	findOptions := &options.FindOptions{
		Skip:  &skipValue,
		Limit: &limitValue,
		Sort:  bson.D{{Key: "createdAt", Value: -1}},
	}

	cursor, err := collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var userModels []models.UserModel
	if err = cursor.All(ctx, &userModels); err != nil {
		return nil, 0, err
	}

	users := r.convertUserModelsToEntities(userModels)

	return users, int(totalCount), nil
}

func (r *Repo) UpdateUser(id string, user domain.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid user ID format")
	}

	collection := r.db.Collection("users")

	updateDoc := bson.M{}
	if user.Name != "" {
		updateDoc["name"] = user.Name
	}
	if user.Email != "" {
		updateDoc["email"] = user.Email
	}

	if len(updateDoc) == 0 {
		return errors.New("no fields to update")
	}

	update := bson.M{"$set": updateDoc}

	result, err := collection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return errors.New("email already exists")
		}
		return err
	}

	if result.MatchedCount == 0 {
		return errors.New("user not found")
	}

	return nil
}

func (r *Repo) DeleteUser(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid user ID format")
	}

	collection := r.db.Collection("users")

	result, err := collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return errors.New("user not found")
	}

	return nil
}

func (r *Repo) GetUserCount() (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := r.db.Collection("users")
	filter := bson.M{}

	count, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *Repo) convertUserModelsToEntities(userModels []models.UserModel) []*domain.User {
	users := make([]*domain.User, 0, len(userModels))

	for _, userModel := range userModels {
		user := r.convertUserModelToEntity(userModel)
		users = append(users, user)
	}

	return users
}

func (r *Repo) convertUserModelToEntity(userModel models.UserModel) *domain.User {
	return &domain.User{
		ID:        userModel.ID.Hex(),
		Name:      userModel.Name,
		Email:     userModel.Email,
		CreatedAt: &userModel.CreatedAt,
	}
}

func (r *Repo) convertUserModelToEntityWithPassword(userModel models.UserModel) *domain.User {
	user := r.convertUserModelToEntity(userModel)
	user.Password = userModel.Password
	return user
}
