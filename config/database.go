package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	var err error
	dsn := os.Getenv("DATABASE_URL")
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	fmt.Println("✅ Database connected successfully")

	// Auto-migrate tables
	autoMigrate()
}

func CloseDB() {
	sqlDB, err := DB.DB()
	if err != nil {
		return
	}
	sqlDB.Close()
}

func autoMigrate() {
	err := DB.AutoMigrate(
		&User{},
		&Company{},
		&Job{},
		&Application{},
	)
	if err != nil {
		log.Printf("Auto-migration error: %v", err)
	} else {
		fmt.Println("✅ All tables migrated successfully")
	}
}

// User model - for job seekers, employers, and admins
type User struct {
	ID        uint      `gorm:"primaryKey"`
	Name      string    `gorm:"type:varchar(255);not null"`
	Email     string    `gorm:"type:varchar(255);unique;not null"`
	Password  string    `gorm:"type:varchar(255);not null"`
	Role      string    `gorm:"type:varchar(50);default:'job_seeker'"` // job_seeker, employer, admin
	Phone     string    `gorm:"type:varchar(20)"`
	Location  string    `gorm:"type:varchar(255)"`
	Bio       string    `gorm:"type:text"`
	ResumeURL string    `gorm:"type:text"` // Cloudinary URL for resume
	Verified  bool      `gorm:"default:false"`
	CompanyID *uint     `gorm:"index"` // For employers - links to Company
	Company   *Company  `gorm:"foreignKey:CompanyID"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

// Company model - for employer organizations
type Company struct {
	ID          uint      `gorm:"primaryKey"`
	Name        string    `gorm:"type:varchar(255);not null"`
	Description string    `gorm:"type:text"`
	Website     string    `gorm:"type:varchar(255)"`
	Location    string    `gorm:"type:varchar(255)"`
	LogoURL     string    `gorm:"type:text"` // Company logo
	Industry    string    `gorm:"type:varchar(100)"`
	Size        string    `gorm:"type:varchar(50)"` // e.g., "1-10", "11-50", "51-200"
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
}

// Job model - job postings
type Job struct {
	ID          uint       `gorm:"primaryKey"`
	Title       string     `gorm:"type:varchar(255);not null"`
	Description string     `gorm:"type:text;not null"`
	CompanyID   uint       `gorm:"not null;index"`
	Company     Company    `gorm:"foreignKey:CompanyID"`
	Location    string     `gorm:"type:varchar(255)"`
	JobType     string     `gorm:"type:varchar(50)"`                  // full-time, part-time, contract, remote
	Salary      string     `gorm:"type:varchar(100)"`                 // e.g., "$50k - $80k"
	Experience  string     `gorm:"type:varchar(100)"`                 // e.g., "2-5 years"
	Skills      string     `gorm:"type:text"`                         // Comma-separated or JSON
	Status      string     `gorm:"type:varchar(50);default:'active'"` // active, closed, pending
	PostedBy    uint       `gorm:"not null;index"`                    // User ID of employer
	Poster      User       `gorm:"foreignKey:PostedBy"`
	CreatedAt   time.Time  `gorm:"autoCreateTime"`
	UpdatedAt   time.Time  `gorm:"autoUpdateTime"`
	DeletedAt   *time.Time `gorm:"index"`
}

// Application model - job applications
type Application struct {
	ID          uint      `gorm:"primaryKey"`
	JobID       uint      `gorm:"not null;index"`
	Job         Job       `gorm:"foreignKey:JobID"`
	UserID      uint      `gorm:"not null;index"`
	User        User      `gorm:"foreignKey:UserID"`
	ResumeURL   string    `gorm:"type:text;not null"` // Cloudinary URL
	CoverLetter string    `gorm:"type:text"`
	Status      string    `gorm:"type:varchar(50);default:'pending'"` // pending, reviewed, shortlisted, rejected, accepted
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
}
