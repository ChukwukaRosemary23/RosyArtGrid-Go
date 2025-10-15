package routes

import (
	"jobconnect-backend/handlers"
	"jobconnect-backend/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	// Auth routes (public)
	auth := r.Group("/api/auth")
	{
		auth.POST("/register", handlers.RegisterUser)
		auth.POST("/login", handlers.LoginUser)
	}

	// Public job routes
	public := r.Group("/api")
	{
		public.GET("/jobs", handlers.GetJobs)
		public.GET("/jobs/:id", handlers.GetJob)
	}

	// Protected routes (job seekers and employers)
	protected := r.Group("/api")
	protected.Use(middleware.AuthMiddleware())
	{
		// Profile
		protected.GET("/profile", handlers.GetProfile)
		protected.PUT("/profile", handlers.UpdateProfile)

		// Resume upload (job seekers)
		protected.POST("/upload_resume", handlers.UploadResume)

		// Job applications (job seekers)
		protected.POST("/jobs/:id/apply", handlers.ApplyForJob)
		protected.GET("/my_applications", handlers.GetMyApplications)
	}

	// Employer routes
	employer := r.Group("/api/employer")
	employer.Use(middleware.AuthMiddleware(), middleware.EmployerMiddleware())
	{
		// Company management
		employer.POST("/company", handlers.CreateCompany)
		employer.GET("/company", handlers.GetMyCompany)
		employer.PUT("/company", handlers.UpdateCompany)
		employer.POST("/company/logo", handlers.UploadCompanyLogo)

		// Job management
		employer.POST("/jobs", handlers.CreateJob)
		employer.GET("/jobs", handlers.GetMyJobs)
		employer.PUT("/jobs/:id", handlers.UpdateJob)
		employer.DELETE("/jobs/:id", handlers.DeleteJob)

		// Application management
		employer.GET("/jobs/:id/applications", handlers.GetJobApplications)
		employer.PUT("/applications/:id/status", handlers.UpdateApplicationStatus)
	}

	// Admin routes
	admin := r.Group("/api/admin")
	admin.Use(middleware.AuthMiddleware(), middleware.AdminMiddleware())
	{
		// User management
		admin.GET("/users", handlers.GetAllUsers)
		admin.DELETE("/users/:id", handlers.DeleteUser)

		// Job management
		admin.GET("/jobs", handlers.GetAllJobsAdmin)
		admin.PUT("/jobs/:id/status", handlers.UpdateJobStatus)

		// Application management
		admin.GET("/applications", handlers.GetAllApplications)

		// Dashboard stats
		admin.GET("/stats", handlers.GetDashboardStats)
	}
}
