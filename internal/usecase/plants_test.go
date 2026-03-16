package usecase

import (
	"context"
	"errors"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/dCatherinee/plant-care-bot/internal/domain"
)

func mustNoErr(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func newService() *PlantService {
	return NewPlantService()
}

func addPlant(t *testing.T, svc *PlantService, ctx context.Context, userID int64, name string) domain.Plant {
	t.Helper()
	plant, err := svc.AddPlant(ctx, userID, name)

	mustNoErr(t, err)

	return plant
}

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
			svc := newService()

			plant, err := svc.AddPlant(ctx, tc.userID, tc.plantName)

			if tc.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}

				var myErr domain.ValidationError
				if !errors.As(err, &myErr) {
					t.Fatalf("expected ValidationError, got %T: %v", err, err)
				}
				if myErr.Field != tc.wantField {
					t.Fatalf("expected field %q, got %q", tc.wantField, myErr.Field)
				}
				if myErr.Problem != tc.wantProblem {
					t.Fatalf("expected problem %q, got %q", tc.wantProblem, myErr.Problem)
				}
				if !errors.Is(err, domain.ErrInvalidArgument) {
					t.Fatalf("expected ErrInvalidArgument, got %v", err)
				}
				return
			}

			mustNoErr(t, err)
			if strings.TrimSpace(tc.plantName) != plant.Name {
				t.Fatalf("expected trimmed name %q, got %q", strings.TrimSpace(tc.plantName), plant.Name)
			}
			if plant.ID != 1 {
				t.Fatalf("expected plantID %v, got %v", 1, plant.ID)
			}
		})
	}
}

func TestPlantServiceAddPlantIDIncreases(t *testing.T) {
	ctx := context.Background()
	svc := newService()

	firstPlant := addPlant(t, svc, ctx, 10, "Monstera")
	secondPlant := addPlant(t, svc, ctx, 10, "Cactus")

	if firstPlant.ID != 1 {
		t.Fatalf("expected first plantID %v, got %v", 1, firstPlant.ID)
	}
	if secondPlant.ID != 2 {
		t.Fatalf("expected second plantID %v, got %v", 2, secondPlant.ID)
	}
	if secondPlant.ID <= firstPlant.ID {
		t.Fatalf("expected second plantID to be greater than first: first=%v second=%v", firstPlant.ID, secondPlant.ID)
	}
}

func TestPlantServiceListPlants(t *testing.T) {
	ctx := context.Background()
	svc := newService()

	plant := addPlant(t, svc, ctx, 10, "Monstera")

	list, err := svc.ListPlants(ctx, 10)
	mustNoErr(t, err)

	if len(list) != 1 {
		t.Fatalf("expected list length %v, got %v", 1, len(list))
	}
	if list[0].ID != plant.ID {
		t.Fatalf("expected plantID %v, got %v", plant.ID, list[0].ID)
	}
	if list[0].Name != plant.Name {
		t.Fatalf("expected plant name %q, got %q", plant.Name, list[0].Name)
	}
}

func TestPlantServiceListPlantsReturnsCopy(t *testing.T) {
	ctx := context.Background()
	svc := newService()

	addPlant(t, svc, ctx, 10, "Monstera")

	list, err := svc.ListPlants(ctx, 10)
	mustNoErr(t, err)

	list[0].Name = "Changed"

	freshList, err := svc.ListPlants(ctx, 10)
	mustNoErr(t, err)

	if freshList[0].Name != "Monstera" {
		t.Fatalf("expected original plant name %q, got %q", "Monstera", freshList[0].Name)
	}
}

func TestPlantServiceListPlantsOtherUserEmpty(t *testing.T) {
	ctx := context.Background()
	svc := newService()

	addPlant(t, svc, ctx, 10, "Monstera")

	list, err := svc.ListPlants(ctx, 20)
	mustNoErr(t, err)

	if len(list) != 0 {
		t.Fatalf("expected empty list for another user, got %v items", len(list))
	}
}

func TestPlantServiceDeletePlant(t *testing.T) {
	ctx := context.Background()
	svc := newService()

	plant := addPlant(t, svc, ctx, 10, "Cactus")

	err := svc.DeletePlant(ctx, 10, plant.ID)
	mustNoErr(t, err)

	list, err := svc.ListPlants(ctx, 10)
	mustNoErr(t, err)

	if len(list) != 0 {
		t.Fatalf("expected empty list after delete, got %v items", len(list))
	}
}

func TestPlantServiceDeletePlantOtherUserReturnsNotFound(t *testing.T) {
	ctx := context.Background()
	svc := newService()

	plant := addPlant(t, svc, ctx, 10, "Cactus")

	err := svc.DeletePlant(ctx, 20, plant.ID)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestPlantServiceCancelContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	svc := newService()

	_, err := svc.AddPlant(ctx, 10, "Cactus")

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context.Canceled, got %v", err)
	}
}
func TestPlantServiceConcurrentAddPlant(t *testing.T) {
	svc := newService()
	ctx := context.Background()
	const N = 100
	start := make(chan struct{})
	errCh := make(chan error, N)
	var wg sync.WaitGroup

	wg.Add(N)
	for i := 1; i <= N; i++ {
		go func() {
			defer wg.Done()
			<-start

			_, err := svc.AddPlant(ctx, 10, "Monstera")

			if err != nil {
				errCh <- err
			}
		}()
	}
	close(start)
	wg.Wait()
	close(errCh)

	for err := range errCh {
		t.Fatalf("AddPlant returned error: %v", err)
	}

	list, err := svc.ListPlants(ctx, 10)
	mustNoErr(t, err)

	if len(list) != N {
		t.Fatalf("expected %d plants, got %d", N, len(list))
	}

	seen := map[int64]struct{}{}

	for _, plant := range list {
		if _, ok := seen[plant.ID]; ok {
			t.Fatalf("duplicate plant ID: %d", plant.ID)
		}
		seen[plant.ID] = struct{}{}
	}
}

func TestPlantServiceDeadlineExceeded(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()
	svc := newService()

	_, err := svc.AddPlant(ctx, 10, "Cactus")

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("expected context.DeadlineExceeded, got %v", err)
	}
}

func TestPlantServiceGetPlantOk(t *testing.T) {
	ctx := context.Background()
	svc := newService()

	plant := addPlant(t, svc, ctx, 10, "Cactus")
	getPlant, err := svc.GetPlant(ctx, 10, plant.ID)
	mustNoErr(t, err)

	if getPlant.ID != plant.ID || getPlant.Name != plant.Name {
		t.Fatalf("expected plant with ID %d, got ID %d", plant.ID, getPlant.ID)
	}
}

func TestPlantServiceGetPlantNotFound(t *testing.T) {
	ctx := context.Background()
	svc := newService()

	addPlant(t, svc, ctx, 10, "Cactus")
	_, err := svc.GetPlant(ctx, 10, 3)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestPlantServiceGetPlantWrongUserID(t *testing.T) {
	ctx := context.Background()
	svc := newService()

	plant := addPlant(t, svc, ctx, 10, "Cactus")
	_, err := svc.GetPlant(ctx, 5, plant.ID)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestPlantServiceUpdatePlantName(t *testing.T) {
	tests := []struct {
		name      string
		userID    int64
		plantID   int64
		plantName string
		wantName  string
		wantErr   error
		wantField string
		wantProb  string
	}{
		{
			name:      "update_plant_name_ok",
			userID:    10,
			plantID:   1,
			plantName: "Cactus",
			wantName:  "Cactus",
		},
		{
			name:      "plant_not_found",
			userID:    10,
			plantID:   3,
			plantName: "Cactus",
			wantErr:   domain.ErrNotFound,
		},
		{
			name:      "wrong_user",
			userID:    17,
			plantID:   1,
			plantName: "Cactus",
			wantErr:   domain.ErrNotFound,
		},
		{
			name:      "empty_name",
			userID:    10,
			plantID:   1,
			plantName: "",
			wantErr:   domain.ErrInvalidArgument,
			wantField: "name",
			wantProb:  "is empty",
		},
		{
			name:      "trim_name",
			userID:    10,
			plantID:   1,
			plantName: "  Cactus ",
			wantName:  "Cactus",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			svc := newService()

			plant, err := svc.AddPlant(ctx, 10, "Monstera")
			mustNoErr(t, err)

			plant, err = svc.UpdatePlantName(ctx, tc.userID, tc.plantID, tc.plantName)

			if tc.wantErr != nil {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				if !errors.Is(err, tc.wantErr) {
					t.Fatalf("expected %v, got %v", tc.wantErr, err)
				}
				if tc.wantField != "" {
					var validationErr domain.ValidationError
					if !errors.As(err, &validationErr) {
						t.Fatalf("expected ValidationError, got %T: %v", err, err)
					}
					if validationErr.Field != tc.wantField {
						t.Fatalf("expected field %q, got %q", tc.wantField, validationErr.Field)
					}
					if validationErr.Problem != tc.wantProb {
						t.Fatalf("expected problem %q, got %q", tc.wantProb, validationErr.Problem)
					}
				}
				return
			}

			mustNoErr(t, err)
			if tc.wantName != plant.Name {
				t.Fatalf("expected new name %q, got %q", tc.wantName, plant.Name)
			}
		})
	}
}

func TestPlantServiceUpdatePlantNameCanceledContext(t *testing.T) {
	ctx := context.Background()
	svc := newService()

	plant := addPlant(t, svc, ctx, 10, "Monstera")

	canceledCtx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := svc.UpdatePlantName(canceledCtx, 10, plant.ID, "Cactus")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context.Canceled, got %v", err)
	}
}
