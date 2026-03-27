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
	"go-gin-ecommerce/internal/platform/httpserver"
	"go-gin-ecommerce/internal/platform/logger"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	appLogger := logger.New(cfg.LogLevel)
	router := routes.New(cfg, appLogger)
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
