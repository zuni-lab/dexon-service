-- name: InsertOrder :one
INSERT INTO orders (wallet, from_token, to_token, side, condition, price, created_at)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: GetOrdersByWallet :many
SELECT * FROM orders
WHERE wallet = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;
