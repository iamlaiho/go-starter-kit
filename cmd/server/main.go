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

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"

	"github.com/iamlaiho/go-starter-kit/internal/config"
	"github.com/iamlaiho/go-starter-kit/internal/handler"
	"github.com/iamlaiho/go-starter-kit/internal/middleware"
	"github.com/iamlaiho/go-starter-kit/internal/observability"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	logger := observability.NewLogger(cfg.AppEnv)
	slog.SetDefault(logger)

	r := chi.NewRouter()

	// Global middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger(logger))
	r.Use(middleware.Recoverer(logger))
	r.Use(middleware.SecureHeaders)
	r.Use(middleware.CORS(cfg.CORSAllowedOrigins))
	r.Use(middleware.Timeout(30 * time.Second))
	r.Use(middleware.BodyLimit(cfg.MaxRequestBodyBytes))
	r.Use(chimiddleware.StripSlashes)

	// Public routes
	r.Get("/health", handler.Health)

	// Protected routes (authenticated)
	r.Group(func(r chi.Router) {
		r.Use(middleware.Authenticate([]byte(cfg.JWTSecret)))
		// Add protected routes here
	})

	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in background
	go func() {
		logger.Info("server starting", "port", cfg.Port, "env", cfg.AppEnv)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("server error", "error", err)
			os.Exit(1)
		}
	}()

	// Graceful shutdown on SIGINT or SIGTERM
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("shutting down server")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("forced shutdown", "error", err)
		os.Exit(1)
	}

	logger.Info("server stopped")
}
