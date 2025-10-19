package handlers

import (
	"context"
	"net/http"
	"os"
	"strconv"

	"jobconnect-backend/config"
	"jobconnect-backend/models"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/gin-gonic/gin"
)

// ApplyForJob - Job seeker applies for a job
func ApplyForJob(c *gin.Context) {
	userID, _ := c.Get("user_id")
	jobID := c.Param("id")

	// Check if user has uploaded resume
	var user config.User
	if err := config.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if user.ResumeURL == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Please upload your resume first"})
		return
	}

	// Check if job exists
	var job config.Job
	if err := config.DB.Where("id = ? AND status = ?", jobID, "active").First(&job).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Job not found or not active"})
		return
	}

	// Check if already applied
	var existingApp config.Application
	if err := config.DB.Where("job_id = ? AND user_id = ?", jobID, userID).First(&existingApp).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "You have already applied for this job"})
		return
	}

	// Get cover letter from request
	var req models.ApplyJobRequest
	c.ShouldBindJSON(&req)

	// Create application
	application := config.Application{
		JobID:       job.ID,
		UserID:      user.ID,
		ResumeURL:   user.ResumeURL,
		CoverLetter: req.CoverLetter,
		Status:      "pending",
	}

	if err := config.DB.Create(&application).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to submit application"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Application submitted successfully",
	})
}

// GetMyApplications - Job seeker gets their applications
func GetMyApplications(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var applications []config.Application
	config.DB.Preload("Job.Company").Where("user_id = ?", userID).Order("created_at DESC").Find(&applications)

	var response []models.ApplicationResponse
	for _, app := range applications {
		response = append(response, models.ApplicationResponse{
			ID: app.ID,
			Job: models.JobResponse{
				ID:    app.Job.ID,
				Title: app.Job.Title,
				Company: models.CompanyResponse{
					ID:      app.Job.Company.ID,
					Name:    app.Job.Company.Name,
					LogoURL: app.Job.Company.LogoURL,
				},
				Location: app.Job.Location,
				JobType:  app.Job.JobType,
				Salary:   app.Job.Salary,
			},
			ResumeURL:   app.ResumeURL,
			CoverLetter: app.CoverLetter,
			Status:      app.Status,
			CreatedAt:   app.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"success":      true,
		"applications": response,
	})
}

// GetJobApplications - Employer gets applications for their job
func GetJobApplications(c *gin.Context) {
	userID, _ := c.Get("user_id")
	jobID := c.Param("id")

	// Check if job belongs to user
	var job config.Job
	if err := config.DB.Where("id = ? AND posted_by = ?", jobID, userID).First(&job).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Job not found or unauthorized"})
		return
	}

	// Get applications
	var applications []config.Application
	config.DB.Preload("User").Where("job_id = ?", jobID).Order("created_at DESC").Find(&applications)

	var response []models.ApplicationResponse
	for _, app := range applications {
		response = append(response, models.ApplicationResponse{
			ID: app.ID,
			User: models.UserResponse{
				ID:        app.User.ID,
				Name:      app.User.Name,
				Email:     app.User.Email,
				Phone:     app.User.Phone,
				Location:  app.User.Location,
				Bio:       app.User.Bio,
				ResumeURL: app.User.ResumeURL,
			},
			ResumeURL:   app.ResumeURL,
			CoverLetter: app.CoverLetter,
			Status:      app.Status,
			CreatedAt:   app.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"success":      true,
		"applications": response,
	})
}

// UpdateApplicationStatus - Employer updates application status
func UpdateApplicationStatus(c *gin.Context) {
	userID, _ := c.Get("user_id")
	appID := c.Param("id")

	var req struct {
		Status string `json:"status" binding:"required"` // reviewed, shortlisted, rejected, accepted
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate status
	validStatuses := map[string]bool{
		"pending":     true,
		"reviewed":    true,
		"shortlisted": true,
		"rejected":    true,
		"accepted":    true,
	}

	if !validStatuses[req.Status] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status"})
		return
	}

	// Get application with job
	var application config.Application
	if err := config.DB.Preload("Job").Where("id = ?", appID).First(&application).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Application not found"})
		return
	}

	// Check if job belongs to user
	if application.Job.PostedBy != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized"})
		return
	}

	// Update status
	config.DB.Model(&application).Update("status", req.Status)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Application status updated",
	})
}

// UploadResume - Job seeker uploads resume
func UploadResume(c *gin.Context) {
	userID, _ := c.Get("user_id")

	// Get uploaded file
	file, err := c.FormFile("resume")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}

	// Validate file type (PDF only for simplicity)
	if file.Header.Get("Content-Type") != "application/pdf" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Only PDF files are allowed"})
		return
	}

	// Open file
	fileContent, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open file"})
		return
	}
	defer fileContent.Close()

	// Setup Cloudinary
	cld, err := cloudinary.NewFromParams(
		os.Getenv("CLOUDINARY_CLOUD_NAME"),
		os.Getenv("CLOUDINARY_API_KEY"),
		os.Getenv("CLOUDINARY_API_SECRET"),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to setup Cloudinary"})
		return
	}

	// Upload to Cloudinary
	ctx := context.Background()
	uploadResult, err := cld.Upload.Upload(ctx, fileContent, uploader.UploadParams{
		Folder:       "jobconnect/resumes",
		ResourceType: "auto",
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload resume"})
		return
	}

	// Update user's resume URL
	config.DB.Model(&config.User{}).Where("id = ?", userID).Update("resume_url", uploadResult.SecureURL)

	c.JSON(http.StatusOK, gin.H{
		"success":    true,
		"message":    "Resume uploaded successfully",
		"resume_url": uploadResult.SecureURL,
	})
}

// GetAllApplications - Admin gets all applications
func GetAllApplications(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset := (page - 1) * limit

	var applications []config.Application
	var totalCount int64

	config.DB.Model(&config.Application{}).Count(&totalCount)
	config.DB.Preload("Job.Company").Preload("User").
		Offset(offset).Limit(limit).
		Order("created_at DESC").
		Find(&applications)

	var response []models.ApplicationResponse
	for _, app := range applications {
		response = append(response, models.ApplicationResponse{
			ID: app.ID,
			Job: models.JobResponse{
				ID:    app.Job.ID,
				Title: app.Job.Title,
				Company: models.CompanyResponse{
					ID:   app.Job.Company.ID,
					Name: app.Job.Company.Name,
				},
			},
			User: models.UserResponse{
				ID:    app.User.ID,
				Name:  app.User.Name,
				Email: app.User.Email,
			},
			Status:    app.Status,
			CreatedAt: app.CreatedAt,
		})
	}

	totalPages := int((totalCount + int64(limit) - 1) / int64(limit))

	c.JSON(http.StatusOK, gin.H{
		"success":      true,
		"applications": response,
		"page":         page,
		"totalPages":   totalPages,
		"totalCount":   totalCount,
	})
}
