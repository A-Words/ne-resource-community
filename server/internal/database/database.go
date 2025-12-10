package database

import (
	"fmt"
	"log"
	"time"

	"github.com/A-Words/ne-resource-community/server/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// New connects to PostgreSQL using GORM with sane defaults.
func New(dsn string) *gorm.DB {
	gormCfg := &gorm.Config{Logger: logger.Default.LogMode(logger.Info)}
	db, err := gorm.Open(postgres.Open(dsn), gormCfg)
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("failed to get sql DB: %v", err)
	}
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)
	return db
}

// AutoMigrate creates tables and full-text index for resources.
func AutoMigrate(db *gorm.DB) error {
	if err := db.AutoMigrate(
		&models.User{},
		&models.Resource{},
		&models.Review{},
		&models.Favorite{},
		&models.DownloadLog{},
		&models.Report{},
	); err != nil {
		return fmt.Errorf("automigrate: %w", err)
	}

	// Ensure search_vector column and GIN index exist for full-text search.
	resourceFTS := `
ALTER TABLE resources
	ADD COLUMN IF NOT EXISTS search_vector tsvector GENERATED ALWAYS AS (
		setweight(to_tsvector('english', coalesce(title, '')), 'A') ||
		setweight(to_tsvector('english', coalesce(description, '')), 'B') ||
		setweight(to_tsvector('english', coalesce(tags, '')), 'C')
	) STORED;
CREATE INDEX IF NOT EXISTS idx_resources_search_vector ON resources USING GIN (search_vector);
`
	if err := db.Exec(resourceFTS).Error; err != nil {
		return fmt.Errorf("ensure fts: %w", err)
	}

	return nil
}
