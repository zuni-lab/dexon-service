-- name: InsertOrder :one
INSERT INTO orders (parent_id, wallet, pool_id, side, status, type, price, amount, twap_amount, twap_parts, filled_at, partial_filled_at, cancelled_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
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

-- name: GetMatchedOrder :one
SELECT * FROM orders
WHERE (
        (side = 'BUY' AND type = 'LIMIT' AND price <= $1)
        OR (side = 'SELL' AND type = 'LIMIT' AND price >= $1)
        OR (side = 'BUY' AND type = 'STOP' AND price >= $1)
        OR (side = 'SELL' AND type 'STOP' AND price <= $1)
        OR (side = 'BUY' AND type = 'TWAP' AND price <= $1)
        OR (side = 'SELL' AND type = 'TWAP' AND price >= $1)
    )
    AND status IN ('PENDING', 'PARTIAL_FILLED')
    AND type <> 'TWAP'
ORDER BY created_at ASC
LIMIT 1;

-- name: UpdateOrder :one
UPDATE orders
SET
    status = COALESCE($2, status),
    twap_amount = COALESCE($3, twap_amount),
    filled_at = COALESCE($4, filled_at),
    cancelled_at = COALESCE($5, cancelled_at),
    twap_amount = COALESCE($6, twap_amount),
    partial_filled_at = COALESCE($7, partial_filled_at)
WHERE id = $1
RETURNING *;
