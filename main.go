package main

import (
	"fmt"
	"log"
	"os"

	"jobconnect-backend/config"
	"jobconnect-backend/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

func createDefaultAdmin() {
	adminEmail := os.Getenv("ADMIN_EMAIL")
	
	// Check if admin already exists
	var count int64
	config.DB.Model(&config.User{}).Where("email = ?", adminEmail).Count(&count)

	if count > 0 {
		fmt.Println("Admin user already exists")
		return
	}

	// Hash the password from environment variable
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(os.Getenv("ADMIN_PASSWORD")), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Warning: Failed to create admin user:", err)
		return
	}

	// Create admin user
	admin := config.User{
		Name:     "Admin User",
		Email:    adminEmail,
		Password: string(hashedPassword),
		Role:     "admin",
		Verified: true,
	}

	if err := config.DB.Create(&admin).Error; err != nil {
		log.Println("Warning: Failed to create admin user:", err)
		return
	}

	fmt.Println("✅ Admin user created successfully!")
}

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Initialize database
	config.InitDB()
	defer config.CloseDB()

	// Auto-create admin user if it doesn't exist
	createDefaultAdmin()

	// Setup Gin router
	r := gin.Default()

	// CORS configuration
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{
		os.Getenv("FRONTEND_URL"),
		"http://localhost:3000",
		"http://localhost:5173",
	}
	corsConfig.AllowCredentials = true
	corsConfig.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	r.Use(cors.New(corsConfig))

	// Setup routes
	routes.SetupRoutes(r)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}

	fmt.Printf("🚀 Server starting on port %s\n", port)
	log.Fatal(r.Run(":" + port))
}