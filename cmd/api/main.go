package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"foodapp/internal/platform/config"
	"foodapp/internal/platform/database"
	"foodapp/internal/platform/security"

	"foodapp/internal/auth"
	"foodapp/internal/ingredients"
	"foodapp/internal/mappings"
	"foodapp/internal/users"
)

func main() {
	// Context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Configuration
	cfg := config.LoadConfig()

	// Tools
	encryption := security.NewEncryption(cfg.Secrets.EmailEncryptionKey)
	hash := security.NewHash(cfg.Secrets.EmailHashKey, cfg.Secrets.TokenHashKey)
	token := security.NewToken(cfg.Secrets.TokenHashKey)
	userMapping := mappings.NewUserMap(encryption)

	// Database
	db, err := database.Open(cfg.DBString)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal("Database not reachable:", err)
	}

	if err := database.RunMigrations(ctx, db, "internal/platform/database/migrations"); err != nil {
		log.Fatal(err)
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
	srv := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	log.Printf("Server starting on %s ...", srv.Addr)
	log.Fatal(srv.ListenAndServe())
}
