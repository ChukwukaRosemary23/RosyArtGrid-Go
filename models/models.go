package models

import "time"

// Auth requests
type RegisterRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Role     string `json:"role"` // creative or company
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// Project requests
type CreateProjectRequest struct {
	Title       string   `json:"title" binding:"required"`
	Description string   `json:"description" binding:"required"`
	CategoryID  uint     `json:"category_id" binding:"required"`
	Tags        string   `json:"tags"`                          // Comma-separated
	ImageURLs   []string `json:"image_urls" binding:"required"` // Already uploaded to Cloudinary
}

type UpdateProjectRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Tags        string `json:"tags"`
}

// Comment request
type CreateCommentRequest struct {
	Content string `json:"content" binding:"required"`
}

// Company request
type CreateCompanyRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Website     string `json:"website"`
	Location    string `json:"location"`
	Industry    string `json:"industry"`
	Size        string `json:"size"`
}

// Response models
type UserResponse struct {
	ID        uint             `json:"id"`
	Name      string           `json:"name"`
	Email     string           `json:"email"`
	Role      string           `json:"role"`
	Bio       string           `json:"bio"`
	Location  string           `json:"location"`
	Skills    string           `json:"skills"`
	Website   string           `json:"website"`
	AvatarURL string           `json:"avatar_url"`
	ForHire   bool             `json:"for_hire"`
	Company   *CompanyResponse `json:"company,omitempty"`
	CreatedAt time.Time        `json:"created_at"`
}

type CompanyResponse struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Website     string `json:"website"`
	Location    string `json:"location"`
	LogoURL     string `json:"logo_url"`
	Industry    string `json:"industry"`
	Size        string `json:"size"`
}

type ProjectResponse struct {
	ID          uint                   `json:"id"`
	Title       string                 `json:"title"`
	Description string                 `json:"description"`
	CoverImage  string                 `json:"cover_image"`
	Images      []ProjectImageResponse `json:"images"`
	User        UserResponse           `json:"user"`
	Category    CategoryResponse       `json:"category"`
	Tags        string                 `json:"tags"`
	Views       int                    `json:"views"`
	LikesCount  int                    `json:"likes_count"`
	IsLiked     bool                   `json:"is_liked"` // If current user liked it
	CreatedAt   time.Time              `json:"created_at"`
}

type ProjectImageResponse struct {
	ID       uint   `json:"id"`
	ImageURL string `json:"image_url"`
	Order    int    `json:"order"`
}

type CategoryResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
	Icon string `json:"icon"`
}

type CommentResponse struct {
	ID        uint         `json:"id"`
	User      UserResponse `json:"user"`
	Content   string       `json:"content"`
	CreatedAt time.Time    `json:"created_at"`
}

type LikeResponse struct {
	ID        uint         `json:"id"`
	User      UserResponse `json:"user"`
	CreatedAt time.Time    `json:"created_at"`
}
