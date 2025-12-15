package handlers

import (
	"net/http"

	"github.com/A-Words/ne-resource-community/server/internal/config"
	"github.com/A-Words/ne-resource-community/server/internal/http/middleware"
	"github.com/A-Words/ne-resource-community/server/internal/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RequestHandler struct {
	db  *gorm.DB
	cfg config.Config
}

func NewRequestHandler(db *gorm.DB, cfg config.Config) *RequestHandler {
	return &RequestHandler{db: db, cfg: cfg}
}

type requestCreateReq struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	Bounty      int    `json:"bounty"`
}

func (h *RequestHandler) Create(c *gin.Context) {
	userID, ok := middleware.UserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req requestCreateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if user has enough points if bounty > 0
	if req.Bounty > 0 {
		var user models.User
		if err := h.db.First(&user, "id = ?", userID).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch user"})
			return
		}
		if user.Points < req.Bounty {
			c.JSON(http.StatusBadRequest, gin.H{"error": "insufficient points"})
			return
		}
		// Deduct points immediately
		if err := h.db.Model(&user).UpdateColumn("points", gorm.Expr("points - ?", req.Bounty)).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to deduct points"})
			return
		}
	}

	request := models.Request{
		Title:       req.Title,
		Description: req.Description,
		Bounty:      req.Bounty,
		UserID:      userID,
		Status:      "open",
	}

	if err := h.db.Create(&request).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create request"})
		return
	}

	c.JSON(http.StatusCreated, request)
}

func (h *RequestHandler) List(c *gin.Context) {
	var requests []models.Request
	if err := h.db.Preload("User").Order("created_at DESC").Find(&requests).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "query failed"})
		return
	}
	c.JSON(http.StatusOK, requests)
}
