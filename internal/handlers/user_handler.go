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
// Register godoc
// @Summary      Register a new user
// @Description  Creates a new user account with a username and password. The password will be hashed using bcrypt.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        user  body      models.AuthRequest  true  "User Registration Details"
// @Success      201   {string}  string              "Created"
// @Failure      400   {string}  string              "Invalid request body"
// @Failure      409   {string}  string              "Username already taken"
// @Failure      500   {string}  string              "Internal server error"
// @Router       /register [post]
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
// Login godoc
// @Summary      User Login
// @Description  Authenticates a user and returns a JWT token valid for 24 hours.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        credentials  body      models.AuthRequest  true  "Login Credentials"
// @Success      200          {object}  map[string]string   "returns {'token': 'jwt_token_here'}"
// @Failure      400          {string}  string              "Invalid request body"
// @Failure      401          {string}  string              "Invalid credentials"
// @Router       /login [post]
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

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		http.Error(w, "Could not generate token", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
}
// ChangePassword godoc
// @Summary      Change user password
// @Description  Updates the authenticated user's password. Requires the old password for verification and a valid JWT.
// @Tags         auth
// @Accept       json
// @Produce      plain
// @Security     ApiKeyAuth
// @Param        request body      models.UpdatePasswordRequest  true  "Old and New Password"
// @Success      200     {string}  string                        "Password updated"
// @Failure      400     {string}  string                        "Invalid request body"
// @Failure      401     {string}  string                        "Unauthorized - Missing token or incorrect old password"
// @Failure      500     {string}  string                        "Internal server error"
// @Router       /password [put]
func (h *UserHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	var req models.UpdatePasswordRequest
	userID, ok := r.Context().Value("user_id").(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user, err := h.Repo.GetByID(userID)
	if err != nil || bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.OldPassword)) != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	newHash, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}
	h.Repo.UpdatePassword(userID, string(newHash))
	w.Write([]byte("Password updated"))
}
