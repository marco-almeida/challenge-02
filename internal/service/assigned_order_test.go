package service

import (
	"strconv"
	"testing"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/marco-almeida/challenge-02/internal/postgresql/db"
)

func TestSortOrdersByClosest(t *testing.T) {
	t.Parallel()
	var orders []db.Order
	for i := 0; i < 5; i++ {
		orders = append(orders, db.Order{
			ID: int64(i),
			Destination: pgtype.Point{
				P: pgtype.Vec2{
					X: float64(i),
					Y: float64(i),
				},
				Valid: true,
			},
		})
	}

	for _, order := range orders {
		t.Run("Order"+strconv.FormatInt(order.ID, 10)+"ShouldBeClosest", func(t *testing.T) {
			sortOrdersByClosest(orders, order.Destination.P.X, order.Destination.P.Y)

			if orders[0].ID != order.ID {
				t.Errorf("expected order ID %d to be the closest, got %d", order.ID, orders[0].ID)
			}
		})
	}

	if len(orders) != 5 {
		t.Errorf("expected 5 orders, got %d", len(orders))
	}

}
