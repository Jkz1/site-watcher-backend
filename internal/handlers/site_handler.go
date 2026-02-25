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
	// Implementation for creating a site
	var req models.Site
	json.NewDecoder(r.Body).Decode(&req)

}
