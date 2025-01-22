package main

import (
	"log"
	"time"

	"github.com/MdSadiqMd/RSS-Scraper/internal/database"
)

func startScrapping(
	db *database.Queries,
	concurrency int,
	timeBetweenRequest time.Duration,
) {
	log.Printf("Scrapping on %v goroutines every %v", concurrency, timeBetweenRequest)
	ticker:=time.NewTicker(timeBetweenRequest)
	for ; ; <-ticker.C {
		for i := 0; i < concurrency; i++ {
			go func() {
				db.Scrapping()
			}()
		}
	}
}
