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
	_, err = r.DB.Exec(
		"INSERT INTO health_checks (site_id, status_code, latency_ms) VALUES ($1, $2, $3)",
		id, status, latency,
	)
	return err
}

func (r *SitesRepo) UpdateActiveStatus(userID int, siteID int, status bool) error {
	_, err := r.DB.Exec(
		"UPDATE sites SET is_active=$1, last_checked=NOW() WHERE id=$2 AND user_id=$3",
		status, siteID, userID,
	)
	if err != nil {
		return err
	}

	return err
}
func (r *SitesRepo) CleanOldLogs() (int64, error) {
	// Delete logs where checked_at is older than 30 days
	result, err := r.DB.Exec("DELETE FROM health_checks WHERE checked_at < NOW() - INTERVAL '30 days'")
	if err != nil {
		return 0, err
	}

	// Return how many rows were deleted for logging purposes
	return result.RowsAffected()
}
func (r *SitesRepo) GetHistoryBySite(userID int, siteID int) ([]models.HealthCheck, error) {
	var history []models.HealthCheck

	// Select automatically maps the columns to the struct tags
	query := `SELECT * FROM health_checks WHERE site_id = $1 AND site_id IN (SELECT id FROM sites WHERE user_id = $2) ORDER BY checked_at DESC LIMIT 100`
	err := r.DB.Select(&history, query, siteID, userID)

	return history, err
}
