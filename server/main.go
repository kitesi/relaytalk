package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/kitesi/relaytalk/api"
	"github.com/kitesi/relaytalk/db"
)

func main() {
	godotenv.Load()

	dsn := os.Getenv("DATABASE_URL")
	port := os.Getenv("PORT")

	if dsn == "" {
		log.Fatal("DATABASE_URL environment variable is not set")
	}

	if port == "" {
		port = "8080"
	}

	pool, err := pgxpool.New(context.Background(), dsn)

	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}

	defer pool.Close()
	store := db.New(pool)

	r := chi.NewRouter()
	api.RegisterRoutes(store, r)

	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
