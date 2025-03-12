package main

import (
	"bytes"
	"context"
	"log"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/html"

	"github.com/divakaivan/rssagg/internal/database"
	"github.com/google/uuid"
)

func startScraping(
	db *database.Queries,
	concurrency int,
	timeBetweenRequest time.Duration,
) {
	log.Printf("[Info] Scraping on %v goroutings every %s duration", concurrency, timeBetweenRequest)
	ticker := time.NewTicker(timeBetweenRequest)
	// this runs the body immediately and THEN it wais
	// if we did range <- ticker.C, it would wait for the first tick
	for ; ; <-ticker.C {
		feeds, err := db.GetNextFeedsToFetch(
			context.Background(),
			int32(concurrency),
		)
		if err != nil {
			log.Printf("[Error] getting feeds to fetch: %v", err)
			continue
		}

		wg := &sync.WaitGroup{}
		for _, feed := range feeds {
			if feed.Url == "" {
				log.Printf("[Info] Feed %s has no URL", feed.Name)
				continue
			}
			wg.Add(1) // increments the counter by 1
			// spawn concurrency number of goroutines
			go scrapeFeed(db, wg, feed)
		}
		wg.Wait() // after all goroutines are done, we continue
	}
}

func scrapeFeed(db *database.Queries, wg *sync.WaitGroup, feed database.Feed) {
	defer wg.Done() // decrements the counter by 1
	_, err := db.MarkFeedAsFetched(context.Background(), feed.ID)
	if err != nil {
		log.Printf("[Error] marking feed as fetched: %v", err)
		return
	}
	rssFeed, err := urlToFeed(feed.Url)
	if err != nil {
		log.Printf("[Error] fetching feed: %v", err)
		return
	}
	for _, entry := range rssFeed.Entry {
		_, err := db.CreatePost(context.Background(),
			database.CreatePostParams{
				ID:          uuid.New(),
				CreatedAt:   time.Now().UTC(),
				UpdatedAt:   time.Now().UTC(),
				Title:       entry.Title,
				Summary:     entry.Summary,
				Content:     stripHTML(entry.Content),
				Url:         entry.Link.Href,
				PublishedAt: entry.Published,
				FeedID:      feed.ID,
			})
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
				continue
			}
			log.Printf("[Error] creating post: %v", err)
			continue
		}
	}
	log.Printf("[Info] Feed %s collected, %v posts found", feed.Name, len(rssFeed.Entry))
}

func stripHTML(htmlStr string) string {
	doc, err := html.Parse(strings.NewReader(htmlStr))
	if err != nil {
		return htmlStr
	}

	var buf bytes.Buffer
	var extractText func(*html.Node)
	extractText = func(n *html.Node) {
		if n.Type == html.TextNode {
			buf.WriteString(n.Data + " ")
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			extractText(c)
		}
	}
	extractText(doc)

	return strings.TrimSpace(buf.String())
}
