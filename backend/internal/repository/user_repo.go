package repository

import (
	"context"

	"github.com/serikkalibeknur/project-clothesstore/config"
	"github.com/serikkalibeknur/project-clothesstore/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository struct {
	client *mongo.Client
}

func NewUserRepository(client *mongo.Client) *UserRepository {
	return &UserRepository{client: client}
}

func (ur *UserRepository) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	collection := config.GetCollection(ur.client, "users")

	result, err := collection.InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}

	user.ID = result.InsertedID.(primitive.ObjectID)
	return user, nil
}

func (ur *UserRepository) GetUserByID(ctx context.Context, userID primitive.ObjectID) (*models.User, error) {
	collection := config.GetCollection(ur.client, "users")

	var user models.User
	err := collection.FindOne(ctx, bson.M{"_id": userID}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (ur *UserRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	collection := config.GetCollection(ur.client, "users")

	var user models.User
	err := collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (ur *UserRepository) UpdateUser(ctx context.Context, userID primitive.ObjectID, user *models.User) (*models.User, error) {
	collection := config.GetCollection(ur.client, "users")

	updateData := bson.M{
		"name":      user.Name,
		"email":     user.Email,
		"role":      user.Role,
		"updatedAt": user.UpdatedAt,
	}

	if user.Password != "" {
		updateData["password"] = user.Password
	}

	result := collection.FindOneAndUpdate(
		ctx,
		bson.M{"_id": userID},
		bson.M{"$set": updateData},
	)

	if result.Err() != nil {
		return nil, result.Err()
	}

	var updatedUser models.User
	if err := result.Decode(&updatedUser); err != nil {
		return nil, err
	}

	return &updatedUser, nil
}

func (ur *UserRepository) DeleteUser(ctx context.Context, userID primitive.ObjectID) error {
	collection := config.GetCollection(ur.client, "users")

	_, err := collection.DeleteOne(ctx, bson.M{"_id": userID})
	return err
}

func (ur *UserRepository) GetAllUsers(ctx context.Context) ([]models.User, error) {
	collection := config.GetCollection(ur.client, "users")

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []models.User
	if err = cursor.All(ctx, &users); err != nil {
		return nil, err
	}

	return users, nil
}

func (ur *UserRepository) GetTotalUsers(ctx context.Context) (int64, error) {
	collection := config.GetCollection(ur.client, "users")

	count, err := collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return 0, err
	}

	return count, nil
}
