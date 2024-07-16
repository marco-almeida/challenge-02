package service

import (
	"context"
	"errors"

	"github.com/marco-almeida/challenge-02/internal"
	"github.com/marco-almeida/challenge-02/internal/postgresql/db"
)

// VehicleRepository defines the methods that any vehicle repository should implement.
type VehicleRepository interface {
	Create(ctx context.Context, arg db.CreateVehicleParams) (db.Vehicle, error)
	Get(ctx context.Context, id int64) (db.Vehicle, error)
	Delete(ctx context.Context, id int64) error
	GetAll(ctx context.Context, arg db.GetVehiclesParams) ([]db.Vehicle, error)
}

type VehicleService struct {
	repo VehicleRepository
}

// NewVehicleService returns a new VehicleService.
func NewVehicleService(repo VehicleRepository) *VehicleService {
	return &VehicleService{
		repo: repo,
	}
}

func (s *VehicleService) Create(ctx context.Context, arg db.CreateVehicleParams) (db.Vehicle, error) {
	vehicle, err := s.repo.Create(ctx, arg)
	if errors.Is(err, internal.ErrUniqueConstraintViolation) {
		return db.Vehicle{}, internal.ErrVehicleAlreadyExists
	}

	return vehicle, err
}

func (s *VehicleService) Get(ctx context.Context, id int64) (db.Vehicle, error) {
	return s.repo.Get(ctx, id)
}

func (s *VehicleService) GetAll(ctx context.Context, arg db.GetVehiclesParams) ([]db.Vehicle, error) {
	return s.repo.GetAll(ctx, arg)
}

func (s *VehicleService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}
