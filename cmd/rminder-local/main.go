package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"rminder/internal/pkg/config"
	"rminder/internal/pkg/logger"
	"rminder/internal/router"
)

func main() {
	log := slog.Default()
	log.Info("Starting Platform (local mode)")

	configDir, err := os.UserConfigDir()
	if err != nil {
		log.Error("Failed to get user config dir", "error", err)
		os.Exit(1)
	}
	dbDir := filepath.Join(configDir, "rminder")
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		log.Error("Failed to create rminder config dir", "error", err)
		os.Exit(1)
	}
	dbPath := filepath.Join(dbDir, "rminder.db")

	cfg := &config.Config{
		Server: config.ServerConfig{
			AuthPort:     4002,
			ReadTimeout:  "15s",
			WriteTimeout: "15s",
		},
		Logging: config.LoggingConfig{
			Level:  "info",
			Format: "text",
			Output: "stdout",
		},
	}

	appLogger, err := logger.New(logger.Config{
		Level:  cfg.Logging.Level,
		Format: cfg.Logging.Format,
		Output: cfg.Logging.Output,
	})
	if err != nil {
		log.Error("Failed to initialize logger", "error", err)
		os.Exit(1)
	}

	rtr := router.NewLocal(appLogger, cfg, dbPath)

	addr := fmt.Sprintf(":%d", cfg.Server.AuthPort)
	appLogger.Info("Rminder (local) starting", "addr", addr, "db", dbPath)

	readTimeout, _ := config.ParseDuration(cfg.Server.ReadTimeout)
	writeTimeout, _ := config.ParseDuration(cfg.Server.WriteTimeout)

	srv := &http.Server{
		Addr:         addr,
		Handler:      rtr,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
	}

	serverErrors := make(chan error, 1)
	go func() {
		serverErrors <- srv.ListenAndServe()
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		appLogger.Error("Server error", "error", err)
		os.Exit(1)
	case sig := <-shutdown:
		appLogger.Info("Shutdown signal received", "signal", sig)

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			appLogger.Error("Graceful shutdown failed", "error", err)
			if err := srv.Close(); err != nil {
				appLogger.Error("Server close failed", "error", err)
			}
		}

		appLogger.Info("Rminder (local) stopped")
		appLogger.Close()
	}
}
