package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CartItem struct {
	ProductID primitive.ObjectID `json:"productID" bson:"productID"`
	Quantity  int                `json:"quantity" bson:"quantity"`
	Size      string             `json:"size" bson:"size"`
	AddedAt   time.Time          `json:"addedAt" bson:"addedAt"`
}

type Cart struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID    primitive.ObjectID `json:"userID" bson:"userID"`
	Products  []CartItem         `json:"products" bson:"products"`
	UpdatedAt time.Time          `json:"updatedAt" bson:"updatedAt"`
}

type AddToCartInput struct {
	ProductID string `json:"productID"`
	Quantity  int    `json:"quantity"`
	Size      string `json:"size"`
}
