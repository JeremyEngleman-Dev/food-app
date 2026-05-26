package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"Feastio/internal/platform/config"
	"Feastio/internal/platform/database"

	"Feastio/internal/ingredients"
	"Feastio/internal/middleware"
	"Feastio/internal/sessions"
	"Feastio/internal/users"
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

	if err := database.RunMigrations(ctx, db, "../../internal/platform/database/migrations"); err != nil {
		log.Fatal(err)
	}

	// HTTP
	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	mw := middleware.New(db)
	usersModule := users.New(db, cfg)
	ingredientsModule := ingredients.New(db)
	sessionsModule := sessions.New(db, cfg, usersModule.Service())

	usersModule.RegisterRoutes(mux, mw)
	ingredientsModule.RegisterRoutes(mux)
	sessionsModule.RegisterRoutes(mux)

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
