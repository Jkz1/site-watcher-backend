package main

import (
	"log"
	"net/http"
	"os"
	"site-checker-backend/internal/handlers"
	auth_middleware "site-checker-backend/internal/middleware"
	"site-checker-backend/internal/monitor"
	"site-checker-backend/internal/repository"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {

	godotenv.Load()

	connstr := os.Getenv("DATABASE_URL")
	db, err := sqlx.Connect("postgres", connstr)
	if err != nil {
		log.Fatal("DB Connection failed:", err)
	}

	db.MustExec(`
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		username TEXT UNIQUE NOT NULL,
		password TEXT NOT NULL
	);`)

	db.MustExec(`
	CREATE TABLE IF NOT EXISTS sites (
		id SERIAL PRIMARY KEY,
		user_id INT REFERENCES users(id) ON DELETE CASCADE, 
		url TEXT NOT NULL,
		last_status INT DEFAULT 0,
		latency_ms INT DEFAULT 0,
		is_active BOOLEAN DEFAULT TRUE,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`)
	db.MustExec(`
	CREATE TABLE health_checks (
		id SERIAL PRIMARY KEY,
		site_id INT REFERENCES sites(id) ON DELETE CASCADE,
		status_code INT,
		latency_ms INT,
		checked_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`)

	repo := &repository.UserRepo{DB: db}
	h := &handlers.UserHandler{Repo: repo}
	sitesRepo := &repository.SitesRepo{DB: db}
	siteHandler := &handlers.SiteHandler{Repo: sitesRepo}

	monitor.StartWorker(sitesRepo)
	monitor.StartJanitor(sitesRepo)
	mux := http.NewServeMux()
	log.Println("Database schema initialized.")
	// Routes
	mux.HandleFunc("POST /register", h.Register)
	mux.HandleFunc("POST /login", h.Login)
	mux.HandleFunc("PUT /password", h.ChangePassword)
	mux.HandleFunc("GET /sites", auth_middleware.AuthMiddleware(siteHandler.GetMySites))
	mux.HandleFunc("POST /sites", auth_middleware.AuthMiddleware(siteHandler.CreateSite))
	mux.HandleFunc("PUT /sites/history", auth_middleware.AuthMiddleware(siteHandler.GetHistory))
	mux.HandleFunc("PUT /sites/activated", auth_middleware.AuthMiddleware(sitehandler))
	log.Println("Server running on :8080")
	http.ListenAndServe(":8080", mux)
}
