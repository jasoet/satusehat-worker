package db

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"os"
	"path/filepath"
	"time"

	_ "modernc.org/sqlite"
)

// InternalConfig holds configuration options for SQLite
type InternalConfig struct {
	FilePath          string
	CacheSize         int    // in pages
	JournalMode       string // DELETE, TRUNCATE, PERSIST, MEMORY, WAL, OFF
	Synchronous       string // EXTRA, FULL, NORMAL, OFF
	ForeignKeys       bool
	RecursiveTriggers bool
	BusyTimeout       time.Duration
	MaxOpenConns      int
	MaxIdleConns      int
	ConnMaxLifetime   time.Duration
}

// defaultInternalConfig returns an InternalConfig with default settings
func defaultInternalConfig(pathElements ...string) (InternalConfig, error) {
	var filePath string

	if len(pathElements) > 0 {
		// Join the provided path elements
		filePath = filepath.Join(pathElements...)
	} else {
		// Use default path in user's home directory
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return InternalConfig{}, fmt.Errorf("failed to get user home directory: %w", err)
		}
		filePath = filepath.Join(homeDir, "internal.db")
	}

	// Ensure the directory exists
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return InternalConfig{}, fmt.Errorf("failed to create directory %s: %w", dir, err)
	}

	return InternalConfig{
		FilePath:          filePath,
		CacheSize:         -2000,    // 2000 pages, negative means kibibytes
		JournalMode:       "WAL",    // Write-Ahead Logging for better concurrency
		Synchronous:       "NORMAL", // Good balance between safety and speed
		ForeignKeys:       true,
		RecursiveTriggers: true,
		BusyTimeout:       10 * time.Second, // Increased for better external access tolerance
		MaxOpenConns:      5,                // Increased to allow more concurrent access
		MaxIdleConns:      5,                // Matching MaxOpenConns
		ConnMaxLifetime:   1 * time.Hour,
	}, nil
}

// Dsn returns the data source name (connection string) for the SQLite database
func (c InternalConfig) Dsn() string {
	return fmt.Sprintf("%s?_cache_size=%d&_journal_mode=%s&_synchronous=%s&_foreign_keys=%t&_recursive_triggers=%t&_busy_timeout=%d",
		c.FilePath, c.CacheSize, c.JournalMode, c.Synchronous, c.ForeignKeys, c.RecursiveTriggers, c.BusyTimeout.Milliseconds())
}

func (c InternalConfig) Pool() (*sqlx.DB, error) {
	db, err := sqlx.Open("sqlite", c.Dsn())
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(c.MaxOpenConns)
	db.SetMaxIdleConns(c.MaxIdleConns)
	db.SetConnMaxLifetime(c.ConnMaxLifetime)

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
