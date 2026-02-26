package models

import "time"

type HealthCheck struct {
	ID         int       `db:"id" json:"id"`
	SiteID     int       `db:"site_id" json:"site_id"`
	StatusCode int       `db:"status_code" json:"status_code"`
	LatencyMS  int       `db:"latency_ms" json:"latency_ms"`
	CheckedAt  time.Time `db:"checked_at" json:"checked_at"`
}

type HistoryRequest struct {
	SiteID int `json:"site_id"`
}
