package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"time"

	"foodapp/internal/logging"
	"foodapp/internal/platform/config"
	"foodapp/internal/platform/database"
	"foodapp/internal/platform/security"

	"foodapp/internal/auth"
	"foodapp/internal/ingredients"
	"foodapp/internal/mappings"
	"foodapp/internal/users"
)

func main() {
	// Logging
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	slog.SetDefault(logger)

	// Context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Configuration
	cfg := config.LoadConfig()

	// Utilities
	encryption := security.NewEncryption(cfg.Secrets.EmailEncryptionKey)
	hash := security.NewHash(cfg.Secrets.EmailHashKey, cfg.Secrets.TokenHashKey)
	token := security.NewToken(cfg.Secrets.TokenHashKey)
	userMapping := mappings.NewUserMap(encryption)

	// Database
	db, err := database.Open(cfg.DBString)
	if err != nil {
		logger.Error("Failed to load database", "error", err.Error())
		os.Exit(1)
	}
	defer db.Close()

	if err := database.RunMigrations(ctx, db, "internal/platform/database/migrations"); err != nil {
		logger.Error("Database migration error", "error", err.Error())
		os.Exit(1)
	}

	// HTTP
	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	// HTTP - Users
	userRepo := users.NewRepo(db)
	userService := users.NewService(userRepo, encryption, hash)
	userHandler := users.NewHandler(userService, userMapping)

	// HTTP - Auth
	authentication := auth.NewMiddleware(token)
	authRepo := auth.NewRepo(db)
	authService := auth.NewService(authRepo, hash, userRepo, token)
	authHandler := auth.NewHandler(authService)

	// HTTP - Ingredients
	ingredientRepo := ingredients.NewRepo(db)
	ingredientService := ingredients.NewService(ingredientRepo)
	ingredientHandler := ingredients.NewHandler(ingredientService)

	// HTTP - Routes
	userHandler.RegisterRoutes(mux, authentication.Authenticate, authentication.Admin)
	authHandler.RegisterRoutes(mux, authentication.Authenticate)
	ingredientHandler.RegisterRoutes(mux, authentication.Authenticate)

	// Server
	serverHandler := logging.Log(logger)(mux)

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      serverHandler,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	logger.Info("Server starting", "address", srv.Addr)
	if err := srv.ListenAndServe(); err != nil {
		logger.Error("server stopped", "err", err)
		os.Exit(1)
	}
}
