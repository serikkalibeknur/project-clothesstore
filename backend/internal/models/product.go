package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Product struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name        string             `json:"name" bson:"name"`
	Description string             `json:"description" bson:"description"`
	Price       float64            `json:"price" bson:"price"`
	Size        []string           `json:"size" bson:"size"`
	Category    string             `json:"category" bson:"category"`
	ImageURL    string             `json:"imageURL" bson:"imageURL"`
	Stock       int                `json:"stock" bson:"stock"`
	CreatedAt   time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt   time.Time          `json:"updatedAt" bson:"updatedAt"`
}

type ProductInput struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Price       float64  `json:"price"`
	Size        []string `json:"size"`
	Category    string   `json:"category"`
	ImageURL    string   `json:"imageURL"`
	Stock       int      `json:"stock"`
}
