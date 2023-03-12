package database

import (
	"database/sql"
	"fmt"
	"os"

	"app/config"
)

func Connect() (*sql.DB, error) {
	// Get database connection string from environment variable
	config.Load()
	connectionString := os.Getenv("DATABASE_URL")
	if connectionString == "" {
		return nil, fmt.Errorf("DATABASE_URL environment variable not set")
	}

	// Open database connection
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %v", err)
	}

	// Test database connection
	if err = db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	fmt.Println("Successfully connected!")
	return db, nil
}

