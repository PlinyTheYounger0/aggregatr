package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
	"time"

	"github.com/PlinyTheYounger0/aggregatr/internal/database"
	"github.com/google/uuid"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func fecthFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return nil, fmt.Errorf("Feed request generation error. %w", err)
	}
	req.Header.Set("User-Agent", "gatr")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Feed request fetching error. %w", err)
	}
	defer res.Body.Close()

	dat, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("Error reading feed. %w", err)
	}

	rssFeed := RSSFeed{}
	if err := xml.Unmarshal(dat, &rssFeed); err != nil {
		return nil, fmt.Errorf("Error unmarshaling feed. %w", err)
	}

	rssFeed.Channel.Title = html.UnescapeString(rssFeed.Channel.Title)
	rssFeed.Channel.Description = html.UnescapeString(rssFeed.Channel.Description)

	for i, item := range rssFeed.Channel.Item {
		item.Title = html.UnescapeString(item.Title)
		item.Description = html.UnescapeString(item.Description)
		rssFeed.Channel.Item[i] = item
	}

	return &rssFeed, nil
}

func handlerAddFeed(s *state, cmd command) error {
	if len(cmd.Args) != 2 {
		return fmt.Errorf("Usage: addfeed <feed-name> <feed-url>")
	}

	currentUser, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
	if err != nil {
		return fmt.Errorf("Error Retrieving User For addFeed: %w", err)
	}

	feed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      cmd.Args[0],
		Url:       cmd.Args[1],
		UserID:    currentUser.ID,
	})
	if err != nil {
		return fmt.Errorf("Error Creating Feed in DB: %w", err)
	}

	fmt.Printf("%s Feed Created\n", feed.Name)
	fmt.Printf("Feed URL: %s\n", feed.Url)
	fmt.Printf("Feed ID: %v\n", feed.ID)
	fmt.Printf("Feed Added By: %v\n", feed.UserID)

	return nil
}

func handlerListFeeds(s *state, cmd command) error {
	feeds, err := s.db.ListFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("Error Fetching Feeds: %w", err)
	}

	for _, feed := range feeds {
		userName, err := s.db.GetUserNameFromID(context.Background(), feed.UserID)
		if err != nil {
			return fmt.Errorf("Error Finding Adding User: %w")
		}

		fmt.Printf("Feed Name: %s\n", feed.Name)
		fmt.Printf("Feed URL: %s\n", feed.Url)
		fmt.Printf("Adding User: %s\n", userName)
	}

	return nil
}
