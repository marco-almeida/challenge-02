package postgresql

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/marco-almeida/challenge-02/internal"
	"github.com/marco-almeida/challenge-02/internal/postgresql/db"
)

// OrderRepository represents the repository used for interacting with order records.
type OrderRepository struct {
	q *db.Queries
}

// NewOrder instantiates the order repository.
func NewOrderRepository(connPool *pgxpool.Pool) *OrderRepository {
	return &OrderRepository{
		q: db.New(connPool),
	}
}

func (orderRepo *OrderRepository) Create(ctx context.Context, arg db.CreateOrderParams) (db.Order, error) {
	order, err := orderRepo.q.CreateOrder(ctx, arg)
	return order, internal.DBErrorToInternal(err)
}

func (orderRepo *OrderRepository) Get(ctx context.Context, id int64) (db.Order, error) {
	order, err := orderRepo.q.GetOrder(ctx, id)
	return order, internal.DBErrorToInternal(err)
}

func (orderRepo *OrderRepository) UpdateObservations(ctx context.Context, arg db.UpdateOrderObservationsParams) (db.Order, error) {
	order, err := orderRepo.q.UpdateOrderObservations(ctx, arg)
	return order, internal.DBErrorToInternal(err)
}
