package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/ledsouza/aggregator/internal/database"
)

func handlerAgg(s *state, cmd command) error {
	url := "https://www.wagslane.dev/index.xml"
	feed, err := fetchFeed(context.Background(), url)
	if err != nil {
		return err
	}

	fmt.Printf("%+v\n", feed)
	return nil
}

func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.Args) < 2 {
		return fmt.Errorf("name and url are required")
	}

	name := cmd.Args[0]
	url := cmd.Args[1]

	feed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      name,
		Url:       url,
		UserID:    user.ID,
	})
	if err != nil {
		return fmt.Errorf("couldn't create the feed: %w", err)
	}

	_, err = s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return fmt.Errorf("couldn't create the feed follow: %w", err)
	}

	fmt.Printf("%+v\n", feed)
	return nil
}

func handlerFeeds(s *state, cmd command) error {
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("couldn't get feeds: %w", err)
	}

	fmt.Println("Feeds details:")
	for _, feed := range feeds {
		user, err := s.db.GetUserByID(context.Background(), feed.UserID)
		if err != nil {
			return fmt.Errorf("couldn't get user: %w", err)
		}

		fmt.Println("---------------------------")
		fmt.Printf("Name: %s\n", feed.Name)
		fmt.Printf("URL: %s\n", feed.Url)
		fmt.Printf("Created By: %s\n", user.Name)
		fmt.Println("---------------------------")
	}

	return nil
}

func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.Args) < 1 {
		return fmt.Errorf("url is required")
	}

	url := cmd.Args[0]

	feed, err := s.db.GetFeedByUrl(context.Background(), url)
	if err != nil {
		return fmt.Errorf("couldn't get feed by url: %w", err)
	}

	followFeed, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return fmt.Errorf("couldn't create follow feed: %w", err)
	}

	fmt.Println("Feed details:")
	fmt.Printf("Name: %s\n", followFeed.FeedName)
	fmt.Printf("User: %s\n", followFeed.UserName)

	return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error {
	feeds, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("couldn't get feeds: %w", err)
	}

	fmt.Println("The feeds that you are following are:")
	for _, feed := range feeds {
		fmt.Printf("%s\n", feed.FeedName)
	}

	return nil
}

func handlerUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.Args) < 1 {
		return fmt.Errorf("url is required")
	}
	url := cmd.Args[0]

	err := s.db.DeleteFeedFollowByUrl(context.Background(), database.DeleteFeedFollowByUrlParams{
		Url:    url,
		UserID: user.ID,
	})
	if err != nil {
		return fmt.Errorf("couldn't delete feed follow: %w", err)
	}

	return nil
}
