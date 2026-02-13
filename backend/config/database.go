package config

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoClient *mongo.Client

// ConnectDB connects to MongoDB
func ConnectDB(ctx context.Context) (*mongo.Client, error) {
	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}

	clientOptions := options.Client().ApplyURI(mongoURI)

	// Set connection timeout
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Println("Failed to create MongoDB client:", err)
		return nil, err
	}

	// Test the connection
	if err := client.Ping(ctx, nil); err != nil {
		log.Println("Failed to ping MongoDB:", err)
		return nil, err
	}

	MongoClient = client
	log.Println("Connected to MongoDB successfully")
	return client, nil
}

// DisconnectDB closes MongoDB connection
func DisconnectDB(ctx context.Context) error {
	if MongoClient == nil {
		return nil
	}
	return MongoClient.Disconnect(ctx)
}

// GetCollection returns a MongoDB collection
func GetCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "clothes_store"
	}
	return client.Database(dbName).Collection(collectionName)
}

// GetDB returns a MongoDB database
func GetDB(client *mongo.Client) *mongo.Database {
	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "clothes_store"
	}
	return client.Database(dbName)
}
