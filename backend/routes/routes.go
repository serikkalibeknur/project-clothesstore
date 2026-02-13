package routes

import (
	"github.com/gorilla/mux"
	"github.com/serikkalibeknur/project-clothesstore/internal/controllers"
	"github.com/serikkalibeknur/project-clothesstore/internal/middleware"
	"go.mongodb.org/mongo-driver/mongo"
)

func SetupRoutes(router *mux.Router, client *mongo.Client) {
	api := router.PathPrefix("/api").Subrouter()

	// Auth routes
	api.HandleFunc("/auth/register", middleware.LoggerMiddleware(controllers.Register(client))).Methods("POST")
	api.HandleFunc("/auth/login", middleware.LoggerMiddleware(controllers.Login(client))).Methods("POST")

	// Product routes (public)
	api.HandleFunc("/products", middleware.LoggerMiddleware(controllers.GetAllProducts(client))).Methods("GET")
	api.HandleFunc("/products/{id}", middleware.LoggerMiddleware(controllers.GetProductByID(client))).Methods("GET")

	// Product routes (admin only)
	api.HandleFunc("/products", middleware.LoggerMiddleware(middleware.AuthMiddleware(middleware.RequireAdmin(controllers.CreateProduct(client))))).Methods("POST")
	api.HandleFunc("/products/{id}", middleware.LoggerMiddleware(middleware.AuthMiddleware(middleware.RequireAdmin(controllers.UpdateProduct(client))))).Methods("PUT")
	api.HandleFunc("/products/{id}", middleware.LoggerMiddleware(middleware.AuthMiddleware(middleware.RequireAdmin(controllers.DeleteProduct(client))))).Methods("DELETE")

	// Cart routes (protected)
	api.HandleFunc("/cart", middleware.LoggerMiddleware(middleware.AuthMiddleware(controllers.GetCart(client)))).Methods("GET")
	api.HandleFunc("/cart", middleware.LoggerMiddleware(middleware.AuthMiddleware(controllers.AddToCart(client)))).Methods("POST")
	api.HandleFunc("/cart", middleware.LoggerMiddleware(middleware.AuthMiddleware(controllers.UpdateCart(client)))).Methods("PUT")
	api.HandleFunc("/cart/{productID}", middleware.LoggerMiddleware(middleware.AuthMiddleware(controllers.RemoveFromCart(client)))).Methods("DELETE")

	// Order routes (protected)
	api.HandleFunc("/orders", middleware.LoggerMiddleware(middleware.AuthMiddleware(controllers.CreateOrder(client)))).Methods("POST")
	api.HandleFunc("/orders", middleware.LoggerMiddleware(middleware.AuthMiddleware(controllers.GetUserOrders(client)))).Methods("GET")

	// Wishlist routes (protected)
	api.HandleFunc("/wishlist", middleware.LoggerMiddleware(middleware.AuthMiddleware(controllers.GetWishlist(client)))).Methods("GET")
	api.HandleFunc("/wishlist", middleware.LoggerMiddleware(middleware.AuthMiddleware(controllers.AddToWishlist(client)))).Methods("POST")
	api.HandleFunc("/wishlist/{productID}", middleware.LoggerMiddleware(middleware.AuthMiddleware(controllers.RemoveFromWishlist(client)))).Methods("DELETE")

	// Admin routes (protected, admin only)
	api.HandleFunc("/admin/statistics", middleware.LoggerMiddleware(middleware.AuthMiddleware(middleware.RequireAdmin(controllers.GetStatistics(client))))).Methods("GET")
	api.HandleFunc("/admin/orders", middleware.LoggerMiddleware(middleware.AuthMiddleware(middleware.RequireAdmin(controllers.GetAllOrders(client))))).Methods("GET")
	api.HandleFunc("/admin/orders/{id}", middleware.LoggerMiddleware(middleware.AuthMiddleware(middleware.RequireAdmin(controllers.UpdateOrderStatus(client))))).Methods("PUT")
}
