package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/kitesi/relaytalk/db"
	"github.com/kitesi/relaytalk/handlers"
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

	http.HandleFunc("/register", handlers.Register(store))
	http.HandleFunc("/login", handlers.Login(store))
	http.HandleFunc("/protected-ping", handlers.AuthMiddleware(store, handlers.ProtectedPing(store)))
	http.HandleFunc("/messages", handlers.AuthMiddleware(store, handlers.SendMessage(store)))

	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
