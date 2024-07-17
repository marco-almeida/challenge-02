-- name: CreateAssignedOrder :one
INSERT INTO "assigned_order" (
    assigned_at,
    order_id,
    vehicle_id
) VALUES (
    $1, $2, $3
) RETURNING *;

-- name: GetVehicleAssignedOrders :many
SELECT "order".*
FROM "assigned_order"
JOIN "order"
ON "order".id = "assigned_order".order_id
WHERE vehicle_id = $1 AND finished = false
ORDER BY "assigned_order".id
LIMIT $2
OFFSET $3;
