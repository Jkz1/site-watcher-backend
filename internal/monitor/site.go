package monitor

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"site-checker-backend/internal/models"
	"site-checker-backend/internal/repository"
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
