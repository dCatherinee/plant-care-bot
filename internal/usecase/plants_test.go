package usecase

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/dCatherinee/plant-care-bot/internal/domain"
)

func TestPlantServiceAddPlant(t *testing.T) {
	tests := []struct {
		name        string
		userID      int64
		plantName   string
		wantErr     bool
		wantField   string
		wantProblem string
	}{
		{"add_plant_ok", 10, "Monstera", false, "", ""},
		{"empty_user_id", 0, "Cactus", true, "userID", "must be positive"},
		{"empty_name", 10, "", true, "name", "is empty"},
		{"trim_name", 10, " Cactus Poppy  ", false, "", ""},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			plantService := NewPlantService()

			plant, err := plantService.AddPlant(ctx, tc.userID, tc.plantName)

			if tc.wantErr {
				if err == nil {
					t.Fatal("Expected error, got nil")
				}

				var myErr domain.ValidationError
				if !errors.As(err, &myErr) {
					t.Fatalf("Expected ValidationError, got %T: %v", err, err)
				}
				if myErr.Field != tc.wantField {
					t.Fatalf("Expected field %q, got %q", tc.wantField, myErr.Field)
				}
				if myErr.Problem != tc.wantProblem {
					t.Fatalf("Expected problem %q, got %q", tc.wantProblem, myErr.Problem)
				}
				if !errors.Is(err, domain.ErrInvalidArgument) {
					t.Fatalf("Expected ErrInvalidArgument, got %v", err)
				}
				return
			}

			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}
			if strings.TrimSpace(tc.plantName) != plant.Name {
				t.Fatalf("Expected trimmed name %q, got %q", strings.TrimSpace(tc.plantName), plant.Name)
			}
			if plant.ID != 1 {
				t.Fatalf("Expected plantID %v, got %v", 1, plant.ID)
			}
		})
	}
}

func TestPlantServiceAddPlantIDIncreases(t *testing.T) {
	ctx := context.Background()
	plantService := NewPlantService()

	firstPlant, err := plantService.AddPlant(ctx, 10, "Monstera")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	secondPlant, err := plantService.AddPlant(ctx, 10, "Cactus")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if firstPlant.ID != 1 {
		t.Fatalf("Expected first plantID %v, got %v", 1, firstPlant.ID)
	}
	if secondPlant.ID != 2 {
		t.Fatalf("Expected second plantID %v, got %v", 2, secondPlant.ID)
	}
	if secondPlant.ID <= firstPlant.ID {
		t.Fatalf("Expected second plantID to be greater than first: first=%v second=%v", firstPlant.ID, secondPlant.ID)
	}
}

func TestPlantServiceListPlants(t *testing.T) {
	ctx := context.Background()
	plantService := NewPlantService()

	plant, err := plantService.AddPlant(ctx, 10, "Monstera")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	list, err := plantService.ListPlants(ctx, 10)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(list) != 1 {
		t.Fatalf("Expected list length %v, got %v", 1, len(list))
	}
	if list[0].ID != plant.ID {
		t.Fatalf("Expected plantID %v, got %v", plant.ID, list[0].ID)
	}
	if list[0].Name != plant.Name {
		t.Fatalf("Expected plant name %q, got %q", plant.Name, list[0].Name)
	}
}

func TestPlantServiceListPlantsReturnsCopy(t *testing.T) {
	ctx := context.Background()
	plantService := NewPlantService()

	_, err := plantService.AddPlant(ctx, 10, "Monstera")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	list, err := plantService.ListPlants(ctx, 10)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	list[0].Name = "Changed"

	freshList, err := plantService.ListPlants(ctx, 10)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if freshList[0].Name != "Monstera" {
		t.Fatalf("Expected original plant name %q, got %q", "Monstera", freshList[0].Name)
	}
}

func TestPlantServiceListPlantsOtherUserEmpty(t *testing.T) {
	ctx := context.Background()
	plantService := NewPlantService()

	_, err := plantService.AddPlant(ctx, 10, "Monstera")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	list, err := plantService.ListPlants(ctx, 20)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(list) != 0 {
		t.Fatalf("Expected empty list for another user, got %v items", len(list))
	}
}

func TestPlantServiceDeletePlant(t *testing.T) {
	ctx := context.Background()
	plantService := NewPlantService()

	plant, err := plantService.AddPlant(ctx, 10, "Cactus")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	err = plantService.DeletePlant(ctx, 10, plant.ID)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	list, err := plantService.ListPlants(ctx, 10)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(list) != 0 {
		t.Fatalf("Expected empty list after delete, got %v items", len(list))
	}
}

func TestPlantServiceDeletePlantOtherUserReturnsNotFound(t *testing.T) {
	ctx := context.Background()
	plantService := NewPlantService()

	plant, err := plantService.AddPlant(ctx, 10, "Cactus")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	err = plantService.DeletePlant(ctx, 20, plant.ID)
	if err == nil {
		t.Fatal("Expected error, got nil")
	}
	if !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("Expected ErrNotFound, got %v", err)
	}
}
