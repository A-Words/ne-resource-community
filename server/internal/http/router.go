package http

import (
	"github.com/A-Words/ne-resource-community/server/internal/config"
	"github.com/A-Words/ne-resource-community/server/internal/http/handlers"
	"github.com/A-Words/ne-resource-community/server/internal/http/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// NewRouter wires up gin.Engine with routes and middleware.
func NewRouter(db *gorm.DB, cfg config.Config) *gin.Engine {
	r := gin.Default()
	r.Use(cors.Default())

	authHandler := handlers.NewAuthHandler(db, cfg)
	resourceHandler := handlers.NewResourceHandler(db, cfg)

	api := r.Group("/api")
	{
		api.GET("/health", handlers.Health)

		auth := api.Group("/auth")
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)

		resources := api.Group("/resources")
		resources.GET("", resourceHandler.List)
		resources.GET(":id", resourceHandler.Get)
		resources.GET(":id/recommendations", resourceHandler.Recommend)

		protected := resources.Group("")
		protected.Use(middleware.AuthMiddleware(cfg))
		protected.POST("", resourceHandler.Create)
		protected.POST(":id/reviews", resourceHandler.Review)
		protected.GET(":id/download", resourceHandler.Download)
	}

	// Serve uploaded files statically for previews.
	r.Static("/uploads", cfg.UploadDir)

	return r
}
