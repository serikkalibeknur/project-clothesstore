package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/serikkalibeknur/project-clothesstore/internal/models"
	"github.com/serikkalibeknur/project-clothesstore/internal/repository"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	email := flag.String("email", "admin@example.com", "Admin email")
	password := flag.String("password", "Admin123456", "Admin password")
	name := flag.String("name", "Admin User", "Admin name")
	flag.Parse()

	// Connect to MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb+srv://sabyrovs_db_user:Rakhmet2007G@nimo.hrfgknv.mongodb.net/?appName=Nimo"))
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}
	defer client.Disconnect(ctx)

	// Create admin user
	userRepo := repository.NewUserRepository(client)

	// Check if user already exists
	existingUser, _ := userRepo.GetUserByEmail(ctx, *email)
	if existingUser != nil {
		fmt.Printf("User with email %s already exists. Updating role to admin...\n", *email)
		existingUser.Role = "admin"
		existingUser.UpdatedAt = time.Now()
		_, err := userRepo.UpdateUser(ctx, existingUser.ID, existingUser)
		if err != nil {
			log.Fatal("Failed to update user:", err)
		}
		fmt.Printf("✓ User %s upgraded to ADMIN role\n", *email)
		return
	}

	admin := &models.User{
		Name:      *name,
		Email:     *email,
		Password:  *password,
		Role:      "admin",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := admin.HashPassword(); err != nil {
		log.Fatal("Failed to hash password:", err)
	}

	createdAdmin, err := userRepo.CreateUser(ctx, admin)
	if err != nil {
		log.Fatal("Failed to create admin user:", err)
	}

	fmt.Println("\n✓ Admin account created successfully!")
	fmt.Printf("Email:    %s\n", createdAdmin.Email)
	fmt.Printf("Password: %s\n", *password)
	fmt.Printf("Role:     %s\n", createdAdmin.Role)
	fmt.Println("\nYou can now login with these credentials at http://localhost:8081/login.html")
}
