package domain

import (
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
		return CareEvent{}, ValidationError{Field: "plantID", Problem: "must be positive"}
	}
	if !kind.Valid() {
		return CareEvent{}, ValidationError{Field: "kind", Problem: "invalid care kind"}
	}
	if occurredAt.IsZero() {
		return CareEvent{}, ValidationError{Field: "occurredAt", Problem: "is zero"}
	}

	return CareEvent{
		PlantID:    plantID,
		Kind:       kind,
		OccurredAt: occurredAt.UTC(),
		CreatedAt:  time.Now().UTC(),
	}, nil
}

func MustCareEvent(id, plantID int64, kind CareKind, occurredAt time.Time, createdAt time.Time) CareEvent {
	res, err := NewCareEvent(plantID, kind, occurredAt)
	if err != nil {
		panic(err)
	}

	res.ID = id
	res.CreatedAt = createdAt

	return res
}
