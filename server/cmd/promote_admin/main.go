package main

import (
	"flag"
	"log"

	"github.com/A-Words/ne-resource-community/server/internal/config"
	"github.com/A-Words/ne-resource-community/server/internal/database"
	"github.com/A-Words/ne-resource-community/server/internal/models"
)

func main() {
	email := flag.String("email", "", "Email of the user to promote to admin")
	flag.Parse()

	if *email == "" {
		log.Fatal("Please provide an email address using -email flag")
	}

	cfg := config.Load()
	db := database.New(cfg.DatabaseDSN)

	var user models.User
	if err := db.Where("email = ?", *email).First(&user).Error; err != nil {
		log.Fatalf("User with email %s not found: %v", *email, err)
	}

	if err := db.Model(&user).Update("role", "admin").Error; err != nil {
		log.Fatalf("Failed to promote user: %v", err)
	}

	log.Printf("Successfully promoted user %s to admin", *email)
}
