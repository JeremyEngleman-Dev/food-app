package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"Feastio/internal/platform/config"
	"Feastio/internal/platform/database"

	"Feastio/internal/ingredients"
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

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		greeting := "Welcome"
		json.NewEncoder(w).Encode(greeting)
	})

	usersModule := users.New(db, cfg)
	usersModule.RegisterRoutes(mux)

	ingredientsModule := ingredients.New(db)
	ingredientsModule.RegisterRoutes(mux)

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
