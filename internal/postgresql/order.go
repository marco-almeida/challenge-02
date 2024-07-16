package postgresql

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/marco-almeida/challenge-02/internal"
	"github.com/marco-almeida/challenge-02/internal/postgresql/db"
)

// OrderRepository represents the repository used for interacting with order records.
type UserRepository struct {
	q *db.Queries
}

// NewOrder instantiates the order repository.
func NewOrderRepository(connPool *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		q: db.New(connPool),
	}
}

func (userRepo *UserRepository) Create(ctx context.Context, arg db.CreateOrderParams) (db.Order, error) {
	order, err := userRepo.q.CreateOrder(ctx, arg)
	return order, internal.DBErrorToInternal(err)
}

func (userRepo *UserRepository) Get(ctx context.Context, id int64) (db.Order, error) {
	order, err := userRepo.q.GetOrder(ctx, id)
	return order, internal.DBErrorToInternal(err)
}

func (userRepo *UserRepository) UpdateObservations(ctx context.Context, arg db.UpdateOrderObservationsParams) (db.Order, error) {
	order, err := userRepo.q.UpdateOrderObservations(ctx, arg)
	return order, internal.DBErrorToInternal(err)
}
