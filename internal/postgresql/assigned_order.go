package postgresql

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/marco-almeida/challenge-02/internal"
	"github.com/marco-almeida/challenge-02/internal/postgresql/db"
)

// AssignedOrderRepository represents the repository used for interacting with assigned order records.
type AssignedOrderRepository struct {
	q *db.Queries
}

// NewAssignedOrderRepository instantiates the assigned order repository.
func NewAssignedOrderRepository(connPool *pgxpool.Pool) *AssignedOrderRepository {
	return &AssignedOrderRepository{
		q: db.New(connPool),
	}
}

func (assignedOrderRepo *AssignedOrderRepository) GetVehicleAssignedOrders(ctx context.Context, arg db.GetVehicleAssignedOrdersParams) ([]db.Order, error) {
	assignedOrder, err := assignedOrderRepo.q.GetVehicleAssignedOrders(ctx, arg)
	return assignedOrder, internal.DBErrorToInternal(err)
}

func (assignedOrderRepo *AssignedOrderRepository) CreateAssignedOrder(ctx context.Context, arg db.CreateAssignedOrderParams) (db.AssignedOrder, error) {
	assignedOrder, err := assignedOrderRepo.q.CreateAssignedOrder(ctx, arg)
	return assignedOrder, internal.DBErrorToInternal(err)
}
