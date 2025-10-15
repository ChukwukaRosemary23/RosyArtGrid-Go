package handlers

import (
	"net/http"
	"strconv"

	"jobconnect-backend/config"
	"jobconnect-backend/models"

	"github.com/gin-gonic/gin"
)

// GetAllUsers - Admin gets all users
func GetAllUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset := (page - 1) * limit

	role := c.Query("role") // Filter by role

	var users []config.User
	var totalCount int64

	query := config.DB.Preload("Company")

	if role != "" {
		query = query.Where("role = ?", role)
	}

	query.Model(&config.User{}).Count(&totalCount)
	query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&users)

	var response []models.UserResponse
	for _, user := range users {
		userResp := models.UserResponse{
			ID:        user.ID,
			Name:      user.Name,
			Email:     user.Email,
			Role:      user.Role,
			Phone:     user.Phone,
			Location:  user.Location,
			ResumeURL: user.ResumeURL,
			CreatedAt: user.CreatedAt,
		}

		if user.Company != nil {
			userResp.Company = &models.CompanyResponse{
				ID:   user.Company.ID,
				Name: user.Company.Name,
			}
		}

		response = append(response, userResp)
	}

	totalPages := int((totalCount + int64(limit) - 1) / int64(limit))

	c.JSON(http.StatusOK, gin.H{
		"success":    true,
		"users":      response,
		"page":       page,
		"totalPages": totalPages,
		"totalCount": totalCount,
	})
}

// GetAllJobs - Admin gets all jobs (including inactive)
func GetAllJobsAdmin(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset := (page - 1) * limit

	status := c.Query("status") // Filter by status

	var jobs []config.Job
	var totalCount int64

	query := config.DB.Preload("Company").Where("deleted_at IS NULL")

	if status != "" {
		query = query.Where("status = ?", status)
	}

	query.Model(&config.Job{}).Count(&totalCount)
	query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&jobs)

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
			Status:     job.Status,
			CreatedAt:  job.CreatedAt,
		})
	}

	totalPages := int((totalCount + int64(limit) - 1) / int64(limit))

	c.JSON(http.StatusOK, gin.H{
		"success":    true,
		"jobs":       response,
		"page":       page,
		"totalPages": totalPages,
		"totalCount": totalCount,
	})
}

// UpdateJobStatus - Admin updates job status (approve/reject)
func UpdateJobStatus(c *gin.Context) {
	jobID := c.Param("id")

	var req struct {
		Status string `json:"status" binding:"required"` // active, closed, pending
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result := config.DB.Model(&config.Job{}).Where("id = ?", jobID).Update("status", req.Status)

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Job not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Job status updated successfully",
	})
}

// DeleteUser - Admin deletes a user
func DeleteUser(c *gin.Context) {
	userID := c.Param("id")

	result := config.DB.Delete(&config.User{}, userID)

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "User deleted successfully",
	})
}

// GetDashboardStats - Admin gets dashboard statistics
func GetDashboardStats(c *gin.Context) {
	var totalUsers int64
	var totalJobs int64
	var totalApplications int64
	var activeJobs int64

	config.DB.Model(&config.User{}).Count(&totalUsers)
	config.DB.Model(&config.Job{}).Where("deleted_at IS NULL").Count(&totalJobs)
	config.DB.Model(&config.Application{}).Count(&totalApplications)
	config.DB.Model(&config.Job{}).Where("status = ? AND deleted_at IS NULL", "active").Count(&activeJobs)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"stats": gin.H{
			"total_users":        totalUsers,
			"total_jobs":         totalJobs,
			"total_applications": totalApplications,
			"active_jobs":        activeJobs,
		},
	})
}
