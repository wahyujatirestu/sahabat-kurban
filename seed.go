package main

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/wahyujatirestu/sahabat-kurban/model"
	"github.com/wahyujatirestu/sahabat-kurban/repository"
	"github.com/wahyujatirestu/sahabat-kurban/utils/security"
)

func seedInitialAdmin(userRepo repository.UserRepository) {
	ctx := context.Background()
	adminEmail := "admin@sahabatkurban.com"

	existing, _ := userRepo.FindByEmailOrUsername(ctx, adminEmail)
	if existing != nil {
		log.Println("Admin already exists. Skipping seed.")
		return
	}

	password, err := security.GeneratePasswordHash("superadmin123")
	if err != nil {
		log.Fatalf("Failed to hash password: %v", err)
	}

	admin := &model.User{
		ID:        uuid.New(),
		Username:  "superadmin",
		Name:      "Super Admin",
		Email:     adminEmail,
		Password:  password,
		Role:      "admin",
		Created_At: time.Now(),
		Updated_At: time.Now(),
	}

	if err := userRepo.Create(ctx, admin); err != nil {
		log.Fatalf("Failed to create initial admin: %v", err)
	}

	log.Println("Initial admin seeded successfully!")
}
