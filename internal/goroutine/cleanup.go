package goroutine

import (
	"log"
	"site-checker-backend/internal/repository"
	"time"
)

func StartCleanup(db *repository.SitesRepo) {
	ticker := time.NewTicker(1 * time.Hour)

	go func() {
		defer ticker.Stop()

		for range ticker.C {
			query := `DELETE FROM health_checks WHERE checked_at < NOW() - INTERVAL '30 days'`

			result, err := db.DB.Exec(query)
			if err != nil {
				log.Printf("Cleanup failed: %v", err)
				continue
			}

			rows, _ := result.RowsAffected()
			if rows > 0 {
				log.Printf("Cleanup successful: %d rows removed", rows)
			}
		}
	}()
}
