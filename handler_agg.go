package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/PlinyTheYounger0/aggregatr/internal/database"
)

func handlerAgg(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("Usage: agg <time-between-reqs>")
	}

	timeBtwReqs, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		return fmt.Errorf("Error Parsing Time Between Reqs: %w", err)
	}

	fmt.Printf("Collecting Feeds Every %v", timeBtwReqs)

	ticker := time.NewTicker(timeBtwReqs)
	for ; ; <-ticker.C {
		scarpeFeeds(s)
	}

	return nil
}

func scarpeFeeds(s *state) error {
	nextFeedToFetch, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return fmt.Errorf("Error Fetching Next Feed: %w", err)
	}

	feed, err := fecthFeed(context.Background(), nextFeedToFetch.Url)
	if err != nil {
		return fmt.Errorf("Error Fetching Feed During Scrape: %w", err)
	}

	for _, item := range feed.Channel.Item {
		fmt.Printf("%s: %s\n", feed.Channel.Title, item.Title)
	}

	err = s.db.MarkFeedFetched(context.Background(), database.MarkFeedFetchedParams{
		LastFetchedAt: sql.NullTime{
			Time:  time.Now().UTC(),
			Valid: true,
		},
		UpdatedAt: time.Now().UTC(),
		ID:        nextFeedToFetch.ID,
	})
	if err != nil {
		return fmt.Errorf("Error Marking Feed As Fetched: %w", err)
	}

	return nil
}
