-- name: GetPools :many
SELECT * FROM pools;

-- name: GetPool :one
SELECT * FROM pools
WHERE id = $1 LIMIT 1;

-- name: GetPoolByToken :one
SELECT * FROM pools
WHERE token0_id = $1 AND token1_id = $2 LIMIT 1;

-- name: CreatePool :one
INSERT INTO pools (id, token0_id, token1_id)
VALUES ($1, $2, $3)
RETURNING *;

-- name: PoolDetails :one
SELECT 
    pools.id,
    pools.token0_id,
    pools.token1_id,
    token0.name AS token0_name,
    token0.symbol AS token0_symbol,
    token0.decimals AS token0_decimals,
    token0.is_stable AS token0_is_stable,
    token1.name AS token1_name,
    token1.symbol AS token1_symbol,
    token1.decimals AS token1_decimals,
    token1.is_stable AS token1_is_stable
FROM pools
JOIN tokens AS token0 ON pools.token0_id = token0.id
JOIN tokens AS token1 ON pools.token1_id = token1.id
WHERE pools.id = $1;