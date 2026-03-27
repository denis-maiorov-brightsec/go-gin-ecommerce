package config_test

import (
	"testing"
	"time"

	"go-gin-ecommerce/internal/platform/config"
)

func TestEffectiveWriteRateLimitDefaultsForTestEnv(t *testing.T) {
	t.Parallel()

	cfg := config.Config{AppEnv: "test"}

	if cfg.EffectiveWriteRateLimitRequests() != 1000 {
		t.Fatalf("expected test write-rate-limit default to be 1000, got %d", cfg.EffectiveWriteRateLimitRequests())
	}
	if cfg.EffectiveWriteRateLimitWindow() != time.Minute {
		t.Fatalf("expected default write-rate-limit window to be 1m, got %s", cfg.EffectiveWriteRateLimitWindow())
	}
}

func TestEffectiveWriteRateLimitOverrides(t *testing.T) {
	t.Parallel()

	cfg := config.Config{
		AppEnv:                 "production",
		WriteRateLimitRequests: 7,
		WriteRateLimitWindow:   5 * time.Second,
	}

	if cfg.EffectiveWriteRateLimitRequests() != 7 {
		t.Fatalf("expected configured write-rate-limit requests to be used, got %d", cfg.EffectiveWriteRateLimitRequests())
	}
	if cfg.EffectiveWriteRateLimitWindow() != 5*time.Second {
		t.Fatalf("expected configured write-rate-limit window to be used, got %s", cfg.EffectiveWriteRateLimitWindow())
	}
}
