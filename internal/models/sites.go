package models

import (
	"time"
)

type Site struct {
	ID         int        `db:"id" json:"id"`
	UserID     int        `db:"user_id" json:"user_id"`
	URL        string     `db:"url" json:"url"`
	LastStatus *int       `db:"last_status" json:"last_status"`
	LatencyMs  *int       `db:"latency_ms" json:"latency_ms"`
	IsActive   *bool      `db:"is_active" json:"is_active"`
	CreatedAt  *time.Time `db:"created_at" json:"created_at"`
}

type UpdateActiveRequest struct {
	SiteID   int  `json:"site_id"`
	IsActive bool `json:"is_active"`
}
