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

// CreateSite godoc
// @Summary      Create a new site
// @Description  Creates a new site record. Requires a valid JWT in the Authorization header.
// @Tags         sites
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        site  body      models.CreateSiteRequest  true  "Site URL"
// @Success      201   {string}  string       "Created"
// @Failure      400   {string}  string       "Invalid request body"
// @Failure      401   {string}  string       "Unauthorized"
// @Failure      500   {string}  string       "Failed to create site"
// @Router       /sites [post]
func (h *SiteHandler) CreateSite(w http.ResponseWriter, r *http.Request) {
	var req models.CreateSiteRequest
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
// GetMySites godoc
// @Summary      Get all sites for current user
// @Description  Retrieves a list of all sites associated with the authenticated user's ID.
// @Tags         sites
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Success      200  {array}   models.Site  "Successfully retrieved list of sites"
// @Failure      401  {string}  string       "Unauthorized - Invalid or missing JWT"
// @Failure      500  {string}  string       "Internal Server Error"
// @Router       /sites [get]
func (h *SiteHandler) GetMySites(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	sites, err := h.Repo.GetUsersSite(userID)
    if err != nil {
        http.Error(w, "Failed to retrieve sites", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sites)
}
// UpdateActiveStatus godoc
// @Summary      Update site active status
// @Description  Toggles the 'is_active' state of a specific site owned by the authenticated user.
// @Tags         sites
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        status  body      models.UpdateActiveRequest  true  "Site ID and new Active status"
// @Success      200     {string}  string                      "OK"
// @Failure      400     {string}  string                      "Invalid request body"
// @Failure      401     {string}  string                      "Unauthorized - Missing or invalid JWT"
// @Failure      500     {string}  string                      "Internal Server Error - Database update failed"
// @Router       /sites/activated [put]
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
// GetHistory godoc
// @Summary      Get check history for a site
// @Description  Retrieves all historical status checks for a specific site. Requires a SiteID in the request body.
// @Tags         sites
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        request body      models.HistoryRequest  true  "History Filter (contains site_id)"
// @Success      200     {array}   models.HealthCheck     "List of historical check results"
// @Failure      400     {string}  string                 "Invalid request body"
// @Failure      401     {string}  string                 "Unauthorized"
// @Failure      500     {string}  string                 "Internal Server Error"
// @Router       /sites/history [put]
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
