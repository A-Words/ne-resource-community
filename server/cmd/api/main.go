package main

import (
	"log"

	"github.com/A-Words/ne-resource-community/server/internal/config"
	"github.com/A-Words/ne-resource-community/server/internal/database"
	httpserver "github.com/A-Words/ne-resource-community/server/internal/http"
)

func main() {
	cfg := config.Load()

	db := database.New(cfg.DatabaseDSN)
	if err := database.AutoMigrate(db); err != nil {
		log.Fatalf("migrations failed: %v", err)
	}

	r := httpserver.NewRouter(db, cfg)
	if err := r.Run(cfg.Addr); err != nil {
		log.Fatalf("server stopped: %v", err)
	}
}
