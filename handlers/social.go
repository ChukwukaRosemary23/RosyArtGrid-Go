package handlers

import (
	"net/http"
	"strconv"

	"jobconnect-backend/config"
	"jobconnect-backend/models"

	"github.com/gin-gonic/gin"
)

// LikeProject - User likes a project
func LikeProject(c *gin.Context) {
	userID, _ := c.Get("user_id")
	projectID := c.Param("id")

	// Check if project exists
	var project config.Project
	if err := config.DB.Where("id = ?", projectID).First(&project).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		return
	}

	// Check if already liked
	var existingLike config.Like
	if err := config.DB.Where("user_id = ? AND project_id = ?", userID, projectID).First(&existingLike).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Already liked"})
		return
	}

	// Create like
	like := config.Like{
		UserID:    userID.(uint),
		ProjectID: project.ID,
	}

	if err := config.DB.Create(&like).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to like project"})
		return
	}

	// Increment likes count
	config.DB.Model(&project).Update("likes_count", project.LikesCount+1)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Project liked",
	})
}

// UnlikeProject - User unlikes a project
func UnlikeProject(c *gin.Context) {
	userID, _ := c.Get("user_id")
	projectID := c.Param("id")

	// Find and delete like
	var like config.Like
	if err := config.DB.Where("user_id = ? AND project_id = ?", userID, projectID).First(&like).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Like not found"})
		return
	}

	config.DB.Delete(&like)

	// Decrement likes count
	var project config.Project
	if err := config.DB.Where("id = ?", projectID).First(&project).Error; err == nil {
		if project.LikesCount > 0 {
			config.DB.Model(&project).Update("likes_count", project.LikesCount-1)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Project unliked",
	})
}

// GetProjectLikes - Get all likes for a project
func GetProjectLikes(c *gin.Context) {
	projectID := c.Param("id")

	var likes []config.Like
	config.DB.Preload("User").Where("project_id = ?", projectID).Order("created_at DESC").Find(&likes)

	var response []models.LikeResponse
	for _, like := range likes {
		response = append(response, models.LikeResponse{
			ID: like.ID,
			User: models.UserResponse{
				ID:        like.User.ID,
				Name:      like.User.Name,
				AvatarURL: like.User.AvatarURL,
			},
			CreatedAt: like.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"likes":   response,
	})
}

// AddComment - User comments on a project
func AddComment(c *gin.Context) {
	userID, _ := c.Get("user_id")
	projectID := c.Param("id")

	var req models.CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if project exists
	var project config.Project
	if err := config.DB.Where("id = ?", projectID).First(&project).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		return
	}

	// Create comment
	comment := config.Comment{
		UserID:    userID.(uint),
		ProjectID: project.ID,
		Content:   req.Content,
	}

	if err := config.DB.Create(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add comment"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success":    true,
		"message":    "Comment added",
		"comment_id": comment.ID,
	})
}

// GetProjectComments - Get all comments for a project
func GetProjectComments(c *gin.Context) {
	projectID := c.Param("id")

	var comments []config.Comment
	config.DB.Preload("User").Where("project_id = ?", projectID).Order("created_at DESC").Find(&comments)

	var response []models.CommentResponse
	for _, comment := range comments {
		response = append(response, models.CommentResponse{
			ID: comment.ID,
			User: models.UserResponse{
				ID:        comment.User.ID,
				Name:      comment.User.Name,
				AvatarURL: comment.User.AvatarURL,
			},
			Content:   comment.Content,
			CreatedAt: comment.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"success":  true,
		"comments": response,
	})
}

// DeleteComment - User deletes their comment
func DeleteComment(c *gin.Context) {
	userID, _ := c.Get("user_id")
	commentID := c.Param("id")

	result := config.DB.Where("id = ? AND user_id = ?", commentID, userID).Delete(&config.Comment{})

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Comment not found or unauthorized"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Comment deleted",
	})
}

// FollowUser - Follow a user
func FollowUser(c *gin.Context) {
	followerID, _ := c.Get("user_id")
	followingID := c.Param("id")

	// Can't follow yourself
	if strconv.Itoa(int(followerID.(uint))) == followingID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot follow yourself"})
		return
	}

	// Check if user exists
	var user config.User
	if err := config.DB.Where("id = ?", followingID).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Check if already following
	var existingFollow config.Follow
	if err := config.DB.Where("follower_id = ? AND following_id = ?", followerID, followingID).First(&existingFollow).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Already following"})
		return
	}

	// Create follow
	followingIDUint, _ := strconv.ParseUint(followingID, 10, 32)
	follow := config.Follow{
		FollowerID:  followerID.(uint),
		FollowingID: uint(followingIDUint),
	}

	if err := config.DB.Create(&follow).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to follow user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "User followed",
	})
}

// UnfollowUser - Unfollow a user
func UnfollowUser(c *gin.Context) {
	followerID, _ := c.Get("user_id")
	followingID := c.Param("id")

	result := config.DB.Where("follower_id = ? AND following_id = ?", followerID, followingID).Delete(&config.Follow{})

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not following this user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "User unfollowed",
	})
}

// GetUserFollowers - Get user's followers
func GetUserFollowers(c *gin.Context) {
	userID := c.Param("id")

	var follows []config.Follow
	config.DB.Preload("Follower").Where("following_id = ?", userID).Find(&follows)

	var response []models.UserResponse
	for _, follow := range follows {
		response = append(response, models.UserResponse{
			ID:        follow.Follower.ID,
			Name:      follow.Follower.Name,
			AvatarURL: follow.Follower.AvatarURL,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"success":   true,
		"followers": response,
	})
}

// GetUserFollowing - Get users that this user follows
func GetUserFollowing(c *gin.Context) {
	userID := c.Param("id")

	var follows []config.Follow
	config.DB.Preload("Following").Where("follower_id = ?", userID).Find(&follows)

	var response []models.UserResponse
	for _, follow := range follows {
		response = append(response, models.UserResponse{
			ID:        follow.Following.ID,
			Name:      follow.Following.Name,
			AvatarURL: follow.Following.AvatarURL,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"success":   true,
		"following": response,
	})
}
