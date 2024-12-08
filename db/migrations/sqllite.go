package migrations

import (
	"database/sql"
	"embed"
	"fmt"
	"log"
	"os"
	"path"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "modernc.org/sqlite"
)

//go:embed *
var migrationsFiles embed.FS

func checkDbFile(dbPath string) error {
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		dirPath := path.Dir(dbPath)
		os.MkdirAll(dirPath, os.ModePerm)
		file, err := os.Create(dbPath)
		if err != nil {
			return fmt.Errorf("failed to create SQLite database file: %w", err)
		}
		file.Close()
	}
	return nil
}

func RunSQLLite(dbPath string) error {
	err := checkDbFile(dbPath)
	if err != nil {
		return fmt.Errorf("failed to check SQLite database file: %w", err)
	}

	// Open SQLite database
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return fmt.Errorf("failed to open SQLite database: %w", err)
	}
	defer db.Close()

	// Initialize SQLite driver
	driver, err := sqlite.WithInstance(db, &sqlite.Config{})
	if err != nil {
		return fmt.Errorf("failed to create SQLite driver: %w", err)
	}

	// Create source driver for embedded migrations
	d, err := iofs.New(migrationsFiles, ".")
	if err != nil {
		return fmt.Errorf("failed to create source driver: %w", err)
	}

	// Run migrations
	m, err := migrate.NewWithInstance("iofs", d, "sqlite", driver)
	if err != nil {
		return fmt.Errorf("failed to initialize migrate: %w", err)
	}
	defer func() {
		if sourceErr, dbErr := m.Close(); sourceErr != nil || dbErr != nil {
			log.Printf("Error closing migrate: %v, %v", sourceErr, dbErr)
		}
	}()

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("migration failed: %w", err)
	}

	return nil
}
