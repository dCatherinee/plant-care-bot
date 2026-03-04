package domain

import (
	"fmt"
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
		return Plant{}, fmt.Errorf("%w: userID must be positive", ErrInvalidArgument)
	}
	if name == "" {
		return Plant{}, fmt.Errorf("%w: name is empty", ErrInvalidArgument)
	}

	return Plant{
		UserID:    userID,
		Name:      name,
		CreatedAt: time.Now().UTC(),
	}, nil
}
