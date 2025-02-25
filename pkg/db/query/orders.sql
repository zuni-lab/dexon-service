-- name: InsertOrder :one
INSERT INTO orders (parent_id, wallet, pool_id, side, status, type, price, amount, twap_total_time, filled_at, cancelled_at, created_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
RETURNING *;

-- name: GetOrdersByWallet :many
SELECT o1.*, sqlc.embed(o2) FROM orders AS o1
LEFT JOIN orders AS o2 ON o1.id = o2.parent_id AND o2.parent_id IS NOT NULL
WHERE o1.wallet = $1
ORDER BY o1.created_at DESC
LIMIT $2 OFFSET $3;

-- name: GetOrdersByStatus :many
SELECT * FROM orders
WHERE status = ANY(@status::varchar[]);
