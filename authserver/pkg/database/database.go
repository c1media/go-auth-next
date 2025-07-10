package database

import (
	"fmt"
	"log"
	"os"

	"github.com/simple-auth-roles/internal/types"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Connect establishes a connection to the database
func Connect(databaseURL string) (*gorm.DB, error) {
	if databaseURL == "" {
		return nil, fmt.Errorf("database URL is required")
	}

	db, err := gorm.Open(postgres.Open(databaseURL), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return db, nil
}

// RunMigrations runs the database migrations safely
func RunMigrations(db *gorm.DB) error {
	logger := log.New(os.Stdout, "[MIGRATIONS] ", log.LstdFlags)
	logger.Println("Starting database migrations...")

	// Check if migrations have already run (safe for Railway deployments)
	if db.Migrator().HasTable(&types.User{}) {
		logger.Println("Database tables already exist, checking for schema updates...")
	} else {
		logger.Println("Creating database tables for the first time...")
	}

	// Run auto-migrations (safe - only adds new columns/tables)
	err := db.AutoMigrate(
		&types.User{},
		&types.WebAuthnCredential{},
	)
	
	if err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	logger.Println("Database migrations completed successfully")
	return nil
}

// SeedAdminUser creates an admin user if none exists
func SeedAdminUser(db *gorm.DB) error {
	logger := log.New(os.Stdout, "[SEED] ", log.LstdFlags)
	
	// Check if any admin users exist
	var adminCount int64
	err := db.Model(&types.User{}).Where("role = ?", "admin").Count(&adminCount).Error
	if err != nil {
		return fmt.Errorf("failed to check for admin users: %w", err)
	}

	if adminCount > 0 {
		logger.Printf("Admin user already exists (count: %d), skipping seed", adminCount)
		return nil
	}

	// Get admin email from environment
	adminEmail := os.Getenv("ADMIN_EMAIL")
	if adminEmail == "" {
		adminEmail = "admin@example.com" // Default fallback
		logger.Printf("ADMIN_EMAIL not set, using default: %s", adminEmail)
	}

	adminName := os.Getenv("ADMIN_NAME")
	if adminName == "" {
		adminName = "System Administrator"
	}

	// Create admin user
	admin := types.User{
		Email: adminEmail,
		Name:  adminName,
		Role:  "admin",
	}

	err = db.Create(&admin).Error
	if err != nil {
		return fmt.Errorf("failed to create admin user: %w", err)
	}

	logger.Printf("✅ Admin user created successfully: %s (%s)", admin.Name, admin.Email)
	logger.Println("💡 Use this email to sign in as admin on first deployment")
	
	return nil
}
