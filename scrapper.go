package main

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/MdSadiqMd/RSS-Scraper/internal/database"
)

func startScrapping(
	db *database.Queries,
	concurrency int,
	timeBetweenRequest time.Duration,
) {
	log.Printf("🪛 Scrapping on %v goroutines every %v", concurrency, timeBetweenRequest)
	ticker := time.NewTicker(timeBetweenRequest)
	for ; ; <-ticker.C {
		feeds, err := db.GetNextFeedsToFetch(
			context.Background(),
			int32(concurrency),
		)
		if err != nil {
			log.Println("⚠️ Unable to fetch feeds: ", err)
			continue
		}

		wg := sync.WaitGroup{}
		for _, feed := range feeds {
			wg.Add(1)
			go func(
				db *database.Queries,
				wg *sync.WaitGroup,
				feed database.Feed,
			) {
				defer wg.Done()
				_, err := db.MarkFeedAsFetched(context.Background(), feed.ID)
				if err != nil {
					log.Println("💀 Unable to mark feed as fetched: ", err)
					return
				}

				rssFeed, err := urlToFeedURL(feed.Url)
				if err != nil {
					log.Println("💀 Unable to fetch feed: ", err)
					return
				}

				for _, item := range rssFeed.Channel.Items {
					log.Println("📰 Adding item: ", item.Title)
				}
				log.Println("✅ Successfully scrapped feed of name: %v, which consists of %v posts: ", feed.Name, len(rssFeed.Channel.Items))
			}(db, &wg, feed)
		}
		wg.Wait()
	}
}
