package handlers

import (
	"encoding/json"
	"net/http"
	"site-checker-backend/internal/models"
	"site-checker-backend/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	Repo *repository.UserRepo
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req models.AuthRequest
	json.NewDecoder(r.Body).Decode(&req)

	// Hash password before saving
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
	json.NewDecoder(r.Body).Decode(&req)

	user, err := h.Repo.GetByUsername(req.Username)
	if err != nil || bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)) != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"message": "Welcome " + user.Username})
}

func (h *UserHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	var req models.UpdatePasswordRequest
	json.NewDecoder(r.Body).Decode(&req)

	user, err := h.Repo.GetByUsername(req.Username)
	if err != nil || bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.OldPassword)) != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	newHash, _ := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	h.Repo.UpdatePassword(req.Username, string(newHash))
	w.Write([]byte("Password updated"))
}
