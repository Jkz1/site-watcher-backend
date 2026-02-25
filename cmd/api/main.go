package main

import (
	"log"
	"net/http"
	"os"
	"site-checker-backend/internal/handlers"
	"site-checker-backend/internal/repository"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	godotenv.Load() // Loads .env into OS env

	connstr := os.Getenv("DATABASE_URL")
	db, err := sqlx.Connect("postgres", connstr)
	if err != nil {
		log.Fatal("DB Connection failed:", err)
	}

	// Auto-create table
	db.MustExec(`
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		username TEXT UNIQUE NOT NULL,
		password TEXT NOT NULL
	);`)
	db.MustExec(`
	CREATE TABLE sites (
		id SERIAL PRIMARY KEY,
		user_id INT REFERENCES users(id) ON DELETE CASCADE, 
		url TEXT NOT NULL,
		last_status INT DEFAULT 0,
		latency_ms INT DEFAULT 0,
		is_active BOOLEAN DEFAULT TRUE,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`)

	repo := &repository.UserRepo{DB: db}
	h := &handlers.UserHandler{Repo: repo}
	mux := http.NewServeMux()
	log.Println("Database schema initialized.")
	// Routes
	mux.HandleFunc("POST /register", h.Register)
	mux.HandleFunc("POST /login", h.Login)
	mux.HandleFunc("PUT /password", h.ChangePassword)
	log.Println("Server running on :8080")
	http.ListenAndServe(":8080", mux)
}
