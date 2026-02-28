package goroutine

import (
	"fmt"
	"log"
	"net/http"
	"site-checker-backend/internal/models"
	"site-checker-backend/internal/repository"
	"time"
)

func StartWorker(repo *repository.SitesRepo) {
	ticker := time.NewTicker(60 * time.Second)
	go func() {
		for range ticker.C {
			sites, err := repo.GetAllActive()
			if err != nil {
				log.Printf("Worker Error: %v", err)
				continue
			}

			if len(sites) == 0 {
				log.Println("No active sites to monitor. Skipping cycle.")
				continue
			}

			for _, site := range sites {
				go func(s models.Site) {
					status, latency := pingSite(s.URL)
					repo.UpdateSiteStatus(s.ID, status, latency)
					fmt.Println("pinged", s.URL, "status:", status, "latency:", latency, "ms")
				}(site)
			}
		}
	}()
}

func pingSite(url string) (int, int) {
	client := http.Client{
		Timeout: 10 * time.Second, // Don't wait forever
	}

	start := time.Now()
	resp, err := client.Get(url)
	if err != nil {
		return 0, 0 // 0 means site is unreachable
	}
	defer resp.Body.Close()

	latency := int(time.Since(start).Milliseconds())
	return resp.StatusCode, latency
}

func StartJanitor(repo *repository.SitesRepo) {
	// Check once every 24 hours
	ticker := time.NewTicker(24 * time.Hour)

	go func() {
		for range ticker.C {
			log.Println("[Janitor] Starting daily database cleanup...")

			rows, err := repo.CleanOldLogs()
			if err != nil {
				log.Printf("[Janitor] Error cleaning logs: %v", err)
				continue
			}

			log.Printf("[Janitor] Successfully deleted %d old health check records.", rows)
		}
	}()
}
