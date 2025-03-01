-- name: InsertOrder :one
INSERT INTO orders (
    pool_ids, parent_id, wallet, status, side, type,
    price, amount, slippage, twap_interval_seconds,
    twap_executed_times, twap_current_executed_times,
    twap_min_price, twap_max_price, deadline,
    signature, paths,
    partial_filled_at, filled_at, rejected_at,
    cancelled_at, created_at)
VALUES ($1, $2, $3, $4, $5, $6,
        $7, $8, $9, $10,
        $11, $12, $13,
        $14, $15, $16,
        $17, $18, $19, $20,
        $21, $22)
RETURNING *;

-- name: GetOrdersByWallet :many
SELECT * FROM orders
WHERE wallet = $1
ORDER BY created_at DESC
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
        OR (side = 'SELL' AND type = 'STOP' AND price <= $1)
        OR (side = 'BUY' AND type = 'TWAP' AND price <= $1)
        OR (side = 'SELL' AND type = 'TWAP' AND price >= $1)
    )
    AND status IN ('PENDING', 'PARTIAL_FILLED')
ORDER BY created_at ASC
LIMIT 1;

-- name: UpdateOrder :one
UPDATE orders
SET
    status = COALESCE($2, status),
    twap_current_executed_times = COALESCE($3, twap_current_executed_times),
    filled_at = COALESCE($4, filled_at),
    cancelled_at = COALESCE($5, cancelled_at),
    partial_filled_at = COALESCE($6, partial_filled_at),
    rejected_at = COALESCE($7, rejected_at)
WHERE id = $1
RETURNING *;
