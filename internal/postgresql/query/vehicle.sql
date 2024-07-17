-- name: CreateVehicle :one
INSERT INTO "vehicle" (
    max_weight_capacity,
    number_plate
) VALUES (
    $1, $2
) RETURNING *;

-- name: GetVehicle :one
SELECT *
FROM "vehicle"
WHERE id = $1
LIMIT 1;

-- name: GetVehicles :many
SELECT *
FROM "vehicle"
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: DeleteVehicle :exec
DELETE FROM "vehicle"
WHERE id = $1;

-- name: UpdateVehicleCurrentWeight :one
UPDATE "vehicle"
SET current_weight = $1
WHERE id = $2
RETURNING *;