package repository

import (
	"site-checker-backend/internal/models"

	"github.com/jmoiron/sqlx"
)

type SitesRepo struct {
	DB *sqlx.DB
}

func (r *SitesRepo) Create(userID int, url string) error {
	_, err := r.DB.Exec("INSERT INTO sites (user_id, url) VALUES ($1, $2)", userID, url)
	return err
}

func (r *SitesRepo) GetUsersSite(userID int) ([]models.Site, error) {
	var sites []models.Site
	err := r.DB.Select(&sites, "SELECT * FROM sites WHERE user_id=$1", userID)
	return sites, err
}
