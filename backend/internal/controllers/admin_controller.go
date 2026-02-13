package controllers

import (
	"net/http"

	"github.com/serikkalibeknur/project-clothesstore/internal/repository"
	"github.com/serikkalibeknur/project-clothesstore/internal/utils"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetStatistics(client *mongo.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orderRepo := repository.NewOrderRepository(client)
		productRepo := repository.NewProductRepository(client)
		userRepo := repository.NewUserRepository(client)

		// Get statistics
		totalOrders, _ := orderRepo.GetTotalOrders(r.Context())
		totalProducts, _ := productRepo.GetTotalProducts(r.Context())
		totalUsers, _ := userRepo.GetTotalUsers(r.Context())
		totalRevenue, _ := orderRepo.GetTotalRevenue(r.Context())

		stats := map[string]interface{}{
			"totalOrders":   totalOrders,
			"totalProducts": totalProducts,
			"totalUsers":    totalUsers,
			"totalRevenue":  totalRevenue,
		}

		utils.SuccessResponse(w, "Statistics fetched successfully", stats)
	}
}
