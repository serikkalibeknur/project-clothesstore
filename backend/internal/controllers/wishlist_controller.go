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

func GetWishlist(client *mongo.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value("userID").(string)
		userObjectID, err := primitive.ObjectIDFromHex(userID)
		if err != nil {
			utils.ErrorResponse(w, "Invalid user ID", http.StatusBadRequest)
			return
		}

		wishlistRepo := repository.NewWishlistRepository(client)
		wishlist, err := wishlistRepo.GetWishlist(r.Context(), userObjectID)
		if err != nil {
			utils.ErrorResponse(w, "Failed to fetch wishlist", http.StatusInternalServerError)
			return
		}

		if wishlist == nil {
			wishlist = &models.Wishlist{
				UserID:    userObjectID,
				Products:  []models.WishlistItem{},
				UpdatedAt: time.Now(),
			}
		}

		utils.SuccessResponse(w, "Wishlist fetched successfully", wishlist)
	}
}

func AddToWishlist(client *mongo.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value("userID").(string)
		userObjectID, err := primitive.ObjectIDFromHex(userID)
		if err != nil {
			utils.ErrorResponse(w, "Invalid user ID", http.StatusBadRequest)
			return
		}

		var input struct {
			ProductID string `json:"productID"`
		}
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			utils.ErrorResponse(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if input.ProductID == "" {
			utils.ErrorResponse(w, "Product ID is required", http.StatusBadRequest)
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

		wishlistRepo := repository.NewWishlistRepository(client)
		err = wishlistRepo.AddToWishlist(r.Context(), userObjectID, productID)
		if err != nil {
			utils.ErrorResponse(w, "Failed to add to wishlist", http.StatusInternalServerError)
			return
		}

		wishlist, _ := wishlistRepo.GetWishlist(r.Context(), userObjectID)
		utils.SuccessResponse(w, "Added to wishlist successfully", wishlist)
	}
}

func RemoveFromWishlist(client *mongo.Client) http.HandlerFunc {
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

		wishlistRepo := repository.NewWishlistRepository(client)
		err = wishlistRepo.RemoveFromWishlist(r.Context(), userObjectID, productObjectID)
		if err != nil {
			utils.ErrorResponse(w, "Failed to remove from wishlist", http.StatusInternalServerError)
			return
		}

		wishlist, _ := wishlistRepo.GetWishlist(r.Context(), userObjectID)
		utils.SuccessResponse(w, "Removed from wishlist successfully", wishlist)
	}
}
