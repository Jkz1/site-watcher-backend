package handlers

import (
	"encoding/json"
	"net/http"
	"site-checker-backend/internal/models"
	"site-checker-backend/internal/repository"
)

type SiteHandler struct {
	Repo *repository.SitesRepo
}

func (h *SiteHandler) CreateSite(w http.ResponseWriter, r *http.Request) {
	var req models.Site
	json.NewDecoder(r.Body).Decode(&req)
	userID := r.Context().Value("user_id").(int)
	err := h.Repo.Create(userID, req.URL)
	if err != nil {
		http.Error(w, "Failed to create site", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *SiteHandler) GetMySites(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int)
	sites, _ := h.Repo.GetUsersSite(userID)
	json.NewEncoder(w).Encode(sites)
}

func (h *SiteHandler) ToggleMonitoring(w http.ResponseWriter, r *http.Request) {
	// 1. Get UserID from Middleware context
	userID := r.Context().Value("user_id").(int)

	// 2. Get SiteID (Assuming it's a URL param or JSON body)
	var req struct {
		SiteID int  `json:"site_id"`
		Active bool `json:"active"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// 3. Update the DB
	err := h.Repo.UpdateActiveStatus(req.SiteID, userID, req.Active)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
