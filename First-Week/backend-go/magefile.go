//go:build mage

package main

import (
	"fmt"
	"os"

	"github.com/magefile/mage/sh"
)

// Build builds the application with CGO enabled
func Build() error {
	os.Setenv("CGO_ENABLED", "1")
	return sh.Run("go", "build", "-o", "bin/api", "./cmd/api")
}

// Run runs the application
func Run() error {
	return sh.Run("go", "run", "./cmd/api/main.go")
}

// Test runs the test suite
func Test() error {
	return sh.Run("go", "test", "./...")
}

// InitDB initializes the SQLite database
func InitDB() error {
	// Create database directory if it doesn't exist
	if err := os.MkdirAll("data", 0755); err != nil {
		return fmt.Errorf("failed to create data directory: %v", err)
	}

	// Touch the database file
	dbFile := "./data/words.db"
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		file, err := os.Create(dbFile)
		if err != nil {
			return fmt.Errorf("failed to create database file: %v", err)
		}
		file.Close()
	}

	return nil
}

// Clean removes the database file
func Clean() error {
	return os.RemoveAll("data")
} 