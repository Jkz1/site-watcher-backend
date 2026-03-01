package handlers

import (
	"encoding/json"
	"net/http"
	"os"
	"site-checker-backend/internal/models"
	"site-checker-backend/internal/repository"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	Repo *repository.UserRepo
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req models.AuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	hashed, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)

	user := &models.User{Username: req.Username, Password: string(hashed)}
	if err := h.Repo.Create(user); err != nil {
		http.Error(w, "Username taken", http.StatusConflict)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req models.AuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user, err := h.Repo.GetByUsername(req.Username)
	if err != nil || bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)) != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // Expires in 24h
	})

	tokenString, _ := token.SignedString([]byte(os.Getenv("JWT_SECRET")))

	json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
}

func (h *UserHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	var req models.UpdatePasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user, err := h.Repo.GetByUsername(req.Username)
	if err != nil || bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.OldPassword)) != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	newHash, _ := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	h.Repo.UpdatePassword(req.Username, string(newHash))
	w.Write([]byte("Password updated"))
}
