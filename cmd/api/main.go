package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"Feastio/internal/config"
	httpHandler "Feastio/internal/http"
	"Feastio/internal/repository"
	"Feastio/internal/service"
)

func main() {
	// Database
	const dbString = `
		host=localhost
		port=5433
		user=postgres
		password=postgres
		dbname=feastio
		sslmode=disable
	`

	cfg := config.LoadConfig()

	db, err := repository.Open(dbString)
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal("Database not reachable:", err)
	}

	// Context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Repo | Service | Handler
	repo := repository.NewRepository(db)
	service := service.NewService(repo, cfg.EmailKey)
	handler := httpHandler.NewHandler(service)

	if err := repo.RunMigrations(ctx, db, "../../internal/repository/migrations"); err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		path := "Welcome"
		json.NewEncoder(w).Encode(path)
	})

	http.HandleFunc("POST /ingredients", handler.CreateIngredient)
	http.HandleFunc("GET /ingredients", handler.ListIngredients)
	http.HandleFunc("GET /ingredients/", handler.GetIngredient)
	http.HandleFunc("DELETE /ingredients/", handler.DeleteIngredient)

	http.HandleFunc("POST /users", handler.CreateUser)
	http.HandleFunc("GET /users/", handler.GetUser)
	http.HandleFunc("GET /users", handler.ListUsers)
	http.HandleFunc("DELETE /users/", handler.DeleteUser)

	// Server
	srv := &http.Server{
		Addr:         ":8080",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	log.Printf("Server starting on %s ...", srv.Addr)
	log.Fatal(srv.ListenAndServe())
}
