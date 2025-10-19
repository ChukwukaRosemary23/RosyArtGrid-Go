package handlers

import (
	"net/http"
	"strconv"

	"jobconnect-backend/config"
	"jobconnect-backend/models"

	"github.com/gin-gonic/gin"
)

// GetJobs - Get all active jobs (public)
func GetJobs(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset := (page - 1) * limit

	search := c.Query("search")
	jobType := c.Query("job_type")
	location := c.Query("location")

	var jobs []config.Job
	var totalCount int64

	query := config.DB.Preload("Company").Where("status = ? AND deleted_at IS NULL", "active")

	// Apply filters
	if search != "" {
		query = query.Where("title ILIKE ? OR description ILIKE ?", "%"+search+"%", "%"+search+"%")
	}
	if jobType != "" {
		query = query.Where("job_type = ?", jobType)
	}
	if location != "" {
		query = query.Where("location ILIKE ?", "%"+location+"%")
	}

	// Get total count
	query.Model(&config.Job{}).Count(&totalCount)

	// Get paginated results
	query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&jobs)

	var response []models.JobResponse
	for _, job := range jobs {
		response = append(response, models.JobResponse{
			ID:          job.ID,
			Title:       job.Title,
			Description: job.Description,
			Company: models.CompanyResponse{
				ID:       job.Company.ID,
				Name:     job.Company.Name,
				LogoURL:  job.Company.LogoURL,
				Location: job.Company.Location,
			},
			Location:   job.Location,
			JobType:    job.JobType,
			Salary:     job.Salary,
			Experience: job.Experience,
			Skills:     job.Skills,
			Status:     job.Status,
			CreatedAt:  job.CreatedAt,
		})
	}

	totalPages := int((totalCount + int64(limit) - 1) / int64(limit))

	c.JSON(http.StatusOK, gin.H{
		"success":    true,
		"jobs":       response,
		"page":       page,
		"limit":      limit,
		"totalPages": totalPages,
		"totalCount": totalCount,
	})
}

// GetJob - Get single job by ID (public)
func GetJob(c *gin.Context) {
	jobID := c.Param("id")

	var job config.Job
	if err := config.DB.Preload("Company").Preload("Poster").Where("id = ? AND deleted_at IS NULL", jobID).First(&job).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Job not found"})
		return
	}

	response := models.JobResponse{
		ID:          job.ID,
		Title:       job.Title,
		Description: job.Description,
		Company: models.CompanyResponse{
			ID:          job.Company.ID,
			Name:        job.Company.Name,
			Description: job.Company.Description,
			Website:     job.Company.Website,
			Location:    job.Company.Location,
			LogoURL:     job.Company.LogoURL,
			Industry:    job.Company.Industry,
			Size:        job.Company.Size,
		},
		Location:   job.Location,
		JobType:    job.JobType,
		Salary:     job.Salary,
		Experience: job.Experience,
		Skills:     job.Skills,
		Status:     job.Status,
		CreatedAt:  job.CreatedAt,
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"job":     response,
	})
}

// CreateJob - Employer creates a job (protected - employer only)
func CreateJob(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req models.CreateJobRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user's company
	var user config.User
	if err := config.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if user.CompanyID == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Please create a company profile first"})
		return
	}

	// Create job
	job := config.Job{
		Title:       req.Title,
		Description: req.Description,
		CompanyID:   *user.CompanyID,
		Location:    req.Location,
		JobType:     req.JobType,
		Salary:      req.Salary,
		Experience:  req.Experience,
		Skills:      req.Skills,
		Status:      "active",
		PostedBy:    user.ID,
	}

	if err := config.DB.Create(&job).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create job"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Job posted successfully",
		"job_id":  job.ID,
	})
}

// GetMyJobs - Employer gets their posted jobs
func GetMyJobs(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var jobs []config.Job
	config.DB.Preload("Company").Where("posted_by = ? AND deleted_at IS NULL", userID).Order("created_at DESC").Find(&jobs)

	var response []models.JobResponse
	for _, job := range jobs {
		response = append(response, models.JobResponse{
			ID:          job.ID,
			Title:       job.Title,
			Description: job.Description,
			Company: models.CompanyResponse{
				ID:      job.Company.ID,
				Name:    job.Company.Name,
				LogoURL: job.Company.LogoURL,
			},
			Location:   job.Location,
			JobType:    job.JobType,
			Salary:     job.Salary,
			Experience: job.Experience,
			Skills:     job.Skills,
			Status:     job.Status,
			CreatedAt:  job.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"jobs":    response,
	})
}

// UpdateJob - Employer updates their job
func UpdateJob(c *gin.Context) {
	userID, _ := c.Get("user_id")
	jobID := c.Param("id")

	var req models.UpdateJobRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if job belongs to user
	var job config.Job
	if err := config.DB.Where("id = ? AND posted_by = ?", jobID, userID).First(&job).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Job not found or unauthorized"})
		return
	}

	// Update job
	updates := make(map[string]interface{})
	if req.Title != "" {
		updates["title"] = req.Title
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.Location != "" {
		updates["location"] = req.Location
	}
	if req.JobType != "" {
		updates["job_type"] = req.JobType
	}
	if req.Salary != "" {
		updates["salary"] = req.Salary
	}
	if req.Experience != "" {
		updates["experience"] = req.Experience
	}
	if req.Skills != "" {
		updates["skills"] = req.Skills
	}
	if req.Status != "" {
		updates["status"] = req.Status
	}

	config.DB.Model(&job).Updates(updates)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Job updated successfully",
	})
}

// DeleteJob - Employer soft deletes their job
func DeleteJob(c *gin.Context) {
	userID, _ := c.Get("user_id")
	jobID := c.Param("id")

	result := config.DB.Where("id = ? AND posted_by = ?", jobID, userID).Delete(&config.Job{})

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Job not found or unauthorized"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Job deleted successfully",
	})
}
