package handlers

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
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
	Title        string `form:"title" binding:"required"`
	Description  string `form:"description"`
	Type         string `form:"type" binding:"required"`
	Vendor       string `form:"vendor"`
	DeviceModel  string `form:"deviceModel"`
	Protocol     string `form:"protocol"`
	Scenario     string `form:"scenario"`
	Tags         string `form:"tags"`
	ParentID     string `form:"parentId"`     // Optional
	Version      string `form:"version"`      // Optional, default 1.0
	ExternalLink string `form:"externalLink"` // Optional
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

	var diskPath, fileName, contentType, fileHash string

	if req.ExternalLink != "" {
		// If external link is provided, we skip file upload checks
		// We can add basic URL validation here if needed
	} else {
		file, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "file or external link is required"})
			return
		}

		// 1. Format Check (Simple extension check)
		ext := strings.ToLower(filepath.Ext(file.Filename))
		allowedExts := map[string]bool{
			".pdf": true, ".docx": true, ".doc": true, ".txt": true, ".md": true,
			".zip": true, ".rar": true, ".7z": true,
			".pcap": true, ".pcapng": true, ".gns3": true, ".pkt": true,
			".mp4": true,
		}
		if !allowedExts[ext] {
			c.JSON(http.StatusBadRequest, gin.H{"error": "unsupported file format"})
			return
		}

		// 2. Duplicate Check (Calculate Hash)
		src, err := file.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to open file"})
			return
		}
		defer src.Close()

		hash := sha256.New()
		if _, err := io.Copy(hash, src); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to calculate hash"})
			return
		}
		fileHash = hex.EncodeToString(hash.Sum(nil))

		// Reset file pointer for saving
		src.Seek(0, 0)

		var existing models.Resource
		if err := h.db.Where("file_hash = ?", fileHash).First(&existing).Error; err == nil {
			c.JSON(http.StatusConflict, gin.H{"error": "duplicate resource detected", "resourceId": existing.ID})
			return
		}

		safeName := fmt.Sprintf("%s%s", uuid.NewString(), ext)
		diskPath = filepath.Join(h.cfg.UploadDir, safeName)
		if err := c.SaveUploadedFile(file, diskPath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save file"})
			return
		}
		fileName = file.Filename
		contentType = file.Header.Get("Content-Type")
	}

	var parentID *uuid.UUID
	if req.ParentID != "" {
		if pid, err := uuid.Parse(req.ParentID); err == nil {
			parentID = &pid
		}
	}

	version := req.Version
	if version == "" {
		version = "1.0"
	}

	resource := models.Resource{
		Title:        req.Title,
		Description:  req.Description,
		Type:         req.Type,
		Vendor:       req.Vendor,
		DeviceModel:  req.DeviceModel,
		Protocol:     req.Protocol,
		Scenario:     req.Scenario,
		Tags:         req.Tags,
		FilePath:     diskPath,
		FileName:     fileName,
		ContentType:  contentType,
		FileHash:     fileHash,
		ExternalLink: req.ExternalLink,
		Status:       "pending", // Default to pending for audit
		UploaderID:   userID,
		ParentID:     parentID,
		Version:      version,
	}

	if err := h.db.Create(&resource).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save resource"})
		return
	}

	// Award points for contribution (Only after approval? Or now? Let's keep it now for simplicity, or maybe move to approval)
	// For better quality control, points should be awarded after approval.
	// h.db.Model(&models.User{}).Where("id = ?", userID).UpdateColumn("points", gorm.Expr("points + ?", 10))

	c.JSON(http.StatusCreated, resource)
}

type resourceQuery struct {
	Search   string `form:"search"`
	Type     string `form:"type"`
	Vendor   string `form:"vendor"`
	Device   string `form:"device"`
	Protocol string `form:"protocol"`
	Scenario string `form:"scenario"`
	Tag      string `form:"tag"`
	Sort     string `form:"sort"` // "newest", "downloads"
	Limit    int    `form:"limit,default=20"`
	Offset   int    `form:"offset,default=0"`
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
		Where("status = ?", "approved") // Only show approved resources

	if q.Sort == "downloads" {
		dbq = dbq.Order("download_count DESC")
	} else {
		dbq = dbq.Order("created_at DESC")
	}

	dbq = dbq.Limit(q.Limit).
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
		// Split search query by space to support multiple keywords
		keywords := strings.Fields(q.Search)
		for _, keyword := range keywords {
			pattern := "%" + keyword + "%"
			dbq = dbq.Where("title ILIKE ? OR description ILIKE ? OR tags ILIKE ? OR vendor ILIKE ? OR device_model ILIKE ?",
				pattern, pattern, pattern, pattern, pattern)
		}
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

// GetVersions returns the version history of a resource.
func (h *ResourceHandler) GetVersions(c *gin.Context) {
	id := c.Param("id")
	// Find the resource to get its chain.
	// This is a simple implementation: find all resources where ParentID is this ID, or this ID is ParentID...
	// Actually, a better way for versioning is usually a common GroupID.
	// But with ParentID linked list:
	// We can just find all resources that are in the same "chain".
	// For simplicity, let's just return resources that have this ID as ParentID (next versions)
	// and the resource pointed by ParentID (previous version).
	// Or, if we want full history, we might need a recursive query or just store a RootID.
	// Let's stick to: Find all resources that are connected.
	// For now, let's just return the immediate parent and children.

	// Better approach for this MVP: Just return resources where ParentID = id (children)
	// and the resource where ID = current.ParentID (parent).

	var current models.Resource
	if err := h.db.First(&current, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}

	var versions []models.Resource
	// Add current
	versions = append(versions, current)

	// Add parent if exists
	if current.ParentID != nil {
		var parent models.Resource
		if err := h.db.First(&parent, "id = ?", *current.ParentID).Error; err == nil {
			versions = append(versions, parent)
		}
	}

	// Add children (newer versions)
	var children []models.Resource
	h.db.Where("parent_id = ? AND status = 'approved'", current.ID).Find(&children)
	versions = append(versions, children...)

	c.JSON(http.StatusOK, versions)
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

	// Award points for interaction
	h.db.Model(&models.User{}).Where("id = ?", uid).UpdateColumn("points", gorm.Expr("points + ?", 2))

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

// ToggleFavorite adds or removes a resource from user's favorites.
func (h *ResourceHandler) ToggleFavorite(c *gin.Context) {
	uid, ok := middleware.UserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	id := c.Param("id")
	rid, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid resource id"})
		return
	}

	var fav models.Favorite
	result := h.db.Where("user_id = ? AND resource_id = ?", uid, rid).First(&fav)
	if result.Error == nil {
		// Exists, delete it
		h.db.Delete(&fav)
		c.JSON(http.StatusOK, gin.H{"status": "removed"})
	} else {
		// Not exists, create it
		h.db.Create(&models.Favorite{UserID: uid, ResourceID: rid})
		c.JSON(http.StatusOK, gin.H{"status": "added"})
	}
}

// ListFavorites returns resources favorited by the user.
func (h *ResourceHandler) ListFavorites(c *gin.Context) {
	uid, ok := middleware.UserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var favorites []models.Favorite
	if err := h.db.Where("user_id = ?", uid).Find(&favorites).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch favorites"})
		return
	}

	var resourceIDs []uuid.UUID
	for _, f := range favorites {
		resourceIDs = append(resourceIDs, f.ResourceID)
	}

	var resources []models.Resource
	if len(resourceIDs) > 0 {
		h.db.Where("id IN ?", resourceIDs).Find(&resources)
	} else {
		resources = []models.Resource{}
	}

	c.JSON(http.StatusOK, resources)
}

// ListDownloads returns resources downloaded by the user.
func (h *ResourceHandler) ListDownloads(c *gin.Context) {
	uid, ok := middleware.UserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var logs []models.DownloadLog
	if err := h.db.Where("user_id = ?", uid).Order("created_at DESC").Find(&logs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch download history"})
		return
	}

	var resourceIDs []uuid.UUID
	uniqueIDs := make(map[uuid.UUID]bool)
	for _, l := range logs {
		if !uniqueIDs[l.ResourceID] {
			uniqueIDs[l.ResourceID] = true
			resourceIDs = append(resourceIDs, l.ResourceID)
		}
	}

	var resources []models.Resource
	if len(resourceIDs) > 0 {
		h.db.Where("id IN ?", resourceIDs).Find(&resources)
	} else {
		resources = []models.Resource{}
	}

	c.JSON(http.StatusOK, resources)
}

// --- Quality Control & Admin ---

// ReportResource allows users to report a resource.
func (h *ResourceHandler) ReportResource(c *gin.Context) {
	uid, ok := middleware.UserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	id := c.Param("id")
	rid, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid resource id"})
		return
	}

	var req struct {
		Reason string `json:"reason" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	report := models.Report{
		UserID:     uid,
		ResourceID: rid,
		Reason:     req.Reason,
		Status:     "pending",
	}
	if err := h.db.Create(&report).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to submit report"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "report submitted"})
}

// AdminListPending returns resources waiting for audit.
func (h *ResourceHandler) AdminListPending(c *gin.Context) {
	// In real app, check if user is admin
	var resources []models.Resource
	if err := h.db.Preload("Uploader").Where("status = ?", "pending").Find(&resources).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "query failed"})
		return
	}
	c.JSON(http.StatusOK, resources)
}

// AdminAuditResource approves or rejects a resource.
func (h *ResourceHandler) AdminAuditResource(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Action string `json:"action" binding:"required,oneof=approve reject"`
		Reason string `json:"reason"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var resource models.Resource
	if err := h.db.First(&resource, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}

	if req.Action == "approve" {
		resource.Status = "approved"
		// Award points now
		h.db.Model(&models.User{}).Where("id = ?", resource.UploaderID).UpdateColumn("points", gorm.Expr("points + ?", 10))
	} else {
		resource.Status = "rejected"
		resource.RejectReason = req.Reason
	}

	if err := h.db.Save(&resource).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update status"})
		return
	}

	c.JSON(http.StatusOK, resource)
}

// AdminListReports returns all pending reports.
func (h *ResourceHandler) AdminListReports(c *gin.Context) {
	var reports []models.Report
	// Preload Resource and User to show details
	if err := h.db.Preload("Resource").Preload("User").Where("status = ?", "pending").Find(&reports).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch reports"})
		return
	}
	c.JSON(http.StatusOK, reports)
}

// AdminResolveReport marks a report as resolved.
func (h *ResourceHandler) AdminResolveReport(c *gin.Context) {
	id := c.Param("id")
	var report models.Report
	if err := h.db.First(&report, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}

	report.Status = "resolved"
	if err := h.db.Save(&report).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update report"})
		return
	}

	c.JSON(http.StatusOK, report)
}

// GetPopularTags returns a list of popular tags.
func (h *ResourceHandler) GetPopularTags(c *gin.Context) {
	type TagResult struct {
		Tag   string `json:"tag"`
		Count int64  `json:"count"`
	}
	var results []TagResult

	// Postgres specific query for comma separated tags
	err := h.db.Raw(`
		SELECT trim(tag) as tag, count(*) as count
		FROM (
			SELECT unnest(string_to_array(tags, ',')) as tag
			FROM resources
			WHERE status = 'approved' AND tags != ''
		) t
		WHERE trim(tag) != ''
		GROUP BY tag
		ORDER BY count DESC
		LIMIT 30
	`).Scan(&results).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch tags"})
		return
	}

	c.JSON(http.StatusOK, results)
}
