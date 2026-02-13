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

type WishlistRepository struct {
	client *mongo.Client
}

func NewWishlistRepository(client *mongo.Client) *WishlistRepository {
	return &WishlistRepository{client: client}
}

func (wr *WishlistRepository) GetWishlist(ctx context.Context, userID primitive.ObjectID) (*models.Wishlist, error) {
	collection := config.GetCollection(wr.client, "wishlists")

	var wishlist models.Wishlist
	err := collection.FindOne(ctx, bson.M{"userID": userID}).Decode(&wishlist)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// Create new wishlist if it doesn't exist
			return &models.Wishlist{
				UserID:    userID,
				Products:  []models.WishlistItem{},
				UpdatedAt: time.Now(),
			}, nil
		}
		return nil, err
	}

	return &wishlist, nil
}

func (wr *WishlistRepository) AddToWishlist(ctx context.Context, userID, productID primitive.ObjectID) error {
	collection := config.GetCollection(wr.client, "wishlists")

	// Check if wishlist exists
	var wishlist models.Wishlist
	err := collection.FindOne(ctx, bson.M{"userID": userID}).Decode(&wishlist)

	if err != nil && err != mongo.ErrNoDocuments {
		return err
	}

	if err == mongo.ErrNoDocuments {
		// Create new wishlist
		newWishlist := models.Wishlist{
			ID:     primitive.NewObjectID(),
			UserID: userID,
			Products: []models.WishlistItem{
				{
					ProductID: productID,
					AddedAt:   time.Now(),
				},
			},
			UpdatedAt: time.Now(),
		}
		_, err := collection.InsertOne(ctx, newWishlist)
		return err
	}

	// Check if product already exists in wishlist
	for _, item := range wishlist.Products {
		if item.ProductID == productID {
			// Product already in wishlist
			return nil
		}
	}

	// Add new item to wishlist
	wishlist.Products = append(wishlist.Products, models.WishlistItem{
		ProductID: productID,
		AddedAt:   time.Now(),
	})

	wishlist.UpdatedAt = time.Now()

	_, err = collection.UpdateOne(
		ctx,
		bson.M{"userID": userID},
		bson.M{"$set": bson.M{"products": wishlist.Products, "updatedAt": wishlist.UpdatedAt}},
	)

	return err
}

func (wr *WishlistRepository) RemoveFromWishlist(ctx context.Context, userID, productID primitive.ObjectID) error {
	collection := config.GetCollection(wr.client, "wishlists")

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

func (wr *WishlistRepository) ClearWishlist(ctx context.Context, userID primitive.ObjectID) error {
	collection := config.GetCollection(wr.client, "wishlists")

	_, err := collection.UpdateOne(
		ctx,
		bson.M{"userID": userID},
		bson.M{"$set": bson.M{"products": []models.WishlistItem{}, "updatedAt": time.Now()}},
	)

	return err
}

func (wr *WishlistRepository) DeleteWishlist(ctx context.Context, userID primitive.ObjectID) error {
	collection := config.GetCollection(wr.client, "wishlists")

	_, err := collection.DeleteOne(ctx, bson.M{"userID": userID})
	return err
}
