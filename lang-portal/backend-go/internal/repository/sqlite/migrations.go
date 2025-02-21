package sqlite

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// RunMigrations executes all SQL migration files in the specified directory
func (db *Database) RunMigrations(migrationsDir string) error {
	// Read migration files
	files, err := os.ReadDir(migrationsDir)
	if err != nil {
		return fmt.Errorf("error reading migrations directory: %v", err)
	}

	// Get only .sql files and sort them
	var migrations []string
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".sql") {
			migrations = append(migrations, file.Name())
		}
	}
	sort.Strings(migrations)

	// Execute each migration file
	for _, migration := range migrations {
		migrationPath := filepath.Join(migrationsDir, migration)
		content, err := os.ReadFile(migrationPath)
		if err != nil {
			return fmt.Errorf("error reading migration file %s: %v", migration, err)
		}

		// Execute the migration within a transaction
		tx, err := db.DB.Begin()
		if err != nil {
			return fmt.Errorf("error starting transaction for %s: %v", migration, err)
		}

		if _, err := tx.Exec(string(content)); err != nil {
			tx.Rollback()
			return fmt.Errorf("error executing migration %s: %v", migration, err)
		}

		if err := tx.Commit(); err != nil {
			return fmt.Errorf("error committing migration %s: %v", migration, err)
		}
	}

	return nil
} 