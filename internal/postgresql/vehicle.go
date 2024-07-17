package postgresql

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/marco-almeida/challenge-02/internal"
	"github.com/marco-almeida/challenge-02/internal/postgresql/db"
)

// VehicleRepository represents the repository used for interacting with vehicle records.
type VehicleRepository struct {
	q *db.Queries
}

// NewVehicleRepository instantiates the vehicle repository.
func NewVehicleRepository(connPool *pgxpool.Pool) *VehicleRepository {
	return &VehicleRepository{
		q: db.New(connPool),
	}
}

func (vehicleRepo *VehicleRepository) Create(ctx context.Context, arg db.CreateVehicleParams) (db.Vehicle, error) {
	vehicle, err := vehicleRepo.q.CreateVehicle(ctx, arg)
	return vehicle, internal.DBErrorToInternal(err)
}

func (vehicleRepo *VehicleRepository) Get(ctx context.Context, id int64) (db.Vehicle, error) {
	vehicle, err := vehicleRepo.q.GetVehicle(ctx, id)
	return vehicle, internal.DBErrorToInternal(err)
}

func (vehicleRepo *VehicleRepository) GetAll(ctx context.Context, arg db.GetVehiclesParams) ([]db.Vehicle, error) {
	vehicle, err := vehicleRepo.q.GetVehicles(ctx, arg)
	return vehicle, internal.DBErrorToInternal(err)
}

func (vehicleRepo *VehicleRepository) Delete(ctx context.Context, id int64) error {
	err := vehicleRepo.q.DeleteVehicle(ctx, id)
	return internal.DBErrorToInternal(err)
}

func (vehicleRepo *VehicleRepository) UpdateCurrentWeight(ctx context.Context, arg db.UpdateVehicleCurrentWeightParams) (db.Vehicle, error) {
	vehicle, err := vehicleRepo.q.UpdateVehicleCurrentWeight(ctx, arg)
	return vehicle, internal.DBErrorToInternal(err)
}
