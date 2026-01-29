package db

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"strings"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	pool          *pgxpool.Pool
	once          sync.Once
	currentDBVersion string
	mu            sync.Mutex
)

// InitDB initializes the database connection pool
func InitDB() error {
	mu.Lock()
	defer mu.Unlock()

	// Check if DB_VERSION changed - if so, recreate the pool
	dbVersion := os.Getenv("DB_VERSION")
	if pool != nil && dbVersion != "" && dbVersion != currentDBVersion {
		pool.Close()
		pool = nil
		once = sync.Once{} // Reset once
	}

	var err error
	once.Do(func() {
		supabaseURL := os.Getenv("SUPABASE_URL")
		password := os.Getenv("DB_PASSWORD")
		
		if supabaseURL == "" {
			err = fmt.Errorf("SUPABASE_URL not set")
			return
		}
		if password == "" {
			err = fmt.Errorf("DB_PASSWORD not set")
			return
		}

		// Extract project reference from Supabase URL
		projectRef := extractProjectRef(supabaseURL)
		if projectRef == "" {
			err = fmt.Errorf("invalid SUPABASE_URL format: %s", supabaseURL)
			return
		}

		// URL encode the password to handle special characters
		encodedPassword := url.QueryEscape(password)

		// Build connection string using Supabase Pooler (port 6543)
		// IMPORTANT: statement_cache_mode=describe disables prepared statements
		// This is required for Supabase Pooler in Transaction Mode
		databaseURL := fmt.Sprintf(
			"postgresql://postgres.%s:%s@aws-1-ap-south-1.pooler.supabase.com:6543/postgres?sslmode=require&statement_cache_mode=describe",
			projectRef,
			encodedPassword,
		)

		pool, err = pgxpool.New(context.Background(), databaseURL)
		if err != nil {
			err = fmt.Errorf("failed to create connection pool: %v", err)
			return
		}

		// Test connection
		err = pool.Ping(context.Background())
		if err != nil {
			err = fmt.Errorf("failed to ping database: %v", err)
			return
		}

		// Store current version
		currentDBVersion = dbVersion
	})

	return err
}

// extractProjectRef extracts the project reference from Supabase URL
func extractProjectRef(supabaseURL string) string {
	// Remove https:// or http://
	supabaseURL = strings.TrimPrefix(supabaseURL, "https://")
	supabaseURL = strings.TrimPrefix(supabaseURL, "http://")

	// Extract project ref (before .supabase.co)
	parts := strings.Split(supabaseURL, ".supabase.co")
	if len(parts) > 0 && parts[0] != "" {
		return parts[0]
	}

	return ""
}

// GetPool returns the database connection pool
func GetPool() *pgxpool.Pool {
	return pool
}

// Close closes the database connection pool
func Close() {
	if pool != nil {
		pool.Close()
	}
}
