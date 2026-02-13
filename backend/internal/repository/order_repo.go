package repository

import (
	"context"

	"github.com/serikkalibeknur/project-clothesstore/config"
	"github.com/serikkalibeknur/project-clothesstore/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type OrderRepository struct {
	client *mongo.Client
}

func NewOrderRepository(client *mongo.Client) *OrderRepository {
	return &OrderRepository{client: client}
}

func (or *OrderRepository) CreateOrder(ctx context.Context, order *models.Order) (*models.Order, error) {
	collection := config.GetCollection(or.client, "orders")

	result, err := collection.InsertOne(ctx, order)
	if err != nil {
		return nil, err
	}

	order.ID = result.InsertedID.(primitive.ObjectID)
	return order, nil
}

func (or *OrderRepository) GetOrderByID(ctx context.Context, orderID primitive.ObjectID) (*models.Order, error) {
	collection := config.GetCollection(or.client, "orders")

	var order models.Order
	err := collection.FindOne(ctx, bson.M{"_id": orderID}).Decode(&order)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &order, nil
}

func (or *OrderRepository) GetUserOrders(ctx context.Context, userID primitive.ObjectID) ([]models.Order, error) {
	collection := config.GetCollection(or.client, "orders")

	cursor, err := collection.Find(ctx, bson.M{"userID": userID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var orders []models.Order
	if err = cursor.All(ctx, &orders); err != nil {
		return nil, err
	}

	return orders, nil
}

func (or *OrderRepository) GetAllOrders(ctx context.Context) ([]models.Order, error) {
	collection := config.GetCollection(or.client, "orders")

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var orders []models.Order
	if err = cursor.All(ctx, &orders); err != nil {
		return nil, err
	}

	return orders, nil
}

func (or *OrderRepository) UpdateOrderStatus(ctx context.Context, orderID primitive.ObjectID, status string) (*models.Order, error) {
	collection := config.GetCollection(or.client, "orders")

	result := collection.FindOneAndUpdate(
		ctx,
		bson.M{"_id": orderID},
		bson.M{"$set": bson.M{"status": status}},
	)

	if result.Err() != nil {
		return nil, result.Err()
	}

	var updatedOrder models.Order
	if err := result.Decode(&updatedOrder); err != nil {
		return nil, err
	}

	return &updatedOrder, nil
}

func (or *OrderRepository) GetOrdersByStatus(ctx context.Context, status string) ([]models.Order, error) {
	collection := config.GetCollection(or.client, "orders")

	cursor, err := collection.Find(ctx, bson.M{"status": status})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var orders []models.Order
	if err = cursor.All(ctx, &orders); err != nil {
		return nil, err
	}

	return orders, nil
}

func (or *OrderRepository) DeleteOrder(ctx context.Context, orderID primitive.ObjectID) error {
	collection := config.GetCollection(or.client, "orders")

	_, err := collection.DeleteOne(ctx, bson.M{"_id": orderID})
	return err
}

func (or *OrderRepository) GetTotalOrders(ctx context.Context) (int64, error) {
	collection := config.GetCollection(or.client, "orders")

	count, err := collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (or *OrderRepository) GetTotalRevenue(ctx context.Context) (float64, error) {
	collection := config.GetCollection(or.client, "orders")

	pipeline := []bson.M{
		{
			"$group": bson.M{
				"_id":          nil,
				"totalRevenue": bson.M{"$sum": "$totalPrice"},
			},
		},
	}

	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return 0, err
	}
	defer cursor.Close(ctx)

	var result []bson.M
	if err = cursor.All(ctx, &result); err != nil {
		return 0, err
	}

	if len(result) == 0 {
		return 0, nil
	}

	revenue := result[0]["totalRevenue"].(float64)
	return revenue, nil
}
