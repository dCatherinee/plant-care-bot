package domain

import (
	"strings"
	"time"
)

type Plant struct {
	ID        int64
	UserID    int64
	Name      string
	Notes     string
	CreatedAt time.Time
}

func NewPlant(userID int64, name string) (Plant, error) {
	name = strings.TrimSpace(name)
	if userID <= 0 {
		return Plant{}, ValidationError{Field: "userID", Problem: "must be positive"}
	}
	if name == "" {
		return Plant{}, ValidationError{Field: "name", Problem: "is empty"}
	}

	return Plant{
		UserID:    userID,
		Name:      name,
		CreatedAt: time.Now().UTC(),
	}, nil
}
