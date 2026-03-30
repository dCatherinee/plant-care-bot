package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/dCatherinee/plant-care-bot/internal/domain"
)

type fakePlantRepo struct {
	createPlantFn      func(ctx context.Context, plant domain.Plant) (int64, error)
	listPlantsByUserFn func(ctx context.Context, userID int64) ([]domain.Plant, error)
	deletePlantFn      func(ctx context.Context, userID int64, plantID int64) error
	getPlantByIDFn     func(ctx context.Context, userID int64, plantID int64) (domain.Plant, error)
	updatePlantNameFn  func(ctx context.Context, userID int64, plantID int64, name string) (domain.Plant, error)
}

func (f *fakePlantRepo) CreatePlant(ctx context.Context, plant domain.Plant) (int64, error) {
	if f.createPlantFn != nil {
		return f.createPlantFn(ctx, plant)
	}
	return 0, nil
}

func (f *fakePlantRepo) ListPlantsByUser(ctx context.Context, userID int64) ([]domain.Plant, error) {
	if f.listPlantsByUserFn != nil {
		return f.listPlantsByUserFn(ctx, userID)
	}
	return nil, nil
}

func (f *fakePlantRepo) DeletePlant(ctx context.Context, userID int64, plantID int64) error {
	if f.deletePlantFn != nil {
		return f.deletePlantFn(ctx, userID, plantID)
	}
	return nil
}

func (f *fakePlantRepo) GetPlantByID(ctx context.Context, userID int64, plantID int64) (domain.Plant, error) {
	if f.getPlantByIDFn != nil {
		return f.getPlantByIDFn(ctx, userID, plantID)
	}
	return domain.Plant{}, nil
}

func (f *fakePlantRepo) UpdatePlantName(ctx context.Context, userID int64, plantID int64, name string) (domain.Plant, error) {
	if f.updatePlantNameFn != nil {
		return f.updatePlantNameFn(ctx, userID, plantID, name)
	}
	return domain.Plant{}, nil
}

func mustNoErr(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func newService(r *fakePlantRepo) *PlantService {
	return NewPlantService(r)
}

// func addPlant(t *testing.T, svc *PlantService, ctx context.Context, userID int64, name string) domain.Plant {
// 	t.Helper()
// 	plant, err := svc.AddPlant(ctx, userID, name)

// 	mustNoErr(t, err)

// 	return plant
// }

func TestPlantServiceAddPlant(t *testing.T) {
	repo := &fakePlantRepo{
		createPlantFn: func(ctx context.Context, plant domain.Plant) (int64, error) {
			if plant.UserID != 10 {
				t.Fatalf("expected userID 10, got %d", plant.UserID)
			}
			if plant.Name != "Cactus" {
				t.Fatalf("expected trimmed name %q, got %q", "Cactus", plant.Name)
			}
			return 42, nil
		},
	}

	svc := newService(repo)

	plant, err := svc.AddPlant(context.Background(), 10, "  Cactus  ")
	mustNoErr(t, err)

	if plant.ID != 42 {
		t.Fatalf("expected plant ID 42, got %d", plant.ID)
	}
	if plant.Name != "Cactus" {
		t.Fatalf("expected name %q, got %q", "Cactus", plant.Name)
	}
}

func TestPlantServiceAddPlantValidationError(t *testing.T) {
	repo := &fakePlantRepo{
		createPlantFn: func(ctx context.Context, plant domain.Plant) (int64, error) {
			t.Fatal("repo should not be called on invalid input")
			return 0, nil
		},
	}

	svc := newService(repo)

	_, err := svc.AddPlant(context.Background(), 0, "")
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	var vErr domain.ValidationError
	if !errors.As(err, &vErr) {
		t.Fatalf("expected ValidationError, got %T: %v", err, err)
	}
	if !errors.Is(err, domain.ErrInvalidArgument) {
		t.Fatalf("expected ErrInvalidArgument, got %v", err)
	}
}

func TestPlantServiceListPlants(t *testing.T) {
	repo := &fakePlantRepo{
		listPlantsByUserFn: func(ctx context.Context, userID int64) ([]domain.Plant, error) {
			if userID != 10 {
				t.Fatalf("expected userID 10, got %d", userID)
			}
			return []domain.Plant{
				{ID: 1, UserID: 10, Name: "Monstera"},
			}, nil
		},
	}

	svc := newService(repo)

	list, err := svc.ListPlants(context.Background(), 10)
	mustNoErr(t, err)

	if len(list) != 1 {
		t.Fatalf("expected 1 plant, got %d", len(list))
	}
	if list[0].Name != "Monstera" {
		t.Fatalf("expected %q, got %q", "Monstera", list[0].Name)
	}
}

func TestPlantServiceListPlantsReturnsCopy(t *testing.T) {
	repo := &fakePlantRepo{
		listPlantsByUserFn: func(ctx context.Context, userID int64) ([]domain.Plant, error) {
			return []domain.Plant{
				{ID: 1, UserID: 10, Name: "Monstera"},
			}, nil
		},
	}

	svc := newService(repo)

	list, err := svc.ListPlants(context.Background(), 10)
	mustNoErr(t, err)

	list[0].Name = "Changed"

	freshList, err := svc.ListPlants(context.Background(), 10)
	mustNoErr(t, err)

	if freshList[0].Name != "Monstera" {
		t.Fatalf("expected original plant name %q, got %q", "Monstera", freshList[0].Name)
	}
}

func TestPlantServiceDeletePlant(t *testing.T) {
	repo := &fakePlantRepo{
		deletePlantFn: func(ctx context.Context, userID int64, plantID int64) error {
			if userID != 10 {
				t.Fatalf("expected userID 10, got %d", userID)
			}
			if plantID != 5 {
				t.Fatalf("expected plantID 5, got %d", plantID)
			}
			return nil
		},
	}

	svc := newService(repo)

	err := svc.DeletePlant(context.Background(), 10, 5)
	mustNoErr(t, err)
}

func TestPlantServiceDeletePlantOtherUserReturnsNotFound(t *testing.T) {
	repo := &fakePlantRepo{
		deletePlantFn: func(ctx context.Context, userID int64, plantID int64) error {
			return domain.ErrNotFound
		},
	}

	svc := newService(repo)

	err := svc.DeletePlant(context.Background(), 20, 1)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestPlantServiceGetPlantOk(t *testing.T) {
	repo := &fakePlantRepo{
		getPlantByIDFn: func(ctx context.Context, userID int64, plantID int64) (domain.Plant, error) {
			return domain.Plant{ID: 1, UserID: 10, Name: "Monstera"}, nil
		},
	}
	ctx := context.Background()
	svc := newService(repo)

	getPlant, err := svc.GetPlant(ctx, 10, 1)
	mustNoErr(t, err)

	if getPlant.ID != 1 || getPlant.Name != "Monstera" {
		t.Fatalf("expected plant with ID %d, got ID %d", 1, getPlant.ID)
	}
}

func TestPlantServiceGetPlantNotFound(t *testing.T) {
	repo := &fakePlantRepo{
		getPlantByIDFn: func(ctx context.Context, userID int64, plantID int64) (domain.Plant, error) {
			return domain.Plant{}, domain.ErrNotFound
		},
	}
	ctx := context.Background()
	svc := newService(repo)

	_, err := svc.GetPlant(ctx, 10, 3)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestPlantServiceGetPlantWrongUserID(t *testing.T) {
	repo := &fakePlantRepo{
		getPlantByIDFn: func(ctx context.Context, userID int64, plantID int64) (domain.Plant, error) {
			if userID != 10 {
				return domain.Plant{}, domain.ErrNotFound
			}
			return domain.Plant{ID: 1, UserID: 10, Name: "Monstera"}, nil
		},
	}
	ctx := context.Background()
	svc := newService(repo)

	_, err := svc.GetPlant(ctx, 5, 1)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestPlantServiceUpdatePlantName(t *testing.T) {
	dbErr := errors.New("db failed")

	tests := []struct {
		name           string
		userID         int64
		plantID        int64
		plantName      string
		repoPlant      domain.Plant
		repoErr        error
		wantPlant      domain.Plant
		wantErr        error
		wantRepoCalled bool
		wantRepoName   string
	}{
		{
			name:           "update_plant_name_ok",
			userID:         10,
			plantID:        1,
			plantName:      "Cactus",
			repoPlant:      domain.Plant{ID: 1, UserID: 10, Name: "Cactus"},
			wantPlant:      domain.Plant{ID: 1, UserID: 10, Name: "Cactus"},
			wantRepoCalled: true,
		},
		{
			name:           "plant_not_found",
			userID:         10,
			plantID:        3,
			plantName:      "Cactus",
			repoErr:        domain.ErrNotFound,
			wantErr:        domain.ErrNotFound,
			wantRepoCalled: true,
		},
		{
			name:           "repo_error",
			userID:         10,
			plantID:        1,
			plantName:      "Cactus",
			repoErr:        dbErr,
			wantErr:        dbErr,
			wantRepoCalled: true,
		},
		{
			name:           "empty_name",
			userID:         10,
			plantID:        1,
			plantName:      "",
			wantErr:        domain.ErrInvalidArgument,
			wantRepoCalled: false,
		},
		{
			name:           "trim_name",
			userID:         10,
			plantID:        1,
			plantName:      "  Cactus   ",
			repoPlant:      domain.Plant{ID: 1, UserID: 10, Name: "Cactus"},
			wantPlant:      domain.Plant{ID: 1, UserID: 10, Name: "Cactus"},
			wantRepoCalled: true,
			wantRepoName:   "Cactus",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			repoCalled := false
			repo := &fakePlantRepo{
				updatePlantNameFn: func(ctx context.Context, userID int64, plantID int64, name string) (domain.Plant, error) {
					repoCalled = true

					if userID != tc.userID {
						t.Fatalf("expected userID %d, got %d", tc.userID, userID)
					}
					if plantID != tc.plantID {
						t.Fatalf("expected plantID %d, got %d", tc.plantID, plantID)
					}
					if tc.wantRepoName != "" && name != tc.wantRepoName {
						t.Fatalf("expected repo name %q, got %q", tc.wantRepoName, name)
					}

					return tc.repoPlant, tc.repoErr
				},
			}

			svc := newService(repo)

			plant, err := svc.UpdatePlantName(context.Background(), tc.userID, tc.plantID, tc.plantName)

			if tc.wantErr != nil {
				if repoCalled != tc.wantRepoCalled {
					t.Fatalf("expected repo called=%v, got %v", tc.wantRepoCalled, repoCalled)
				}
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				if !errors.Is(err, tc.wantErr) {
					t.Fatalf("expected %v, got %v", tc.wantErr, err)
				}
				return
			}

			mustNoErr(t, err)
			if repoCalled != tc.wantRepoCalled {
				t.Fatalf("expected repo called=%v, got %v", tc.wantRepoCalled, repoCalled)
			}
			if plant != tc.wantPlant {
				t.Fatalf("expected plant %+v, got %+v", tc.wantPlant, plant)
			}
		})
	}
}

func TestPlantServiceUpdatePlantNameCanceledContext(t *testing.T) {
	canceledCtx, cancel := context.WithCancel(context.Background())
	cancel()

	repoCalled := false
	repo := &fakePlantRepo{
		updatePlantNameFn: func(ctx context.Context, userID int64, plantID int64, name string) (domain.Plant, error) {
			repoCalled = true
			return domain.Plant{}, nil
		},
	}

	svc := newService(repo)

	_, err := svc.UpdatePlantName(canceledCtx, 10, 1, "Cactus")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context.Canceled, got %v", err)
	}
	if repoCalled {
		t.Fatal("repo should not be called when context is canceled")
	}
}
