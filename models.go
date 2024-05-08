package main

import (
	"time"

	"github.com/google/uuid"
)

// this stores all of our custom types and handlers

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
}
