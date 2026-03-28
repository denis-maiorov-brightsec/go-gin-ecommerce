package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-gin-ecommerce/internal/http/routes"
	"go-gin-ecommerce/internal/platform/config"
	platformdb "go-gin-ecommerce/internal/platform/db"
	"go-gin-ecommerce/internal/platform/httpserver"
	"go-gin-ecommerce/internal/platform/logger"
)

// @title Go Gin Ecommerce Backoffice API
// @version 1.0
// @description Versioned backoffice API for ecommerce management.
// @BasePath /v1
// @schemes http
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	cfg, err := config.Load()
	if err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	appLogger := logger.New(cfg.LogLevel)
	database, err := platformdb.Open(cfg)
	if err != nil {
		appLogger.Error("failed to open database", "error", err)
		os.Exit(1)
	}

	router := routes.NewWithDB(cfg, appLogger, database)
	server := httpserver.New(cfg.HTTPAddr(), router)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		appLogger.Info("starting HTTP server", "addr", cfg.HTTPAddr(), "env", cfg.AppEnv)

		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			appLogger.Error("server stopped unexpectedly", "error", err)
			os.Exit(1)
		}
	}()

	<-ctx.Done()
	appLogger.Info("shutdown signal received")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		appLogger.Error("graceful shutdown failed", "error", err)
		os.Exit(1)
	}
}
