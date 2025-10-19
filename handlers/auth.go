package handlers

import (
	"net/http"

	"jobconnect-backend/config"
	"jobconnect-backend/models"
	"jobconnect-backend/utils"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func RegisterUser(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate role - changed from job_seeker/employer to creative/company
	if req.Role != "creative" && req.Role != "company" {
		req.Role = "creative" // Default to creative
	}

	// Check if email already exists
	var existingUser config.User
	if err := config.DB.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email already registered"})
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// Create user
	user := config.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashedPassword),
		Role:     req.Role,
		Verified: true,
	}

	if err := config.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	// Generate JWT token
	token, err := utils.GenerateToken(user.ID, user.Email, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Registration successful",
		"token":   token,
		"user": models.UserResponse{
			ID:       user.ID,
			Name:     user.Name,
			Email:    user.Email,
			Role:     user.Role,
			Location: user.Location,
		},
	})
}

func LoginUser(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Find user
	var user config.User
	if err := config.DB.Preload("Company").Where("email = ?", req.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Check password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Generate JWT token
	token, err := utils.GenerateToken(user.ID, user.Email, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// Prepare response
	userResponse := models.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      user.Role,
		Location:  user.Location,
		Bio:       user.Bio,
		CreatedAt: user.CreatedAt,
	}

	if user.Company != nil {
		userResponse.Company = &models.CompanyResponse{
			ID:          user.Company.ID,
			Name:        user.Company.Name,
			Description: user.Company.Description,
			Website:     user.Company.Website,
			Location:    user.Company.Location,
			LogoURL:     user.Company.LogoURL,
			Industry:    user.Company.Industry,
			Size:        user.Company.Size,
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Login successful",
		"token":   token,
		"user":    userResponse,
	})
}

func GetProfile(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var user config.User
	if err := config.DB.Preload("Company").Where("id = ?", userID).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	userResponse := models.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      user.Role,
		Location:  user.Location,
		Bio:       user.Bio,
		CreatedAt: user.CreatedAt,
	}

	if user.Company != nil {
		userResponse.Company = &models.CompanyResponse{
			ID:          user.Company.ID,
			Name:        user.Company.Name,
			Description: user.Company.Description,
			Website:     user.Company.Website,
			Location:    user.Company.Location,
			LogoURL:     user.Company.LogoURL,
			Industry:    user.Company.Industry,
			Size:        user.Company.Size,
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"user":    userResponse,
	})
}

func UpdateProfile(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req struct {
		Name     string `json:"name"`
		Location string `json:"location"`
		Bio      string `json:"bio"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update user
	result := config.DB.Model(&config.User{}).Where("id = ?", userID).Updates(map[string]interface{}{
		"name":     req.Name,
		"location": req.Location,
		"bio":      req.Bio,
	})

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Profile updated successfully",
	})
}
