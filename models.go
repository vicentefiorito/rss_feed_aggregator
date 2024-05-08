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
}

// handles the creating of a user in the db
// to be used with json
func dbUserToUser(user database.User) User {
	return User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Name:      user.Name,
	}
}
