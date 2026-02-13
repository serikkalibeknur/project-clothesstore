package controllers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/serikkalibeknur/project-clothesstore/internal/models"
	"github.com/serikkalibeknur/project-clothesstore/internal/repository"
	"github.com/serikkalibeknur/project-clothesstore/internal/utils"
	"go.mongodb.org/mongo-driver/mongo"
)

func Register(client *mongo.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input models.RegisterInput
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			utils.ErrorResponse(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if input.Email == "" || input.Password == "" || input.Name == "" {
			utils.ErrorResponse(w, "Email, password, and name are required", http.StatusBadRequest)
			return
		}

		userRepo := repository.NewUserRepository(client)

		// Check if user already exists
		existingUser, _ := userRepo.GetUserByEmail(r.Context(), input.Email)
		if existingUser != nil {
			utils.ErrorResponse(w, "Email already registered", http.StatusConflict)
			return
		}

		user := &models.User{
			Name:      input.Name,
			Email:     input.Email,
			Password:  input.Password,
			Role:      "user",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		if err := user.HashPassword(); err != nil {
			utils.ErrorResponse(w, "Failed to process password", http.StatusInternalServerError)
			return
		}

		createdUser, err := userRepo.CreateUser(r.Context(), user)
		if err != nil {
			utils.ErrorResponse(w, "Failed to register user", http.StatusInternalServerError)
			return
		}

		createdUser.Password = ""
		utils.SuccessResponse(w, "User registered successfully", createdUser)
	}
}

func Login(client *mongo.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input models.LoginInput
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			utils.ErrorResponse(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if input.Email == "" || input.Password == "" {
			utils.ErrorResponse(w, "Email and password are required", http.StatusBadRequest)
			return
		}

		userRepo := repository.NewUserRepository(client)
		user, err := userRepo.GetUserByEmail(r.Context(), input.Email)
		if err != nil || user == nil {
			utils.ErrorResponse(w, "Invalid email or password", http.StatusUnauthorized)
			return
		}

		if !user.ComparePassword(input.Password) {
			utils.ErrorResponse(w, "Invalid email or password", http.StatusUnauthorized)
			return
		}

		token, err := utils.GenerateJWT(user.ID.Hex(), user.Email, user.Role)
		if err != nil {
			utils.ErrorResponse(w, "Failed to generate token", http.StatusInternalServerError)
			return
		}

		user.Password = ""
		response := map[string]interface{}{
			"user":  user,
			"token": token,
		}
		utils.SuccessResponse(w, "Login successful", response)
	}
}
