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

func GetCart(client *mongo.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value("userID").(string)
		userObjectID, err := primitive.ObjectIDFromHex(userID)
		if err != nil {
			utils.ErrorResponse(w, "Invalid user ID", http.StatusBadRequest)
			return
		}

		cartRepo := repository.NewCartRepository(client)
		cart, err := cartRepo.GetCart(r.Context(), userObjectID)
		if err != nil {
			utils.ErrorResponse(w, "Failed to fetch cart", http.StatusInternalServerError)
			return
		}

		if cart == nil {
			cart = &models.Cart{
				UserID:    userObjectID,
				Products:  []models.CartItem{},
				UpdatedAt: time.Now(),
			}
		}

		utils.SuccessResponse(w, "Cart fetched successfully", cart)
	}
}

func AddToCart(client *mongo.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value("userID").(string)
		userObjectID, err := primitive.ObjectIDFromHex(userID)
		if err != nil {
			utils.ErrorResponse(w, "Invalid user ID", http.StatusBadRequest)
			return
		}

		var input models.AddToCartInput
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			utils.ErrorResponse(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if input.ProductID == "" || input.Quantity <= 0 {
			utils.ErrorResponse(w, "Product ID and quantity are required", http.StatusBadRequest)
			return
		}

		productID, err := primitive.ObjectIDFromHex(input.ProductID)
		if err != nil {
			utils.ErrorResponse(w, "Invalid product ID", http.StatusBadRequest)
			return
		}

		// Verify product exists
		productRepo := repository.NewProductRepository(client)
		product, err := productRepo.GetProductByID(r.Context(), productID)
		if err != nil || product == nil {
			utils.ErrorResponse(w, "Product not found", http.StatusNotFound)
			return
		}

		cartRepo := repository.NewCartRepository(client)
		err = cartRepo.AddToCart(r.Context(), userObjectID, productID, input.Quantity, input.Size)
		if err != nil {
			utils.ErrorResponse(w, "Failed to add to cart", http.StatusInternalServerError)
			return
		}

		cart, _ := cartRepo.GetCart(r.Context(), userObjectID)
		utils.SuccessResponse(w, "Added to cart successfully", cart)
	}
}

func UpdateCart(client *mongo.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value("userID").(string)
		userObjectID, err := primitive.ObjectIDFromHex(userID)
		if err != nil {
			utils.ErrorResponse(w, "Invalid user ID", http.StatusBadRequest)
			return
		}

		var input models.AddToCartInput
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			utils.ErrorResponse(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if input.ProductID == "" || input.Quantity < 0 {
			utils.ErrorResponse(w, "Product ID and quantity are required", http.StatusBadRequest)
			return
		}

		productID, err := primitive.ObjectIDFromHex(input.ProductID)
		if err != nil {
			utils.ErrorResponse(w, "Invalid product ID", http.StatusBadRequest)
			return
		}

		cartRepo := repository.NewCartRepository(client)
		if input.Quantity == 0 {
			err = cartRepo.RemoveFromCart(r.Context(), userObjectID, productID)
		} else {
			err = cartRepo.UpdateCartItemQuantity(r.Context(), userObjectID, productID, input.Quantity)
		}

		if err != nil {
			utils.ErrorResponse(w, "Failed to update cart", http.StatusInternalServerError)
			return
		}

		cart, _ := cartRepo.GetCart(r.Context(), userObjectID)
		utils.SuccessResponse(w, "Cart updated successfully", cart)
	}
}

func RemoveFromCart(client *mongo.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value("userID").(string)
		userObjectID, err := primitive.ObjectIDFromHex(userID)
		if err != nil {
			utils.ErrorResponse(w, "Invalid user ID", http.StatusBadRequest)
			return
		}

		vars := mux.Vars(r)
		productID := vars["productID"]

		productObjectID, err := primitive.ObjectIDFromHex(productID)
		if err != nil {
			utils.ErrorResponse(w, "Invalid product ID", http.StatusBadRequest)
			return
		}

		cartRepo := repository.NewCartRepository(client)
		err = cartRepo.RemoveFromCart(r.Context(), userObjectID, productObjectID)
		if err != nil {
			utils.ErrorResponse(w, "Failed to remove from cart", http.StatusInternalServerError)
			return
		}

		cart, _ := cartRepo.GetCart(r.Context(), userObjectID)
		utils.SuccessResponse(w, "Removed from cart successfully", cart)
	}
}
