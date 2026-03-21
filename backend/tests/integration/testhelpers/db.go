package testhelpers

import (
	"fmt"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/luisfucros/expense-tracker-app/internal/domain/model"
)

// SetupTestDB connects to the test database using environment variables,
// creates the database if it does not exist, and auto-migrates the schema.
// Returns an error (callers should skip the test) if the DB is unavailable.
func SetupTestDB() (*gorm.DB, error) {
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "3307")
	user := getEnv("DB_USER", "root")
	password := getEnv("DB_PASSWORD", "secret")
	dbName := getEnv("DB_NAME", "expense_tracker_test")

	// Connect without specifying the database to create it if needed.
	rootDSN := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/?charset=utf8mb4&parseTime=True&loc=UTC",
		user, password, host, port,
	)
	rootDB, err := gorm.Open(mysql.Open(rootDSN), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("could not connect to test DB: %w", err)
	}
	if err := rootDB.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s`", dbName)).Error; err != nil {
		return nil, fmt.Errorf("could not create test database: %w", err)
	}
	sqlDB, _ := rootDB.DB()
	sqlDB.Close()

	// Now connect to the test database.
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=UTC",
		user, password, host, port, dbName,
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("could not connect to test DB: %w", err)
	}

	if err := db.AutoMigrate(&model.User{}, &model.Expense{}); err != nil {
		return nil, fmt.Errorf("auto-migrate failed: %w", err)
	}

	return db, nil
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
