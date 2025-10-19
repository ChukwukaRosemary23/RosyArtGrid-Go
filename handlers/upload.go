package handlers

import (
	"context"
	"jobconnect-backend/config"
	"net/http"
	"os"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/gin-gonic/gin"
)

// UploadImage - Upload single image to Cloudinary
func UploadImage(c *gin.Context) {
	// Get uploaded file
	file, err := c.FormFile("image")
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
		Folder: "rosyartgrid/projects",
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload image"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":   true,
		"image_url": uploadResult.SecureURL,
	})
}

// UploadMultipleImages - Upload multiple images
func UploadMultipleImages(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse form"})
		return
	}

	files := form.File["images"]
	if len(files) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No files uploaded"})
		return
	}

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

	var imageURLs []string
	ctx := context.Background()

	// Upload each file
	for _, file := range files {
		fileContent, err := file.Open()
		if err != nil {
			continue
		}

		uploadResult, err := cld.Upload.Upload(ctx, fileContent, uploader.UploadParams{
			Folder: "rosyartgrid/projects",
		})
		fileContent.Close()

		if err == nil {
			imageURLs = append(imageURLs, uploadResult.SecureURL)
		}
	}

	if len(imageURLs) == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload images"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":    true,
		"image_urls": imageURLs,
	})
}

// UploadAvatar - Upload user avatar
func UploadAvatar(c *gin.Context) {
	userID, _ := c.Get("user_id")

	// Get uploaded file
	file, err := c.FormFile("avatar")
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
		Folder: "rosyartgrid/avatars",
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload avatar"})
		return
	}

	// Update user avatar
	config.DB.Model(&config.User{}).Where("id = ?", userID).Update("avatar_url", uploadResult.SecureURL)

	c.JSON(http.StatusOK, gin.H{
		"success":    true,
		"avatar_url": uploadResult.SecureURL,
	})
}
