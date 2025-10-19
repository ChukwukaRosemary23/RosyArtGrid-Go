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

	// Public routes
	public := r.Group("/api")
	{
		// Browse projects
		public.GET("/projects", handlers.GetProjects)
		public.GET("/projects/:id", handlers.GetProject)

		// Categories
		public.GET("/categories", handlers.GetCategories)

		// Project likes and comments (public view)
		public.GET("/projects/:id/likes", handlers.GetProjectLikes)
		public.GET("/projects/:id/comments", handlers.GetProjectComments)

		// User followers/following (public view)
		public.GET("/users/:id/followers", handlers.GetUserFollowers)
		public.GET("/users/:id/following", handlers.GetUserFollowing)
	}

	// Protected routes (requires authentication)
	protected := r.Group("/api")
	protected.Use(middleware.AuthMiddleware())
	{
		// Profile
		protected.GET("/profile", handlers.GetProfile)
		protected.PUT("/profile", handlers.UpdateProfile)

		// Upload
		protected.POST("/upload/image", handlers.UploadImage)
		protected.POST("/upload/images", handlers.UploadMultipleImages)
		protected.POST("/upload/avatar", handlers.UploadAvatar)

		// Projects
		protected.POST("/projects", handlers.CreateProject)
		protected.GET("/my-projects", handlers.GetMyProjects)
		protected.PUT("/projects/:id", handlers.UpdateProject)
		protected.DELETE("/projects/:id", handlers.DeleteProject)

		// Social features
		protected.POST("/projects/:id/like", handlers.LikeProject)
		protected.DELETE("/projects/:id/unlike", handlers.UnlikeProject)
		protected.POST("/projects/:id/comments", handlers.AddComment)
		protected.DELETE("/comments/:id", handlers.DeleteComment)

		// Follow/Unfollow
		protected.POST("/users/:id/follow", handlers.FollowUser)
		protected.DELETE("/users/:id/unfollow", handlers.UnfollowUser)
	}
}
