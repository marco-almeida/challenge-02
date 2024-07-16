package service

import (
	"context"

	"github.com/marco-almeida/challenge-02/internal/postgresql/db"
)

// OrderRepository defines the methods that any User repository should implement.
type OrderRepository interface {
	Create(ctx context.Context, arg db.CreateOrderParams) (db.Order, error)
	Get(ctx context.Context, id int64) (db.Order, error)
	UpdateObservations(ctx context.Context, arg db.UpdateOrderObservationsParams) (db.Order, error)
}

type OrderService struct {
	repo OrderRepository
}

// NewOrderService returns a new OrderService.
func NewOrderService(repo OrderRepository) *OrderService {
	return &OrderService{
		repo: repo,
	}
}

func (s *OrderService) Create(ctx context.Context, arg db.CreateOrderParams) (db.Order, error) {
	return s.repo.Create(ctx, arg)
}

func (s *OrderService) Get(ctx context.Context, id int64) (db.Order, error) {
	return s.repo.Get(ctx, id)
}

func (s *OrderService) UpdateObservations(ctx context.Context, arg db.UpdateOrderObservationsParams) (db.Order, error) {
	return s.repo.UpdateObservations(ctx, arg)
}
