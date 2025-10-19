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
		&Project{},
		&ProjectImage{},
		&Like{},
		&Comment{},
		&Follow{},
		&Category{},
	)
	if err != nil {
		log.Printf("Auto-migration error: %v", err)
	} else {
		fmt.Println("✅ All tables migrated successfully")
	}
}

// User model - for creatives, companies, and admins
type User struct {
	ID          uint      `gorm:"primaryKey"`
	Name        string    `gorm:"type:varchar(255);not null"`
	Email       string    `gorm:"type:varchar(255);unique;not null"`
	Password    string    `gorm:"type:varchar(255);not null"`
	Role        string    `gorm:"type:varchar(50);default:'creative'"` // creative, company, admin
	Bio         string    `gorm:"type:text"`
	Location    string    `gorm:"type:varchar(255)"`
	Skills      string    `gorm:"type:text"` // Comma-separated
	Website     string    `gorm:"type:varchar(255)"`
	BehanceURL  string    `gorm:"type:varchar(255)"`
	DribbbleURL string    `gorm:"type:varchar(255)"`
	LinkedInURL string    `gorm:"type:varchar(255)"`
	TwitterURL  string    `gorm:"type:varchar(255)"`
	AvatarURL   string    `gorm:"type:text"` // Profile picture
	ForHire     bool      `gorm:"default:false"`
	Verified    bool      `gorm:"default:false"`
	CompanyID   *uint     `gorm:"index"` // For company accounts
	Company     *Company  `gorm:"foreignKey:CompanyID"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
}

// Company model - for recruiters/organizations
type Company struct {
	ID          uint      `gorm:"primaryKey"`
	Name        string    `gorm:"type:varchar(255);not null"`
	Description string    `gorm:"type:text"`
	Website     string    `gorm:"type:varchar(255)"`
	Location    string    `gorm:"type:varchar(255)"`
	LogoURL     string    `gorm:"type:text"`
	Industry    string    `gorm:"type:varchar(100)"`
	Size        string    `gorm:"type:varchar(50)"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
}

// Project model - creative work showcase (replaces Job)
type Project struct {
	ID          uint           `gorm:"primaryKey"`
	Title       string         `gorm:"type:varchar(255);not null"`
	Description string         `gorm:"type:text;not null"`
	UserID      uint           `gorm:"not null;index"`
	User        User           `gorm:"foreignKey:UserID"`
	CategoryID  uint           `gorm:"not null;index"`
	Category    Category       `gorm:"foreignKey:CategoryID"`
	Tags        string         `gorm:"type:text"` // Comma-separated
	CoverImage  string         `gorm:"type:text"` // Main cover image URL
	Images      []ProjectImage `gorm:"foreignKey:ProjectID"`
	Views       int            `gorm:"default:0"`
	LikesCount  int            `gorm:"default:0"`
	Featured    bool           `gorm:"default:false"` // For admin to feature projects
	CreatedAt   time.Time      `gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime"`
	DeletedAt   *time.Time     `gorm:"index"`
}

// ProjectImage model - multiple images per project
type ProjectImage struct {
	ID        uint      `gorm:"primaryKey"`
	ProjectID uint      `gorm:"not null;index"`
	ImageURL  string    `gorm:"type:text;not null"`
	Order     int       `gorm:"default:0"` // For ordering images
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

// Category model - creative categories
type Category struct {
	ID        uint      `gorm:"primaryKey"`
	Name      string    `gorm:"type:varchar(100);unique;not null"` // e.g., Graphic Design, Photography
	Slug      string    `gorm:"type:varchar(100);unique;not null"` // URL-friendly name
	Icon      string    `gorm:"type:varchar(50)"`                  // Emoji or icon class
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

// Like model - users like projects
type Like struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    uint      `gorm:"not null;index"`
	User      User      `gorm:"foreignKey:UserID"`
	ProjectID uint      `gorm:"not null;index"`
	Project   Project   `gorm:"foreignKey:ProjectID"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

// Comment model - users comment on projects
type Comment struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    uint      `gorm:"not null;index"`
	User      User      `gorm:"foreignKey:UserID"`
	ProjectID uint      `gorm:"not null;index"`
	Project   Project   `gorm:"foreignKey:ProjectID"`
	Content   string    `gorm:"type:text;not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

// Follow model - users follow other users
type Follow struct {
	ID          uint      `gorm:"primaryKey"`
	FollowerID  uint      `gorm:"not null;index"` // User who follows
	FollowingID uint      `gorm:"not null;index"` // User being followed
	Follower    User      `gorm:"foreignKey:FollowerID"`
	Following   User      `gorm:"foreignKey:FollowingID"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
}
