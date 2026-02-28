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

func (h *SiteHandler) UpdateActiveStatus(w http.ResponseWriter, r *http.Request) {
	var req models.UpdateActiveRequest
	json.NewDecoder(r.Body).Decode(&req)
	userID := r.Context().Value("user_id").(int)
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
	json.NewDecoder(r.Body).Decode(&req)
	userID := r.Context().Value("user_id").(int)
	sites, err := h.Repo.GetHistoryBySite(userID, req.SiteID)
	if err != nil {
		http.Error(w, "Failed to get history", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(sites)
}
