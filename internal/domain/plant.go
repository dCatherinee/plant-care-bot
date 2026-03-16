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
	if userID <= 0 {
		return Plant{}, ValidationError{Field: "userID", Problem: "must be positive"}
	}
	normalizedName, err := normalizePlantName(name)
	if err != nil {
		return Plant{}, err
	}

	return Plant{
		UserID:    userID,
		Name:      normalizedName,
		CreatedAt: time.Now().UTC(),
	}, nil
}

func (p *Plant) Rename(name string) error {
	normalizedName, err := normalizePlantName(name)
	if err != nil {
		return err
	}

	p.Name = normalizedName
	return nil
}

func normalizePlantName(name string) (string, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return "", ValidationError{
			Field:   "name",
			Problem: "is empty",
		}
	}
	return name, nil
}
