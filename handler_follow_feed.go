package main

import (
	"context"
	"fmt"
	"time"

	"github.com/PlinyTheYounger0/aggregatr/internal/database"
	"github.com/google/uuid"
)

func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("Usage: follow <feed-url>")
	}

	feed, err := s.db.GetFeedByURL(context.Background(), cmd.Args[0])
	if err != nil {
		return fmt.Errorf("Error Fetching Feed to Follow: %w", err)
	}

	feedFollow, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return fmt.Errorf("Error Creating Feed Follow: %w", err)
	}

	fmt.Printf("%s Feed Followed Sucessfully by %s", feedFollow.FeedName, user.Name)

	return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 0 {
		return fmt.Errorf("Usage: following")
	}

	followFeeds, err := s.db.GetFeedFollowsByUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("Error Fetching User Follow Feeds: %w", err)
	}

	for i, feed := range followFeeds {
		feedName, err := s.db.GetFeedNameByID(context.Background(), feed.FeedID)
		if err != nil {
			return fmt.Errorf("Error Fetching Feed Name by ID: %w", err)
		}

		fmt.Printf("%d. %s\n", i, feedName)
	}

	return nil
}

func handlerUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("Usage: unfollow <feed-url>")
	}

	feed, err := s.db.GetFeedByURL(context.Background(), cmd.Args[0])
	if err != nil {
		return fmt.Errorf("Error Fetching Feed to Unfollow: %w", err)
	}

	err = s.db.DeleteFeedFollow(context.Background(), database.DeleteFeedFollowParams{
		UserID: user.ID,
		FeedID: feed.ID,
	})
	if err != nil {
		return fmt.Errorf("Error Deleting Follow Feed: %w", err)
	}

	fmt.Printf("%s Feed Successfully Unfollowed by %s\n", feed.Name, user.Name)

	return nil
}
