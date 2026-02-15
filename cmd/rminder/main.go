package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"rminder/internal/login/authenticator"
	"rminder/internal/pkg/config"
	"rminder/internal/pkg/logger"
	"rminder/internal/router"
)

func main() {
	// Initialise Logger
	log := slog.Default()
	log.Info("Starting Platform")

	// Load configuration
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "configs/config.json"
	}

	cfg, err := config.Load(configPath)
	if err != nil {
		log.Error("Failed to load config:", "error", err)
	}

	// Initialise App Logger
	appLogger, err := logger.New(logger.Config{
		Level:  cfg.Logging.Level,
		Format: cfg.Logging.Format,
		Output: cfg.Logging.Output,
	})
	if err != nil {
		log.Error("Failed to initialize logger", "error", err)
		os.Exit(1)
	}

	auth := authenticator.New(cfg.Auth)

	// Initialize Routes
	rtr := router.New(auth, appLogger, cfg)

	// Starting server
	addr := fmt.Sprintf(":%d", cfg.Server.AuthPort)
	appLogger.Info("Auth service starting", "addr", addr)

	readTimeout, _ := config.ParseDuration(cfg.Server.ReadTimeout)
	writeTimeout, _ := config.ParseDuration(cfg.Server.WriteTimeout)

	srv := &http.Server{
		Addr:         addr,
		Handler:      rtr,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
	}

	// Start server in a goroutine
	serverErrors := make(chan error, 1)
	go func() {
		appLogger.Info("Rminder started", "port", addr)
		serverErrors <- srv.ListenAndServe()
	}()

	// Handle graceful shutdown
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

		appLogger.Info("Rminder stopped")
		appLogger.Close()
	}
}
