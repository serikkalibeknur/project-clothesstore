package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type WishlistItem struct {
	ProductID primitive.ObjectID `json:"productID" bson:"productID"`
	AddedAt   time.Time          `json:"addedAt" bson:"addedAt"`
}

type Wishlist struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID    primitive.ObjectID `json:"userID" bson:"userID"`
	Products  []WishlistItem     `json:"products" bson:"products"`
	UpdatedAt time.Time          `json:"updatedAt" bson:"updatedAt"`
}
