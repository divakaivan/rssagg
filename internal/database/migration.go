package database

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/pressly/goose"
)

func MigrateUp(db *sql.DB, path string) error {
	currentVersion, err := goose.GetDBVersion(db)
	if err != nil {
		return fmt.Errorf("migrate: could not get current version: %w", err)
	}

	if err = goose.Up(db, path); err != nil && !errors.Is(err, goose.ErrNoNextVersion) {
		return fmt.Errorf("migrate: failed to migrate up to the latest version: %w", err)
	}

	newVersion, err := goose.GetDBVersion(db)
	if err != nil {
		return fmt.Errorf("migrate: could not get new version: %w", err)
	}

	log.Printf("Migrated up from version %d to %d\n", currentVersion, newVersion)

	return nil
}
