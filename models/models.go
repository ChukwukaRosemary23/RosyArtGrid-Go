package models

import "time"

// Auth requests
type RegisterRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Role     string `json:"role"` // job_seeker or employer
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// Job requests
type CreateJobRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description" binding:"required"`
	Location    string `json:"location"`
	JobType     string `json:"job_type"` // full-time, part-time, contract, remote
	Salary      string `json:"salary"`
	Experience  string `json:"experience"`
	Skills      string `json:"skills"`
}

type UpdateJobRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Location    string `json:"location"`
	JobType     string `json:"job_type"`
	Salary      string `json:"salary"`
	Experience  string `json:"experience"`
	Skills      string `json:"skills"`
	Status      string `json:"status"` // active, closed
}

// Application requests
type ApplyJobRequest struct {
	CoverLetter string `json:"cover_letter"`
}

// Company requests
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
	Phone     string           `json:"phone"`
	Location  string           `json:"location"`
	Bio       string           `json:"bio"`
	ResumeURL string           `json:"resume_url"`
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

type JobResponse struct {
	ID          uint            `json:"id"`
	Title       string          `json:"title"`
	Description string          `json:"description"`
	Company     CompanyResponse `json:"company"`
	Location    string          `json:"location"`
	JobType     string          `json:"job_type"`
	Salary      string          `json:"salary"`
	Experience  string          `json:"experience"`
	Skills      string          `json:"skills"`
	Status      string          `json:"status"`
	CreatedAt   time.Time       `json:"created_at"`
}

type ApplicationResponse struct {
	ID          uint         `json:"id"`
	Job         JobResponse  `json:"job"`
	User        UserResponse `json:"user"`
	ResumeURL   string       `json:"resume_url"`
	CoverLetter string       `json:"cover_letter"`
	Status      string       `json:"status"`
	CreatedAt   time.Time    `json:"created_at"`
}
