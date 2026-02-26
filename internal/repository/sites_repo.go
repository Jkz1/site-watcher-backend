package repository

import (
	"site-checker-backend/internal/models"
	"time"

	"github.com/jmoiron/sqlx"
)

type SitesRepo struct {
	DB *sqlx.DB
}

func (r *SitesRepo) Create(userID int, url string) error {
	_, err := r.DB.Exec("INSERT INTO sites (user_id, url, is_active) VALUES ($1, $2, false)", userID, url)
	return err
}

func (r *SitesRepo) GetUsersSite(userID int) ([]models.Site, error) {
	var sites []models.Site
	err := r.DB.Select(&sites, "SELECT * FROM sites WHERE user_id=$1", userID)
	return sites, err
}

func (r *SitesRepo) GetAllActive() ([]models.Site, error) {
	var sites []models.Site
	err := r.DB.Select(&sites, "SELECT * FROM sites WHERE is_active=true")
	return sites, err
}

func (r *SitesRepo) UpdateSiteStatus(id int, status int, latency int) error {
	_, err := r.DB.Exec(
		"UPDATE sites SET last_status=$1, latency_ms=$2, last_checked=$3 WHERE id=$4",
		status, latency, time.Now(), id,
	)
	return err
}

func (r *SitesRepo) UpdateActiveStatus(siteID int, userID int, status bool) error {
	// We include userID to ensure User A can't start/stop User B's sites!
	_, err := r.DB.Exec(
		"UPDATE sites SET is_active=$1 WHERE id=$2 AND user_id=$3",
		status, siteID, userID,
	)
	return err
}
