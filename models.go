package main

import (
	"time"

	"github.com/google/uuid"
	"github.com/vicentefiorito/rss_feed_aggregator/internal/database"
)

// this stores all of our custom types and handlers

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	ApiKey    string    `json:"api_key"`
}

// helper that converts a database User to a regular User'
func databaseUserToUser(user database.User) User {
	return User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Name:      user.Name,
		ApiKey:    user.ApiKey,
	}
}

type Feed struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	Url       string    `json:"url"`
	UserID    uuid.UUID `json:"user_id"`
}

type FeedFollow struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	UserID    uuid.UUID `json:"user_id"`
	FeedID    uuid.UUID `json:"feed_id"`
}

// handler that converts a database Feed to a regular feed
func databaseFeedToFeed(feed database.Feed) Feed {
	return Feed{
		ID:        feed.ID,
		CreatedAt: feed.CreatedAt,
		UpdatedAt: feed.UpdatedAt,
		Name:      feed.Name,
		Url:       feed.Url,
		UserID:    feed.UserID,
	}
}

// handler that converts more than one database feed
func databaseFeedsToFeeds(feeds []database.Feed) []Feed {
	res := make([]Feed, len(feeds))
	for i, feed := range feeds {
		res[i] = databaseFeedToFeed(feed)
	}
	return res
}

// handler tha converts db feed_follows to a regular one
func databaseFeedFollowToFeedFollow(feedFollow database.FeedFollow) FeedFollow {
	return FeedFollow{
		ID:        feedFollow.ID,
		CreatedAt: feedFollow.CreatedAt,
		UpdatedAt: feedFollow.UpdatedAt,
		UserID:    feedFollow.UserID,
		FeedID:    feedFollow.FeedID,
	}
}

// returns a slice of feed follows for a single user
func databaseFeedFollowsToFeedFollows(feedFollows []database.FeedFollow) []FeedFollow {
	res := make([]FeedFollow, len(feedFollows))
	for i, feedFollow := range feedFollows {
		res[i] = databaseFeedFollowToFeedFollow(feedFollow)
	}
	return res

}
