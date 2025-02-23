-- name: CreatePrice :one
INSERT INTO prices (id, pool_id, price_usd)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetPrices :many
SELECT * FROM prices
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: GetPriceByPoolID :one
SELECT * FROM prices
WHERE pool_id = $1
ORDER BY created_at DESC
LIMIT 1;

-- name: GetMarketData :many
SELECT 
    time_bucket($1, created_at) AS bucket_time,
    FIRST(price_usd, created_at) AS open_price,
    MAX(price_usd) AS high_price,
    MIN(price_usd) AS low_price,
    LAST(price_usd, created_at) AS close_price,
    AVG(price_usd) AS avg_price,
    COUNT(*) AS number_of_trades
FROM prices
WHERE pool_id = $2
GROUP BY bucket_time
ORDER BY bucket_time DESC;