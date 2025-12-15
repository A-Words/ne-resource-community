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
	requestHandler := handlers.NewRequestHandler(db, cfg)

	api := r.Group("/api")
	{
		api.GET("/health", handlers.Health)

		auth := api.Group("/auth")
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)

		resources := api.Group("/resources")
		resources.GET("", resourceHandler.List)
		resources.GET("/tags/popular", resourceHandler.GetPopularTags)
		resources.GET(":id", resourceHandler.Get)
		resources.GET(":id/recommendations", resourceHandler.Recommend)
		resources.GET(":id/versions", resourceHandler.GetVersions)

		protected := resources.Group("")
		protected.Use(middleware.AuthMiddleware(cfg))
		protected.POST("", resourceHandler.Create)
		protected.POST(":id/reviews", resourceHandler.Review)
		protected.POST(":id/favorite", resourceHandler.ToggleFavorite)
		protected.POST(":id/report", resourceHandler.ReportResource)
		protected.GET(":id/download", resourceHandler.Download)
		protected.POST(":id/progress", resourceHandler.UpdateProgress)
		protected.GET(":id/progress", resourceHandler.GetProgress)

		user := api.Group("/user")
		user.Use(middleware.AuthMiddleware(cfg))
		user.POST("/change-password", authHandler.ChangePassword)
		user.GET("/favorites", resourceHandler.ListFavorites)
		user.GET("/downloads", resourceHandler.ListDownloads)
		user.GET("/uploads", resourceHandler.ListMyUploads)

		// Admin routes (simplified, reusing auth middleware but should check role)
		admin := api.Group("/admin")
		admin.Use(middleware.AuthMiddleware(cfg))
		admin.GET("/pending", resourceHandler.AdminListPending)
		admin.POST("/resources/:id/audit", resourceHandler.AdminAuditResource)
		admin.GET("/reports", resourceHandler.AdminListReports)
		admin.POST("/reports/:id/resolve", resourceHandler.AdminResolveReport)

		requests := api.Group("/requests")
		requests.GET("", requestHandler.List)
		requests.POST("", middleware.AuthMiddleware(cfg), requestHandler.Create)
	}

	// Serve uploaded files statically for previews.
	r.Static("/uploads", cfg.UploadDir)

	return r
}
