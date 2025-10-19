package handlers

import (
	"context"
	"net/http"
	"os"

	"jobconnect-backend/config"
	"jobconnect-backend/models"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/gin-gonic/gin"
)

// CreateCompany - Employer creates company profile
func CreateCompany(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req models.CreateCompanyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if user already has a company
	var user config.User
	if err := config.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if user.CompanyID != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "You already have a company profile"})
		return
	}

	// Create company
	company := config.Company{
		Name:        req.Name,
		Description: req.Description,
		Website:     req.Website,
		Location:    req.Location,
		Industry:    req.Industry,
		Size:        req.Size,
	}

	if err := config.DB.Create(&company).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create company"})
		return
	}

	// Link company to user
	config.DB.Model(&user).Update("company_id", company.ID)

	c.JSON(http.StatusCreated, gin.H{
		"success":    true,
		"message":    "Company created successfully",
		"company_id": company.ID,
	})
}

// GetMyCompany - Employer gets their company
func GetMyCompany(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var user config.User
	if err := config.DB.Preload("Company").Where("id = ?", userID).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if user.Company == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No company profile found"})
		return
	}

	response := models.CompanyResponse{
		ID:          user.Company.ID,
		Name:        user.Company.Name,
		Description: user.Company.Description,
		Website:     user.Company.Website,
		Location:    user.Company.Location,
		LogoURL:     user.Company.LogoURL,
		Industry:    user.Company.Industry,
		Size:        user.Company.Size,
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"company": response,
	})
}

// UpdateCompany - Employer updates their company
func UpdateCompany(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req models.CreateCompanyRequest
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
		c.JSON(http.StatusNotFound, gin.H{"error": "No company profile found"})
		return
	}

	// Update company
	updates := map[string]interface{}{
		"name":        req.Name,
		"description": req.Description,
		"website":     req.Website,
		"location":    req.Location,
		"industry":    req.Industry,
		"size":        req.Size,
	}

	config.DB.Model(&config.Company{}).Where("id = ?", user.CompanyID).Updates(updates)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Company updated successfully",
	})
}

// UploadCompanyLogo - Upload company logo to Cloudinary
func UploadCompanyLogo(c *gin.Context) {
	userID, _ := c.Get("user_id")

	// Get user's company
	var user config.User
	if err := config.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if user.CompanyID == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No company profile found"})
		return
	}

	// Get uploaded file
	file, err := c.FormFile("logo")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
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
		Folder: "jobconnect/logos",
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload logo"})
		return
	}

	// Update company logo URL
	config.DB.Model(&config.Company{}).Where("id = ?", user.CompanyID).Update("logo_url", uploadResult.SecureURL)

	c.JSON(http.StatusOK, gin.H{
		"success":  true,
		"message":  "Logo uploaded successfully",
		"logo_url": uploadResult.SecureURL,
	})
}
