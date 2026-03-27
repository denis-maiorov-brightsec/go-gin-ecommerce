package testutil

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"testing"
	"time"

	"go-gin-ecommerce/internal/platform/config"
	platformdb "go-gin-ecommerce/internal/platform/db"

	"gorm.io/gorm"
)

const defaultDatabaseURL = "postgres://postgres:postgres@localhost:5432/go_gin_ecommerce?sslmode=disable"

func NewTestConfig(t *testing.T) config.Config {
	t.Helper()

	baseDatabaseURL := os.Getenv("DATABASE_URL")
	if baseDatabaseURL == "" {
		baseDatabaseURL = defaultDatabaseURL
	}

	schemaName := fmt.Sprintf("itest_%d", time.Now().UnixNano())
	adminConfig := config.Config{
		AppEnv:      "test",
		Port:        "0",
		LogLevel:    "error",
		DatabaseURL: baseDatabaseURL,
	}
	adminDB, err := platformdb.Open(adminConfig)
	if err != nil {
		t.Fatalf("failed to connect to postgres: %v", err)
	}

	if err := adminDB.Exec(fmt.Sprintf(`CREATE SCHEMA "%s"`, schemaName)).Error; err != nil {
		t.Fatalf("failed to create test schema: %v", err)
	}

	t.Cleanup(func() {
		if err := adminDB.Exec(fmt.Sprintf(`DROP SCHEMA IF EXISTS "%s" CASCADE`, schemaName)).Error; err != nil {
			t.Fatalf("failed to drop test schema: %v", err)
		}

		sqlDB, err := adminDB.DB()
		if err == nil {
			_ = sqlDB.Close()
		}
	})

	return config.Config{
		AppEnv:      "test",
		Port:        "0",
		LogLevel:    "error",
		DatabaseURL: withSearchPath(t, baseDatabaseURL, schemaName),
	}
}

func NewTestDatabase(t *testing.T, cfg config.Config) *gorm.DB {
	t.Helper()

	database, err := platformdb.Open(cfg)
	if err != nil {
		t.Fatalf("failed to open test database: %v", err)
	}
	t.Cleanup(func() {
		sqlDB, err := database.DB()
		if err == nil {
			_ = sqlDB.Close()
		}
	})

	applyMigrations(t, database)

	return database
}

func applyMigrations(t *testing.T, database *gorm.DB) {
	t.Helper()

	files, err := filepath.Glob(filepath.Join(repositoryRoot(t), "db", "migrations", "*.up.sql"))
	if err != nil {
		t.Fatalf("failed to discover migrations: %v", err)
	}
	sort.Strings(files)

	for _, path := range files {
		contents, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("failed to read migration %s: %v", path, err)
		}

		sql := strings.TrimSpace(string(contents))
		if sql == "" {
			continue
		}

		if err := database.Exec(sql).Error; err != nil {
			t.Fatalf("failed to apply migration %s: %v", path, err)
		}
	}
}

func withSearchPath(t *testing.T, rawURL string, schemaName string) string {
	t.Helper()

	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		t.Fatalf("failed to parse DATABASE_URL: %v", err)
	}

	query := parsedURL.Query()
	query.Set("search_path", schemaName)
	parsedURL.RawQuery = query.Encode()

	return parsedURL.String()
}

func repositoryRoot(t *testing.T) string {
	t.Helper()

	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("failed to resolve test helper path")
	}

	return filepath.Clean(filepath.Join(filepath.Dir(file), "..", "..", ".."))
}
