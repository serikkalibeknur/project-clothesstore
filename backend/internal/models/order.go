package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OrderItem struct {
	ProductID primitive.ObjectID `json:"productID" bson:"productID"`
	Name      string             `json:"name" bson:"name"`
	Price     float64            `json:"price" bson:"price"`
	Quantity  int                `json:"quantity" bson:"quantity"`
	Size      string             `json:"size" bson:"size"`
}

type Address struct {
	Street     string `json:"street" bson:"street"`
	City       string `json:"city" bson:"city"`
	PostalCode string `json:"postalCode" bson:"postalCode"`
	Country    string `json:"country" bson:"country"`
}

type Order struct {
	ID              primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID          primitive.ObjectID `json:"userID" bson:"userID"`
	Products        []OrderItem        `json:"products" bson:"products"`
	TotalPrice      float64            `json:"totalPrice" bson:"totalPrice"`
	Status          string             `json:"status" bson:"status"` // pending, processing, shipped, delivered, cancelled
	ShippingAddress Address            `json:"shippingAddress" bson:"shippingAddress"`
	CreatedAt       time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt       time.Time          `json:"updatedAt" bson:"updatedAt"`
}

type CreateOrderInput struct {
	ShippingAddress Address `json:"shippingAddress"`
}
