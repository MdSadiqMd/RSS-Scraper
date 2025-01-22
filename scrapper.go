package main

import (
	"context"
	"database/sql"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/MdSadiqMd/RSS-Scraper/internal/database"
	"github.com/google/uuid"
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
					description := sql.NullString{}
					if item.Description != "" {
						description = sql.NullString{
							String: item.Description,
							Valid:  true,
						}
					}

					pubAt, err := time.Parse(time.RFC1123, item.PubDate)
					if err != nil {
						log.Println("🗓️ Unable to parse Date: ", err)
					}

					_, err = db.CreatePost(context.Background(), database.CreatePostParams{
						ID:          uuid.New(),
						CreatedAt:   time.Now(),
						UpdatedAt:   time.Now(),
						Title:       item.Title,
						Description: description,
						PublishedAt: pubAt,
						Url:         item.Link,
						FeedID:      feed.ID,
					})
					if err != nil {
						if strings.Contains(err.Error(), "duplicate key") {
							continue
						}
						log.Println("🌋 Failed to Create Post: ", err)
					}
					log.Println("📰 Adding item: ", item.Title)
				}
				log.Printf("✅ Successfully scrapped %v posts", len(rssFeed.Channel.Items))
			}(db, &wg, feed)
		}
		wg.Wait()
	}
}
