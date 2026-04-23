package main

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/PlinyTheYounger0/aggregatr/internal/database"
	"github.com/google/uuid"
)

func handlerAgg(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("Usage: agg <time-between-reqs>")
	}

	timeBtwReqs, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		return fmt.Errorf("Error Parsing Time Between Reqs: %w", err)
	}

	fmt.Printf("Collecting Feeds Every %v\n", timeBtwReqs)

	ticker := time.NewTicker(timeBtwReqs)
	for ; ; <-ticker.C {
		err := scarpeFeeds(s)
		if err != nil {
			fmt.Printf("Scraping Error: %v", err)
		}
	}
}

func scarpeFeeds(s *state) error {
	nextFeedToFetch, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return fmt.Errorf("Error Fetching Next Feed: %w", err)
	}

	rssFeed, err := fecthFeed(context.Background(), nextFeedToFetch.Url)
	if err != nil {
		return fmt.Errorf("Error Fetching Feed During Scrape: %w", err)
	}
	fmt.Printf("Fetched: %s\n", rssFeed.Channel.Title)
	fmt.Printf("# of Items Found: %d\n", len(rssFeed.Channel.Item))

	for _, item := range rssFeed.Channel.Item {

		pubDate, err := parseTime(item.PubDate)
		if err != nil {
			fmt.Printf("Error Parsing pubDate: %v", err)
		}

		_, err = s.db.CreatePost(context.Background(), database.CreatePostParams{
			ID:        uuid.New(),
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
			Title: sql.NullString{
				String: item.Title,
				Valid:  true,
			},
			Url: item.Link,
			Description: sql.NullString{
				String: item.Description,
				Valid:  true,
			},
			PublishedAt: sql.NullTime{
				Time:  pubDate,
				Valid: true,
			},
			FeedID: nextFeedToFetch.ID,
		})
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
				continue
			}
			fmt.Printf("Error Creating Post: %v", err)
			continue
		}
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

func parseTime(timeToParse string) (time.Time, error) {
	formats := []string{
		time.Layout,
		time.ANSIC,
		time.UnixDate,
		time.RubyDate,
		time.RFC822,
		time.RFC822Z,
		time.RFC850,
		time.RFC1123,
		time.RFC1123Z,
		time.RFC3339,
		time.RFC3339Nano,
		time.Kitchen,
	}

	for _, format := range formats {
		t, err := time.Parse(format, timeToParse)
		if err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("Error Parsing Time: %s", timeToParse)
}

func handlerBrowse(s *state, cmd command, user database.User) error {
	postLimit := int32(2)
	if len(cmd.Args) > 0 {
		n, err := strconv.Atoi(cmd.Args[0])
		if err != nil {
			return fmt.Errorf("Failed to Convert Arg to Int: %w", err)
		}
		postLimit = int32(n)
	}

	posts, err := s.db.GetPostsForUser(context.Background(), database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  postLimit,
	})
	if err != nil {
		return fmt.Errorf("Error Fetching Posts for User: %w", err)
	}

	if len(posts) == 0 {
		fmt.Println("No posts available.")
		return nil
	}
	for i, post := range posts {
		fmt.Printf("=============Post: %d=================\n", i+1)
		fmt.Printf("Title: %v\n", post.Title.String)
		fmt.Printf("URL: %v\n", post.Url)
		fmt.Printf("Description: %v\n", post.Description.String)
		fmt.Println()
	}

	return nil
}
