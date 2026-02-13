package repository

import (
	"context"

	"github.com/serikkalibeknur/project-clothesstore/config"
	"github.com/serikkalibeknur/project-clothesstore/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ProductRepository struct {
	client *mongo.Client
}

func NewProductRepository(client *mongo.Client) *ProductRepository {
	return &ProductRepository{client: client}
}

func (pr *ProductRepository) CreateProduct(ctx context.Context, product *models.Product) (*models.Product, error) {
	collection := config.GetCollection(pr.client, "products")

	result, err := collection.InsertOne(ctx, product)
	if err != nil {
		return nil, err
	}

	product.ID = result.InsertedID.(primitive.ObjectID)
	return product, nil
}

func (pr *ProductRepository) GetProductByID(ctx context.Context, productID primitive.ObjectID) (*models.Product, error) {
	collection := config.GetCollection(pr.client, "products")

	var product models.Product
	err := collection.FindOne(ctx, bson.M{"_id": productID}).Decode(&product)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &product, nil
}

func (pr *ProductRepository) GetAllProducts(ctx context.Context) ([]models.Product, error) {
	collection := config.GetCollection(pr.client, "products")

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var products []models.Product
	if err = cursor.All(ctx, &products); err != nil {
		return nil, err
	}

	return products, nil
}

func (pr *ProductRepository) GetProductsByCategory(ctx context.Context, category string) ([]models.Product, error) {
	collection := config.GetCollection(pr.client, "products")

	cursor, err := collection.Find(ctx, bson.M{"category": category})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var products []models.Product
	if err = cursor.All(ctx, &products); err != nil {
		return nil, err
	}

	return products, nil
}

func (pr *ProductRepository) UpdateProduct(ctx context.Context, productID primitive.ObjectID, product *models.Product) (*models.Product, error) {
	collection := config.GetCollection(pr.client, "products")

	updateData := bson.M{
		"name":        product.Name,
		"description": product.Description,
		"price":       product.Price,
		"size":        product.Size,
		"category":    product.Category,
		"imageURL":    product.ImageURL,
		"stock":       product.Stock,
		"updatedAt":   product.UpdatedAt,
	}

	result := collection.FindOneAndUpdate(
		ctx,
		bson.M{"_id": productID},
		bson.M{"$set": updateData},
	)

	if result.Err() != nil {
		return nil, result.Err()
	}

	var updatedProduct models.Product
	if err := result.Decode(&updatedProduct); err != nil {
		return nil, err
	}

	return &updatedProduct, nil
}

func (pr *ProductRepository) DeleteProduct(ctx context.Context, productID primitive.ObjectID) error {
	collection := config.GetCollection(pr.client, "products")

	_, err := collection.DeleteOne(ctx, bson.M{"_id": productID})
	return err
}

func (pr *ProductRepository) GetTotalProducts(ctx context.Context) (int64, error) {
	collection := config.GetCollection(pr.client, "products")

	count, err := collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (pr *ProductRepository) SearchProducts(ctx context.Context, query string) ([]models.Product, error) {
	collection := config.GetCollection(pr.client, "products")

	cursor, err := collection.Find(ctx, bson.M{
		"$or": []bson.M{
			{"name": bson.M{"$regex": query, "$options": "i"}},
			{"description": bson.M{"$regex": query, "$options": "i"}},
			{"category": bson.M{"$regex": query, "$options": "i"}},
		},
	})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var products []models.Product
	if err = cursor.All(ctx, &products); err != nil {
		return nil, err
	}

	return products, nil
}
