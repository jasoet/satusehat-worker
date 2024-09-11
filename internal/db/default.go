package db

import (
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "modernc.org/sqlite"
	"sync"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS
var (
	instance *Repository
	once     sync.Once
)

func DefaultRepository(filePath ...string) (*Repository, error) {
	var err error
	once.Do(func() {
		config, ierr := defaultInternalConfig(filePath...)
		if ierr != nil {
			err = ierr
			return
		}

		pool, ierr := config.Pool()
		if ierr != nil {
			err = ierr
			return
		}

		// Run migrations
		if ierr = runMigrations(pool.DB); ierr != nil {
			err = ierr
			return
		}

		instance, err = newRepository(pool)
	})

	return instance, err

}

func runMigrations(db *sql.DB) error {
	driver, err := sqlite.WithInstance(db, &sqlite.Config{})
	if err != nil {
		return fmt.Errorf("failed to create database driver: %w", err)
	}

	d, err := iofs.New(migrationsFS, "migrations")
	if err != nil {
		return fmt.Errorf("failed to create migration source: %w", err)
	}

	m, err := migrate.NewWithInstance("iofs", d, "sqlite", driver)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	return nil
}
