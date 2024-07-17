// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: assigned_order.sql

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createAssignedOrder = `-- name: CreateAssignedOrder :one
INSERT INTO "assigned_order" (
    assigned_at,
    order_id,
    vehicle_id
) VALUES (
    $1, $2, $3
) RETURNING id, assigned_at, vehicle_id, order_id
`

type CreateAssignedOrderParams struct {
	AssignedAt pgtype.Timestamp `json:"assigned_at"`
	OrderID    int64            `json:"order_id"`
	VehicleID  int64            `json:"vehicle_id"`
}

func (q *Queries) CreateAssignedOrder(ctx context.Context, arg CreateAssignedOrderParams) (AssignedOrder, error) {
	row := q.db.QueryRow(ctx, createAssignedOrder, arg.AssignedAt, arg.OrderID, arg.VehicleID)
	var i AssignedOrder
	err := row.Scan(
		&i.ID,
		&i.AssignedAt,
		&i.VehicleID,
		&i.OrderID,
	)
	return i, err
}

const getVehicleAssignedOrders = `-- name: GetVehicleAssignedOrders :many
SELECT "order".id, "order".weight, "order".destination, "order".observations, "order".finished
FROM "assigned_order"
JOIN "order"
ON "order".id = "assigned_order".order_id
WHERE vehicle_id = $1 AND finished = false
ORDER BY "assigned_order".id
LIMIT $2
OFFSET $3
`

type GetVehicleAssignedOrdersParams struct {
	VehicleID int64 `json:"vehicle_id"`
	Limit     int32 `json:"limit"`
	Offset    int32 `json:"offset"`
}

func (q *Queries) GetVehicleAssignedOrders(ctx context.Context, arg GetVehicleAssignedOrdersParams) ([]Order, error) {
	rows, err := q.db.Query(ctx, getVehicleAssignedOrders, arg.VehicleID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Order{}
	for rows.Next() {
		var i Order
		if err := rows.Scan(
			&i.ID,
			&i.Weight,
			&i.Destination,
			&i.Observations,
			&i.Finished,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
