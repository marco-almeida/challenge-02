// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0

package db

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type AssignedOrder struct {
	ID         int64       `json:"id"`
	AssignedAt interface{} `json:"assigned_at"`
	VehicleID  int64       `json:"vehicle_id"`
	OrderID    int64       `json:"order_id"`
}

type Order struct {
	ID           int64        `json:"id"`
	Weight       float32      `json:"weight"`
	Destination  pgtype.Point `json:"destination"`
	Observations string       `json:"observations"`
	Finished     bool         `json:"finished"`
}

type Vehicle struct {
	ID                int64   `json:"id"`
	MaxWeightCapacity float32 `json:"max_weight_capacity"`
	NumberPlate       string  `json:"number_plate"`
}
