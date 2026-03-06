package domain

import (
	"fmt"
	"time"
)

type CareKind string

const (
	CareKindWater     CareKind = "water"
	CareKindFertilize CareKind = "fertilize"
)

func (k CareKind) Valid() bool {
	switch k {
	case CareKindWater, CareKindFertilize:
		return true
	default:
		return false
	}
}

type CareEvent struct {
	ID         int64
	PlantID    int64
	Kind       CareKind
	OccurredAt time.Time
	CreatedAt  time.Time
}

func NewCareEvent(plantID int64, kind CareKind, occurredAt time.Time) (CareEvent, error) {
	if plantID <= 0 {
		return CareEvent{}, fmt.Errorf("%w: plantID must be positive", ErrInvalidArgument)
	}
	if !kind.Valid() {
		return CareEvent{}, fmt.Errorf("%w: invalid care kind %q", ErrInvalidArgument, kind)
	}
	if occurredAt.IsZero() {
		return CareEvent{}, fmt.Errorf("%w: occurredAt is zero", ErrInvalidArgument)
	}

	return CareEvent{
		PlantID:    plantID,
		Kind:       kind,
		OccurredAt: occurredAt.UTC(),
		CreatedAt:  time.Now().UTC(),
	}, nil
}
