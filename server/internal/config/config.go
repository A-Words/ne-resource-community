package config

import (
	"log"
	"os"
)

// Config holds server configuration loaded from environment variables.
type Config struct {
	Addr        string
	DatabaseDSN string
	JWTSecret   string
	UploadDir   string
	Env         string
	ClamAVAddr  string
}

// Load builds Config with sensible defaults; environment variables can override them.
func Load() Config {
	cfg := Config{
		Addr:        getEnv("SERVER_ADDR", ":8080"),
		DatabaseDSN: getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/ne_resource?sslmode=disable"),
		JWTSecret:   getEnv("JWT_SECRET", "dev-secret-change-me"),
		UploadDir:   getEnv("UPLOAD_DIR", "uploads"),
		Env:         getEnv("ENV", "dev"),
		ClamAVAddr:  getEnv("CLAMAV_ADDR", "tcp://localhost:3310"),
	}

	if err := os.MkdirAll(cfg.UploadDir, 0o755); err != nil {
		log.Fatalf("cannot create upload dir %s: %v", cfg.UploadDir, err)
	}
	return cfg
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
