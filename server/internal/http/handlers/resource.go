package handlers

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/A-Words/ne-resource-community/server/internal/config"
	"github.com/A-Words/ne-resource-community/server/internal/http/middleware"
	"github.com/A-Words/ne-resource-community/server/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ResourceHandler manages resource CRUD and search.
type ResourceHandler struct {
	db  *gorm.DB
	cfg config.Config
}

func NewResourceHandler(db *gorm.DB, cfg config.Config) *ResourceHandler {
	return &ResourceHandler{db: db, cfg: cfg}
}

type resourceCreateReq struct {
	Title       string `form:"title" binding:"required"`
	Description string `form:"description"`
	Type        string `form:"type" binding:"required"`
	Vendor      string `form:"vendor"`
	DeviceModel string `form:"deviceModel"`
	Protocol    string `form:"protocol"`
	Scenario    string `form:"scenario"`
	Tags        string `form:"tags"`
}

// Create handles multipart upload, stores file locally, and records resource metadata.
func (h *ResourceHandler) Create(c *gin.Context) {
	userID, ok := middleware.UserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req resourceCreateReq
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
		return
	}

	safeName := fmt.Sprintf("%s%s", uuid.NewString(), filepath.Ext(file.Filename))
	diskPath := filepath.Join(h.cfg.UploadDir, safeName)
	if err := c.SaveUploadedFile(file, diskPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save file"})
		return
	}

	resource := models.Resource{
		Title:       req.Title,
		Description: req.Description,
		Type:        req.Type,
		Vendor:      req.Vendor,
		DeviceModel: req.DeviceModel,
		Protocol:    req.Protocol,
		Scenario:    req.Scenario,
		Tags:        req.Tags,
		FilePath:    diskPath,
		FileName:    file.Filename,
		ContentType: file.Header.Get("Content-Type"),
		UploaderID:  userID,
	}

	if err := h.db.Create(&resource).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save resource"})
		return
	}

	c.JSON(http.StatusCreated, resource)
}

type resourceQuery struct {
	Search     string `form:"search"`
	Type       string `form:"type"`
	Vendor     string `form:"vendor"`
	Device     string `form:"device"`
	Protocol   string `form:"protocol"`
	Scenario   string `form:"scenario"`
	Tag        string `form:"tag"`
	Limit      int    `form:"limit,default=20"`
	Offset     int    `form:"offset,default=0"`
}

// List returns resources with filters and full-text search.
func (h *ResourceHandler) List(c *gin.Context) {
	var q resourceQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	dbq := h.db.Model(&models.Resource{}).
		Preload("Uploader").
		Order("created_at DESC").
		Limit(q.Limit).
		Offset(q.Offset)

	if q.Type != "" {
		dbq = dbq.Where("type = ?", q.Type)
	}
	if q.Vendor != "" {
		dbq = dbq.Where("vendor ILIKE ?", "%"+q.Vendor+"%")
	}
	if q.Device != "" {
		dbq = dbq.Where("device_model ILIKE ?", "%"+q.Device+"%")
	}
	if q.Protocol != "" {
		dbq = dbq.Where("protocol ILIKE ?", "%"+q.Protocol+"%")
	}
	if q.Scenario != "" {
		dbq = dbq.Where("scenario ILIKE ?", "%"+q.Scenario+"%")
	}
	if q.Tag != "" {
		tag := strings.ToLower(q.Tag)
		dbq = dbq.Where("LOWER(tags) LIKE ?", "%"+tag+"%")
	}
	if q.Search != "" {
		dbq = dbq.Where("search_vector @@ websearch_to_tsquery('english', ?)", q.Search)
	}

	var resources []models.Resource
	if err := dbq.Find(&resources).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "query failed"})
		return
	}

	c.JSON(http.StatusOK, resources)
}

// Get returns a single resource by id.
func (h *ResourceHandler) Get(c *gin.Context) {
	id := c.Param("id")
	var resource models.Resource
	if err := h.db.Preload("Uploader").First(&resource, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "query failed"})
		return
	}
	c.JSON(http.StatusOK, resource)
}

// Download increments counters and streams the file.
func (h *ResourceHandler) Download(c *gin.Context) {
	id := c.Param("id")
	var resource models.Resource
	if err := h.db.First(&resource, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}

	h.db.Model(&resource).UpdateColumn("download_count", gorm.Expr("download_count + 1"))
	if uid, ok := middleware.UserID(c); ok {
		h.db.Create(&models.DownloadLog{UserID: uid, ResourceID: resource.ID})
	}

	c.FileAttachment(resource.FilePath, resource.FileName)
}

type reviewReq struct {
	Score   int    `json:"score" binding:"required,min=1,max=5"`
	Comment string `json:"comment"`
}

// Review allows authenticated users to rate and comment on a resource.
func (h *ResourceHandler) Review(c *gin.Context) {
	uid, ok := middleware.UserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	id := c.Param("id")
	var resource models.Resource
	if err := h.db.First(&resource, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}

	var req reviewReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	rev := models.Review{ResourceID: resource.ID, UserID: uid, Score: req.Score, Comment: req.Comment}
	if err := h.db.Create(&rev).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save review"})
		return
	}

	h.db.Model(&resource).Updates(map[string]interface{}{
		"rating_count":   gorm.Expr("rating_count + 1"),
		"rating_average": gorm.Expr("((rating_average * rating_count) + ?) / (rating_count + 1)", req.Score),
	})

	c.JSON(http.StatusCreated, rev)
}

// Recommend suggests similar resources by vendor/type.
func (h *ResourceHandler) Recommend(c *gin.Context) {
	id := c.Param("id")
	var resource models.Resource
	if err := h.db.First(&resource, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}

	var related []models.Resource
	h.db.Where("id <> ? AND (vendor = ? OR type = ?)", resource.ID, resource.Vendor, resource.Type).
		Order("rating_average DESC, download_count DESC").
		Limit(5).
		Find(&related)

	c.JSON(http.StatusOK, related)
}
