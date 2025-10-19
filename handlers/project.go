package handlers

import (
	"net/http"
	"strconv"

	"jobconnect-backend/config"
	"jobconnect-backend/models"

	"github.com/gin-gonic/gin"
)

// GetProjects - Browse all projects (homepage)
func GetProjects(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset := (page - 1) * limit

	categorySlug := c.Query("category")
	search := c.Query("search")

	var projects []config.Project
	var totalCount int64

	query := config.DB.Preload("User").Preload("Category").Preload("Images").Where("deleted_at IS NULL")

	// Filter by category
	if categorySlug != "" {
		var category config.Category
		if err := config.DB.Where("slug = ?", categorySlug).First(&category).Error; err == nil {
			query = query.Where("category_id = ?", category.ID)
		}
	}

	// Search
	if search != "" {
		query = query.Where("title ILIKE ? OR description ILIKE ? OR tags ILIKE ?",
			"%"+search+"%", "%"+search+"%", "%"+search+"%")
	}

	// Get count
	query.Model(&config.Project{}).Count(&totalCount)

	// Get projects
	query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&projects)

	// Build response
	var response []models.ProjectResponse
	for _, project := range projects {
		var images []models.ProjectImageResponse
		for _, img := range project.Images {
			images = append(images, models.ProjectImageResponse{
				ID:       img.ID,
				ImageURL: img.ImageURL,
				Order:    img.Order,
			})
		}

		response = append(response, models.ProjectResponse{
			ID:          project.ID,
			Title:       project.Title,
			Description: project.Description,
			CoverImage:  project.CoverImage,
			Images:      images,
			User: models.UserResponse{
				ID:        project.User.ID,
				Name:      project.User.Name,
				AvatarURL: project.User.AvatarURL,
			},
			Category: models.CategoryResponse{
				ID:   project.Category.ID,
				Name: project.Category.Name,
				Slug: project.Category.Slug,
				Icon: project.Category.Icon,
			},
			Tags:       project.Tags,
			Views:      project.Views,
			LikesCount: project.LikesCount,
			CreatedAt:  project.CreatedAt,
		})
	}

	totalPages := int((totalCount + int64(limit) - 1) / int64(limit))

	c.JSON(http.StatusOK, gin.H{
		"success":    true,
		"projects":   response,
		"page":       page,
		"totalPages": totalPages,
		"totalCount": totalCount,
	})
}

// GetProject - Get single project details
func GetProject(c *gin.Context) {
	projectID := c.Param("id")
	userID, userExists := c.Get("user_id")

	var project config.Project
	if err := config.DB.Preload("User").Preload("Category").Preload("Images").
		Where("id = ? AND deleted_at IS NULL", projectID).First(&project).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		return
	}

	// Increment views
	config.DB.Model(&project).Update("views", project.Views+1)

	// Check if user liked this project
	isLiked := false
	if userExists {
		var like config.Like
		if err := config.DB.Where("user_id = ? AND project_id = ?", userID, projectID).First(&like).Error; err == nil {
			isLiked = true
		}
	}

	// Build images response
	var images []models.ProjectImageResponse
	for _, img := range project.Images {
		images = append(images, models.ProjectImageResponse{
			ID:       img.ID,
			ImageURL: img.ImageURL,
			Order:    img.Order,
		})
	}

	response := models.ProjectResponse{
		ID:          project.ID,
		Title:       project.Title,
		Description: project.Description,
		CoverImage:  project.CoverImage,
		Images:      images,
		User: models.UserResponse{
			ID:        project.User.ID,
			Name:      project.User.Name,
			Email:     project.User.Email,
			Bio:       project.User.Bio,
			Location:  project.User.Location,
			AvatarURL: project.User.AvatarURL,
			ForHire:   project.User.ForHire,
		},
		Category: models.CategoryResponse{
			ID:   project.Category.ID,
			Name: project.Category.Name,
			Slug: project.Category.Slug,
			Icon: project.Category.Icon,
		},
		Tags:       project.Tags,
		Views:      project.Views + 1,
		LikesCount: project.LikesCount,
		IsLiked:    isLiked,
		CreatedAt:  project.CreatedAt,
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"project": response,
	})
}

// CreateProject - User uploads a project
func CreateProject(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req models.CreateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if len(req.ImageURLs) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "At least one image is required"})
		return
	}

	// Create project
	project := config.Project{
		Title:       req.Title,
		Description: req.Description,
		UserID:      userID.(uint),
		CategoryID:  req.CategoryID,
		Tags:        req.Tags,
		CoverImage:  req.ImageURLs[0], // First image is cover
	}

	if err := config.DB.Create(&project).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create project"})
		return
	}

	// Add images
	for i, imageURL := range req.ImageURLs {
		projectImage := config.ProjectImage{
			ProjectID: project.ID,
			ImageURL:  imageURL,
			Order:     i,
		}
		config.DB.Create(&projectImage)
	}

	c.JSON(http.StatusCreated, gin.H{
		"success":    true,
		"message":    "Project created successfully",
		"project_id": project.ID,
	})
}

// GetMyProjects - User gets their own projects
func GetMyProjects(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var projects []config.Project
	config.DB.Preload("Category").Preload("Images").
		Where("user_id = ? AND deleted_at IS NULL", userID).
		Order("created_at DESC").
		Find(&projects)

	var response []models.ProjectResponse
	for _, project := range projects {
		var images []models.ProjectImageResponse
		for _, img := range project.Images {
			images = append(images, models.ProjectImageResponse{
				ID:       img.ID,
				ImageURL: img.ImageURL,
				Order:    img.Order,
			})
		}

		response = append(response, models.ProjectResponse{
			ID:          project.ID,
			Title:       project.Title,
			Description: project.Description,
			CoverImage:  project.CoverImage,
			Images:      images,
			Category: models.CategoryResponse{
				ID:   project.Category.ID,
				Name: project.Category.Name,
				Slug: project.Category.Slug,
			},
			Tags:       project.Tags,
			Views:      project.Views,
			LikesCount: project.LikesCount,
			CreatedAt:  project.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"success":  true,
		"projects": response,
	})
}

// UpdateProject - User updates their project
func UpdateProject(c *gin.Context) {
	userID, _ := c.Get("user_id")
	projectID := c.Param("id")

	var req models.UpdateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check ownership
	var project config.Project
	if err := config.DB.Where("id = ? AND user_id = ?", projectID, userID).First(&project).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project not found or unauthorized"})
		return
	}

	// Update
	updates := make(map[string]interface{})
	if req.Title != "" {
		updates["title"] = req.Title
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.Tags != "" {
		updates["tags"] = req.Tags
	}

	config.DB.Model(&project).Updates(updates)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Project updated successfully",
	})
}

// DeleteProject - User deletes their project
func DeleteProject(c *gin.Context) {
	userID, _ := c.Get("user_id")
	projectID := c.Param("id")

	result := config.DB.Where("id = ? AND user_id = ?", projectID, userID).Delete(&config.Project{})

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project not found or unauthorized"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Project deleted successfully",
	})
}

// GetCategories - Get all categories
func GetCategories(c *gin.Context) {
	var categories []config.Category
	config.DB.Order("name").Find(&categories)

	var response []models.CategoryResponse
	for _, cat := range categories {
		response = append(response, models.CategoryResponse{
			ID:   cat.ID,
			Name: cat.Name,
			Slug: cat.Slug,
			Icon: cat.Icon,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"success":    true,
		"categories": response,
	})
}
