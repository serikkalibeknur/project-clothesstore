package repository

import (
	"context"
	"time"

	"github.com/serikkalibeknur/project-clothesstore/config"
	"github.com/serikkalibeknur/project-clothesstore/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CartRepository struct {
	client *mongo.Client
}

func NewCartRepository(client *mongo.Client) *CartRepository {
	return &CartRepository{client: client}
}

func (cr *CartRepository) GetCart(ctx context.Context, userID primitive.ObjectID) (*models.Cart, error) {
	collection := config.GetCollection(cr.client, "carts")

	var cart models.Cart
	err := collection.FindOne(ctx, bson.M{"userID": userID}).Decode(&cart)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// Create new cart if it doesn't exist
			return &models.Cart{
				UserID:    userID,
				Products:  []models.CartItem{},
				UpdatedAt: time.Now(),
			}, nil
		}
		return nil, err
	}

	return &cart, nil
}

func (cr *CartRepository) AddToCart(ctx context.Context, userID, productID primitive.ObjectID, quantity int, size string) error {
	collection := config.GetCollection(cr.client, "carts")

	// Check if item already exists
	var cart models.Cart
	err := collection.FindOne(ctx, bson.M{"userID": userID}).Decode(&cart)

	if err != nil && err != mongo.ErrNoDocuments {
		return err
	}

	if err == mongo.ErrNoDocuments {
		// Create new cart
		newCart := models.Cart{
			ID:     primitive.NewObjectID(),
			UserID: userID,
			Products: []models.CartItem{
				{
					ProductID: productID,
					Quantity:  quantity,
					Size:      size,
					AddedAt:   time.Now(),
				},
			},
			UpdatedAt: time.Now(),
		}
		_, err := collection.InsertOne(ctx, newCart)
		return err
	}

	// Check if product already exists in cart
	existingIndex := -1
	for i, item := range cart.Products {
		if item.ProductID == productID && item.Size == size {
			existingIndex = i
			break
		}
	}

	if existingIndex != -1 {
		// Update quantity
		cart.Products[existingIndex].Quantity += quantity
	} else {
		// Add new item
		cart.Products = append(cart.Products, models.CartItem{
			ProductID: productID,
			Quantity:  quantity,
			Size:      size,
			AddedAt:   time.Now(),
		})
	}

	cart.UpdatedAt = time.Now()

	_, err = collection.UpdateOne(
		ctx,
		bson.M{"userID": userID},
		bson.M{"$set": bson.M{"products": cart.Products, "updatedAt": cart.UpdatedAt}},
	)

	return err
}

func (cr *CartRepository) RemoveFromCart(ctx context.Context, userID, productID primitive.ObjectID) error {
	collection := config.GetCollection(cr.client, "carts")

	_, err := collection.UpdateOne(
		ctx,
		bson.M{"userID": userID},
		bson.M{
			"$pull": bson.M{"products": bson.M{"productID": productID}},
			"$set":  bson.M{"updatedAt": time.Now()},
		},
	)

	return err
}

func (cr *CartRepository) UpdateCartItemQuantity(ctx context.Context, userID, productID primitive.ObjectID, quantity int) error {
	collection := config.GetCollection(cr.client, "carts")

	var cart models.Cart
	err := collection.FindOne(ctx, bson.M{"userID": userID}).Decode(&cart)
	if err != nil {
		return err
	}

	// Find and update item quantity
	for i, item := range cart.Products {
		if item.ProductID == productID {
			cart.Products[i].Quantity = quantity
			break
		}
	}

	cart.UpdatedAt = time.Now()

	_, err = collection.UpdateOne(
		ctx,
		bson.M{"userID": userID},
		bson.M{"$set": bson.M{"products": cart.Products, "updatedAt": cart.UpdatedAt}},
	)

	return err
}

func (cr *CartRepository) ClearCart(ctx context.Context, userID primitive.ObjectID) error {
	collection := config.GetCollection(cr.client, "carts")

	_, err := collection.UpdateOne(
		ctx,
		bson.M{"userID": userID},
		bson.M{"$set": bson.M{"products": []models.CartItem{}, "updatedAt": time.Now()}},
	)

	return err
}

func (cr *CartRepository) DeleteCart(ctx context.Context, userID primitive.ObjectID) error {
	collection := config.GetCollection(cr.client, "carts")

	_, err := collection.DeleteOne(ctx, bson.M{"userID": userID})
	return err
}
