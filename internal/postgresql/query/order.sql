-- name: CreateOrder :one
INSERT INTO "order" (
    weight,
    destination,
    observations,
    finished
) VALUES (
    $1, $2, $3, $4
) RETURNING *;

-- name: GetOrder :one
SELECT *
FROM "order"
WHERE id = $1
LIMIT 1;

-- name: UpdateOrderObservations :one
UPDATE "order"
SET observations = $2
WHERE id = $1
RETURNING *;

-- name: UpdateOrderFinished :one
UPDATE "order"
SET finished = $2
WHERE id = $1
RETURNING *;

-- name: ListOrders :many
SELECT * FROM "order"
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: DeleteOrder :exec
DELETE FROM "order"
WHERE id = $1;
