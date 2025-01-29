package models

import (
	"log"
	"os"
	"path/filepath"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	// Ensure the "db" directory exists
	dbPath := "db/wiki.db"
	dbDir := filepath.Dir(dbPath)

	if err := os.MkdirAll(dbDir, os.ModePerm); err != nil {
		log.Fatalf("Failed to create database directory: %v", err)
	}

	// Open the SQLite database
	var err error
	DB, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{})

	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Run AutoMigrate for all tables
	err = DB.AutoMigrate(&User{}, &Role{}, &UserRole{}, &Page{})
	if err != nil {
		log.Fatal("Migration failed:", err)
	}

	log.Println("Database migration completed successfully.")

	// Ensure "Home" page exists
	var homePage Page
	result := DB.Where("title = ?", "Home").First(&homePage)

	if result.Error != nil && result.Error == gorm.ErrRecordNotFound {
		log.Println("Creating default Home page...")
		homePage = Page{
			Title:   "Home",
			Content: "Welcome to your wiki! Click 'Edit' to add content.",
		}

		if err := DB.Create(&homePage).Error; err != nil {
			log.Fatalf("Failed to create Home page: %v", err) // Logs the actual error if insertion fails
		} else {
			log.Println("Home page created successfully!")
		}
	}
}
