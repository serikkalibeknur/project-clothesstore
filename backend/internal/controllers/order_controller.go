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

func CreateOrder(client *mongo.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value("userID").(string)
		userObjectID, err := primitive.ObjectIDFromHex(userID)
		if err != nil {
			utils.ErrorResponse(w, "Invalid user ID", http.StatusBadRequest)
			return
		}

		var input models.CreateOrderInput
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			utils.ErrorResponse(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if input.ShippingAddress.Street == "" || input.ShippingAddress.City == "" {
			utils.ErrorResponse(w, "Shipping address is required", http.StatusBadRequest)
			return
		}

		// Get user's cart
		cartRepo := repository.NewCartRepository(client)
		cart, err := cartRepo.GetCart(r.Context(), userObjectID)
		if err != nil || cart == nil || len(cart.Products) == 0 {
			utils.ErrorResponse(w, "Cart is empty", http.StatusBadRequest)
			return
		}

		// Convert cart items to order items and calculate total
		productRepo := repository.NewProductRepository(client)
		var orderItems []models.OrderItem
		totalPrice := 0.0

		for _, cartItem := range cart.Products {
			product, err := productRepo.GetProductByID(r.Context(), cartItem.ProductID)
			if err != nil || product == nil {
				utils.ErrorResponse(w, "Product not found", http.StatusNotFound)
				return
			}

			if product.Stock < cartItem.Quantity {
				utils.ErrorResponse(w, "Insufficient stock for product: "+product.Name, http.StatusBadRequest)
				return
			}

			orderItem := models.OrderItem{
				ProductID: cartItem.ProductID,
				Name:      product.Name,
				Price:     product.Price,
				Quantity:  cartItem.Quantity,
				Size:      cartItem.Size,
			}
			orderItems = append(orderItems, orderItem)
			totalPrice += product.Price * float64(cartItem.Quantity)

			// Decrease product stock
			newStock := product.Stock - cartItem.Quantity
			product.Stock = newStock
			productRepo.UpdateProduct(r.Context(), product.ID, product)
		}

		// Create order
		order := &models.Order{
			UserID:          userObjectID,
			Products:        orderItems,
			TotalPrice:      totalPrice,
			Status:          "pending",
			ShippingAddress: input.ShippingAddress,
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		}

		orderRepo := repository.NewOrderRepository(client)
		createdOrder, err := orderRepo.CreateOrder(r.Context(), order)
		if err != nil {
			utils.ErrorResponse(w, "Failed to create order", http.StatusInternalServerError)
			return
		}

		// Clear cart after successful order
		cartRepo.ClearCart(r.Context(), userObjectID)

		utils.SuccessResponse(w, "Order created successfully", createdOrder)
	}
}

func GetUserOrders(client *mongo.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value("userID").(string)
		userObjectID, err := primitive.ObjectIDFromHex(userID)
		if err != nil {
			utils.ErrorResponse(w, "Invalid user ID", http.StatusBadRequest)
			return
		}

		orderRepo := repository.NewOrderRepository(client)
		orders, err := orderRepo.GetUserOrders(r.Context(), userObjectID)
		if err != nil {
			utils.ErrorResponse(w, "Failed to fetch orders", http.StatusInternalServerError)
			return
		}

		if orders == nil {
			orders = []models.Order{}
		}

		utils.SuccessResponse(w, "Orders fetched successfully", orders)
	}
}

func GetAllOrders(client *mongo.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orderRepo := repository.NewOrderRepository(client)
		orders, err := orderRepo.GetAllOrders(r.Context())
		if err != nil {
			utils.ErrorResponse(w, "Failed to fetch orders", http.StatusInternalServerError)
			return
		}

		if orders == nil {
			orders = []models.Order{}
		}

		utils.SuccessResponse(w, "Orders fetched successfully", orders)
	}
}

func UpdateOrderStatus(client *mongo.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		orderID := vars["id"]

		objectID, err := primitive.ObjectIDFromHex(orderID)
		if err != nil {
			utils.ErrorResponse(w, "Invalid order ID", http.StatusBadRequest)
			return
		}

		var input struct {
			Status string `json:"status"`
		}
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			utils.ErrorResponse(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		validStatuses := map[string]bool{
			"pending":    true,
			"processing": true,
			"shipped":    true,
			"delivered":  true,
			"cancelled":  true,
		}

		if !validStatuses[input.Status] {
			utils.ErrorResponse(w, "Invalid status", http.StatusBadRequest)
			return
		}

		orderRepo := repository.NewOrderRepository(client)
		updatedOrder, err := orderRepo.UpdateOrderStatus(r.Context(), objectID, input.Status)
		if err != nil {
			utils.ErrorResponse(w, "Failed to update order status", http.StatusInternalServerError)
			return
		}

		utils.SuccessResponse(w, "Order status updated successfully", updatedOrder)
	}
}
