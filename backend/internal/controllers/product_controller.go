package controllers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/serikkalibeknur/project-clothesstore/internal/models"
	"github.com/serikkalibeknur/project-clothesstore/internal/repository"
	"github.com/serikkalibeknur/project-clothesstore/internal/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetAllProducts(client *mongo.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		productRepo := repository.NewProductRepository(client)

		products, err := productRepo.GetAllProducts(r.Context())
		if err != nil {
			utils.ErrorResponse(w, "Failed to fetch products", http.StatusInternalServerError)
			return
		}

		if products == nil {
			products = []models.Product{}
		}

		utils.SuccessResponse(w, "Products fetched successfully", products)
	}
}

func GetProductByID(client *mongo.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		productID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			utils.ErrorResponse(w, "Invalid product ID", http.StatusBadRequest)
			return
		}

		productRepo := repository.NewProductRepository(client)
		product, err := productRepo.GetProductByID(r.Context(), productID)
		if err != nil || product == nil {
			utils.ErrorResponse(w, "Product not found", http.StatusNotFound)
			return
		}

		utils.SuccessResponse(w, "Product fetched successfully", product)
	}
}

func CreateProduct(client *mongo.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input models.ProductInput
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			utils.ErrorResponse(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if input.Name == "" || input.Price <= 0 {
			utils.ErrorResponse(w, "Product name and price are required and price must be positive", http.StatusBadRequest)
			return
		}

		product := &models.Product{
			Name:        input.Name,
			Description: input.Description,
			Price:       input.Price,
			Size:        input.Size,
			Category:    input.Category,
			ImageURL:    input.ImageURL,
			Stock:       input.Stock,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		productRepo := repository.NewProductRepository(client)
		createdProduct, err := productRepo.CreateProduct(r.Context(), product)
		if err != nil {
			utils.ErrorResponse(w, "Failed to create product", http.StatusInternalServerError)
			return
		}

		utils.SuccessResponse(w, "Product created successfully", createdProduct)
	}
}

func UpdateProduct(client *mongo.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		productID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			utils.ErrorResponse(w, "Invalid product ID", http.StatusBadRequest)
			return
		}

		var input models.ProductInput
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			utils.ErrorResponse(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		productRepo := repository.NewProductRepository(client)

		// Get existing product
		existingProduct, err := productRepo.GetProductByID(r.Context(), productID)
		if err != nil || existingProduct == nil {
			utils.ErrorResponse(w, "Product not found", http.StatusNotFound)
			return
		}

		// Update fields
		existingProduct.Name = input.Name
		existingProduct.Description = input.Description
		existingProduct.Price = input.Price
		existingProduct.Size = input.Size
		existingProduct.Category = input.Category
		existingProduct.ImageURL = input.ImageURL
		existingProduct.Stock = input.Stock
		existingProduct.UpdatedAt = time.Now()

		updatedProduct, err := productRepo.UpdateProduct(r.Context(), productID, existingProduct)
		if err != nil {
			utils.ErrorResponse(w, "Failed to update product", http.StatusInternalServerError)
			return
		}

		utils.SuccessResponse(w, "Product updated successfully", updatedProduct)
	}
}

func DeleteProduct(client *mongo.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		productID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			utils.ErrorResponse(w, "Invalid product ID", http.StatusBadRequest)
			return
		}

		productRepo := repository.NewProductRepository(client)
		err = productRepo.DeleteProduct(r.Context(), productID)
		if err != nil {
			utils.ErrorResponse(w, "Failed to delete product", http.StatusInternalServerError)
			return
		}

		utils.SuccessResponse(w, "Product deleted successfully", nil)
	}
}
