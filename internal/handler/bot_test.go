package handler

import (
	"context"
	"errors"
	"testing"

	"github.com/dCatherinee/plant-care-bot/internal/domain"
)

type plantUsecaseStub struct {
	addPlantFn func(ctx context.Context, userID int64, name string) (domain.Plant, error)
	listFn     func(ctx context.Context, userID int64) ([]domain.Plant, error)
}

func (s plantUsecaseStub) AddPlant(ctx context.Context, userID int64, name string) (domain.Plant, error) {
	if s.addPlantFn == nil {
		return domain.Plant{}, nil
	}

	return s.addPlantFn(ctx, userID, name)
}

func (s plantUsecaseStub) ListPlants(ctx context.Context, userID int64) ([]domain.Plant, error) {
	if s.listFn == nil {
		return nil, nil
	}

	return s.listFn(ctx, userID)
}

func (s plantUsecaseStub) GetPlant(ctx context.Context, userID int64, plantID int64) (domain.Plant, error) {
	return domain.Plant{}, errors.New("not implemented")
}

func (s plantUsecaseStub) UpdatePlantName(ctx context.Context, userID int64, plantID int64, name string) (domain.Plant, error) {
	return domain.Plant{}, errors.New("not implemented")
}

func (s plantUsecaseStub) DeletePlant(ctx context.Context, userID int64, plantID int64) error {
	return errors.New("not implemented")
}

func TestBotHandlerHandleTextStart(t *testing.T) {
	handler := NewBotHandler(plantUsecaseStub{})

	got := handler.HandleText(context.Background(), 42, "/start")

	if got != defaultHelpText {
		t.Fatalf("expected help text %q, got %q", defaultHelpText, got)
	}
}

func TestBotHandlerHandleTextAddPlant(t *testing.T) {
	ctx := context.Background()
	var gotUserID int64
	var gotName string

	handler := NewBotHandler(plantUsecaseStub{
		addPlantFn: func(ctx context.Context, userID int64, name string) (domain.Plant, error) {
			gotUserID = userID
			gotName = name

			return domain.Plant{
				ID:     1,
				UserID: userID,
				Name:   name,
			}, nil
		},
	})

	reply := handler.HandleText(ctx, 7, "добавить Monstera")

	if gotUserID != 7 {
		t.Fatalf("expected AddPlant userID %d, got %d", 7, gotUserID)
	}
	if gotName != "Monstera" {
		t.Fatalf("expected AddPlant name %q, got %q", "Monstera", gotName)
	}
	if reply != "Добавила растение: Monstera" {
		t.Fatalf("expected add reply %q, got %q", "Добавила растение: Monstera", reply)
	}
}

func TestBotHandlerHandleTextAddPlantEmptyName(t *testing.T) {
	handler := NewBotHandler(plantUsecaseStub{})

	reply := handler.HandleText(context.Background(), 7, "добавить   ")

	if reply != addUsageText {
		t.Fatalf("expected reply %q, got %q", addUsageText, reply)
	}
}

func TestBotHandlerHandleTextListPlants(t *testing.T) {
	ctx := context.Background()
	var gotUserID int64

	handler := NewBotHandler(plantUsecaseStub{
		listFn: func(ctx context.Context, userID int64) ([]domain.Plant, error) {
			gotUserID = userID

			return []domain.Plant{
				{ID: 1, UserID: userID, Name: "Monstera"},
				{ID: 2, UserID: userID, Name: "Cactus"},
			}, nil
		},
	})

	reply := handler.HandleText(ctx, 9, "список")

	if gotUserID != 9 {
		t.Fatalf("expected ListPlants userID %d, got %d", 9, gotUserID)
	}

	want := "Твои растения:\n1. Monstera\n2. Cactus"
	if reply != want {
		t.Fatalf("expected list reply %q, got %q", want, reply)
	}
}

func TestBotHandlerHandleTextListPlantsEmpty(t *testing.T) {
	handler := NewBotHandler(plantUsecaseStub{
		listFn: func(ctx context.Context, userID int64) ([]domain.Plant, error) {
			return []domain.Plant{}, nil
		},
	})

	reply := handler.HandleText(context.Background(), 9, "список")

	if reply != "Список растений пуст." {
		t.Fatalf("expected empty list reply %q, got %q", "Список растений пуст.", reply)
	}
}
