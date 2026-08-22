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
	"foodapp/internal/sessions"
	"foodapp/internal/users"
)

func main() {
	// Context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Database
	cfg := config.LoadConfig()

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

	// Users
	userRepo := users.NewRepo(db)
	userService := users.NewService(userRepo, security.SecurityKeys{
		Encryption: cfg.EmailEncryptionKey,
		Hash:       cfg.EmailHashKey,
	})
	userHandler := users.NewHandler(userService)
	userModule := users.NewModule(userHandler)

	// Sessions
	sessionRepo := sessions.Newrepo(db)
	sessionService := sessions.NewService(sessionRepo, security.SecurityKeys{
		Encryption: cfg.EmailEncryptionKey,
		Hash:       cfg.EmailHashKey,
	}, userRepo)
	sessionHandler := sessions.NewHandler(sessionService)
	sessionModule := sessions.NewModule(sessionHandler)

	// Auth
	authService := auth.NewService(userRepo, sessionRepo)
	authMW := auth.NewMiddleware(authService)

	// Ingredients
	ingredientRepo := ingredients.Newrepo(db)
	ingredientService := ingredients.NewService(ingredientRepo)
	ingredientHandler := ingredients.NewHandler(ingredientService)
	ingredientModule := ingredients.NewModule(ingredientHandler)

	// Routes
	userModule.RegisterRoutes(mux, authMW.Authenticate, authMW.Admin)
	sessionModule.RegisterRoutes(mux)
	ingredientModule.RegisterRoutes(mux, authMW.Authenticate)

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
