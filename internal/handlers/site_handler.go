package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"site-checker-backend/internal/models"
	"site-checker-backend/internal/repository"
)

type SiteHandler struct {
	Repo *repository.SitesRepo
}

// Create Site godoc
// @Summary      Create a new site by user
// @Description  Creates a new site for the authenticated user
// @Tags         sites
// @Accept       json
// @Produce      json
// @Success      201  {object}  models.Site
// @Failure      404  {object}  string
// @Router       /sites [post]
func (h *SiteHandler) CreateSite(w http.ResponseWriter, r *http.Request) {
	var req models.Site
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	userID, ok := r.Context().Value("user_id").(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	err := h.Repo.Create(userID, req.URL)
	if err != nil {
		http.Error(w, "Failed to create site", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *SiteHandler) GetMySites(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	sites, _ := h.Repo.GetUsersSite(userID)
	json.NewEncoder(w).Encode(sites)
}

func (h *SiteHandler) UpdateActiveStatus(w http.ResponseWriter, r *http.Request) {
	var req models.UpdateActiveRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	userID, ok := r.Context().Value("user_id").(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	err := h.Repo.UpdateActiveStatus(userID, req.SiteID, req.IsActive)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Failed to update site status", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *SiteHandler) GetHistory(w http.ResponseWriter, r *http.Request) {
	var req models.HistoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	userID, ok := r.Context().Value("user_id").(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	sites, err := h.Repo.GetHistoryBySite(userID, req.SiteID)
	if err != nil {
		http.Error(w, "Failed to get history", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(sites)
}
