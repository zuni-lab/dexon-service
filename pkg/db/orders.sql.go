// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: orders.sql

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const cancelAllOrders = `-- name: CancelAllOrders :exec
UPDATE orders
SET
    status = 'CANCELLED',
    cancelled_at = $1
WHERE wallet = $2 AND status NOT IN ('REJECTED', 'FILLED')
    RETURNING
    id, pool_ids, parent_id, wallet, status, side, type,
    price, amount, slippage, twap_interval_seconds,
    twap_executed_times, twap_current_executed_times,
    twap_min_price, twap_max_price, deadline, nonce,
    paths, tx_hash, partial_filled_at, filled_at, rejected_at,
    cancelled_at, created_at
`

type CancelAllOrdersParams struct {
	CancelledAt pgtype.Timestamp `json:"cancelledAt"`
	Wallet      string           `json:"wallet"`
}

func (q *Queries) CancelAllOrders(ctx context.Context, arg CancelAllOrdersParams) error {
	_, err := q.db.Exec(ctx, cancelAllOrders, arg.CancelledAt, arg.Wallet)
	return err
}

const cancelOrder = `-- name: CancelOrder :one
UPDATE orders
SET
    status = 'CANCELLED',
    cancelled_at = $1
WHERE id = $2 AND wallet = $3 AND status NOT IN ('REJECTED', 'FILLED')
RETURNING
    id, pool_ids, parent_id, wallet, status, side, type,
    price, amount, slippage, twap_interval_seconds,
    twap_executed_times, twap_current_executed_times,
    twap_min_price, twap_max_price, deadline, nonce,
    paths, tx_hash, partial_filled_at, filled_at, rejected_at,
    cancelled_at, created_at
`

type CancelOrderParams struct {
	CancelledAt pgtype.Timestamp `json:"cancelledAt"`
	ID          int64            `json:"id"`
	Wallet      string           `json:"wallet"`
}

type CancelOrderRow struct {
	ID                       int64            `json:"id"`
	PoolIds                  []string         `json:"poolIds"`
	ParentID                 pgtype.Int8      `json:"parentId"`
	Wallet                   string           `json:"wallet"`
	Status                   OrderStatus      `json:"status"`
	Side                     OrderSide        `json:"side"`
	Type                     OrderType        `json:"type"`
	Price                    pgtype.Numeric   `json:"price"`
	Amount                   pgtype.Numeric   `json:"amount"`
	Slippage                 pgtype.Float8    `json:"slippage"`
	TwapIntervalSeconds      pgtype.Int4      `json:"twapIntervalSeconds"`
	TwapExecutedTimes        pgtype.Int4      `json:"twapExecutedTimes"`
	TwapCurrentExecutedTimes pgtype.Int4      `json:"twapCurrentExecutedTimes"`
	TwapMinPrice             pgtype.Numeric   `json:"twapMinPrice"`
	TwapMaxPrice             pgtype.Numeric   `json:"twapMaxPrice"`
	Deadline                 pgtype.Timestamp `json:"deadline"`
	Nonce                    int64            `json:"nonce"`
	Paths                    string           `json:"paths"`
	TxHash                   pgtype.Text      `json:"txHash"`
	PartialFilledAt          pgtype.Timestamp `json:"partialFilledAt"`
	FilledAt                 pgtype.Timestamp `json:"filledAt"`
	RejectedAt               pgtype.Timestamp `json:"rejectedAt"`
	CancelledAt              pgtype.Timestamp `json:"cancelledAt"`
	CreatedAt                pgtype.Timestamp `json:"createdAt"`
}

func (q *Queries) CancelOrder(ctx context.Context, arg CancelOrderParams) (CancelOrderRow, error) {
	row := q.db.QueryRow(ctx, cancelOrder, arg.CancelledAt, arg.ID, arg.Wallet)
	var i CancelOrderRow
	err := row.Scan(
		&i.ID,
		&i.PoolIds,
		&i.ParentID,
		&i.Wallet,
		&i.Status,
		&i.Side,
		&i.Type,
		&i.Price,
		&i.Amount,
		&i.Slippage,
		&i.TwapIntervalSeconds,
		&i.TwapExecutedTimes,
		&i.TwapCurrentExecutedTimes,
		&i.TwapMinPrice,
		&i.TwapMaxPrice,
		&i.Deadline,
		&i.Nonce,
		&i.Paths,
		&i.TxHash,
		&i.PartialFilledAt,
		&i.FilledAt,
		&i.RejectedAt,
		&i.CancelledAt,
		&i.CreatedAt,
	)
	return i, err
}

const countOrdersByWallet = `-- name: CountOrdersByWallet :one
SELECT COUNT(*) AS total_counts
FROM orders
WHERE wallet = $1
    AND (
        ARRAY_LENGTH($2::order_status[], 1) IS NULL
        OR (
            status = ANY($2)
            AND (
                status <> 'PENDING'
                OR deadline IS NULL
                OR deadline > NOW() --Skip expired orders
            )
        )
    )
    AND (
        ARRAY_LENGTH($3::order_status[], 1) IS NULL
        OR (
            status <> ANY($3)
            AND (
                status <> 'PENDING'
                OR deadline IS NULL
                OR (status = 'PENDING' AND deadline <= NOW())
            )
        )
    )
    AND (
        ARRAY_LENGTH($4::order_type[], 1) IS NULL
        OR type = ANY($4)
    )
    AND (
        $5::order_side IS NULL
        OR side = $5
    )
`

type CountOrdersByWalletParams struct {
	Wallet    string        `json:"wallet"`
	Status    []OrderStatus `json:"status"`
	NotStatus []OrderStatus `json:"notStatus"`
	Types     []OrderType   `json:"types"`
	Side      NullOrderSide `json:"side"`
}

func (q *Queries) CountOrdersByWallet(ctx context.Context, arg CountOrdersByWalletParams) (int64, error) {
	row := q.db.QueryRow(ctx, countOrdersByWallet,
		arg.Wallet,
		arg.Status,
		arg.NotStatus,
		arg.Types,
		arg.Side,
	)
	var total_counts int64
	err := row.Scan(&total_counts)
	return total_counts, err
}

const fillOrder = `-- name: FillOrder :one
UPDATE orders
SET
    status = 'FILLED',
    filled_at = $1,
    tx_hash = $2,
    actual_amount = $3
WHERE id = $4
RETURNING id, pool_ids, paths, wallet, status, side, type, price, actual_amount, amount, slippage, nonce, signature, tx_hash, parent_id, twap_interval_seconds, twap_executed_times, twap_current_executed_times, twap_min_price, twap_max_price, deadline, partial_filled_at, filled_at, rejected_at, cancelled_at, created_at
`

type FillOrderParams struct {
	FilledAt     pgtype.Timestamp `json:"filledAt"`
	TxHash       pgtype.Text      `json:"txHash"`
	ActualAmount pgtype.Numeric   `json:"actualAmount"`
	ID           int64            `json:"id"`
}

func (q *Queries) FillOrder(ctx context.Context, arg FillOrderParams) (Order, error) {
	row := q.db.QueryRow(ctx, fillOrder,
		arg.FilledAt,
		arg.TxHash,
		arg.ActualAmount,
		arg.ID,
	)
	var i Order
	err := row.Scan(
		&i.ID,
		&i.PoolIds,
		&i.Paths,
		&i.Wallet,
		&i.Status,
		&i.Side,
		&i.Type,
		&i.Price,
		&i.ActualAmount,
		&i.Amount,
		&i.Slippage,
		&i.Nonce,
		&i.Signature,
		&i.TxHash,
		&i.ParentID,
		&i.TwapIntervalSeconds,
		&i.TwapExecutedTimes,
		&i.TwapCurrentExecutedTimes,
		&i.TwapMinPrice,
		&i.TwapMaxPrice,
		&i.Deadline,
		&i.PartialFilledAt,
		&i.FilledAt,
		&i.RejectedAt,
		&i.CancelledAt,
		&i.CreatedAt,
	)
	return i, err
}

const fillTwapOrder = `-- name: FillTwapOrder :one
UPDATE orders
SET
    status = $1,
    twap_current_executed_times = $2,
    partial_filled_at = COALESCE($3, partial_filled_at),
    filled_at = $4,
    tx_hash = $5
WHERE id = $6
RETURNING id, pool_ids, paths, wallet, status, side, type, price, actual_amount, amount, slippage, nonce, signature, tx_hash, parent_id, twap_interval_seconds, twap_executed_times, twap_current_executed_times, twap_min_price, twap_max_price, deadline, partial_filled_at, filled_at, rejected_at, cancelled_at, created_at
`

type FillTwapOrderParams struct {
	Status                   OrderStatus      `json:"status"`
	TwapCurrentExecutedTimes pgtype.Int4      `json:"twapCurrentExecutedTimes"`
	PartialFilledAt          pgtype.Timestamp `json:"partialFilledAt"`
	FilledAt                 pgtype.Timestamp `json:"filledAt"`
	TxHash                   pgtype.Text      `json:"txHash"`
	ID                       int64            `json:"id"`
}

func (q *Queries) FillTwapOrder(ctx context.Context, arg FillTwapOrderParams) (Order, error) {
	row := q.db.QueryRow(ctx, fillTwapOrder,
		arg.Status,
		arg.TwapCurrentExecutedTimes,
		arg.PartialFilledAt,
		arg.FilledAt,
		arg.TxHash,
		arg.ID,
	)
	var i Order
	err := row.Scan(
		&i.ID,
		&i.PoolIds,
		&i.Paths,
		&i.Wallet,
		&i.Status,
		&i.Side,
		&i.Type,
		&i.Price,
		&i.ActualAmount,
		&i.Amount,
		&i.Slippage,
		&i.Nonce,
		&i.Signature,
		&i.TxHash,
		&i.ParentID,
		&i.TwapIntervalSeconds,
		&i.TwapExecutedTimes,
		&i.TwapCurrentExecutedTimes,
		&i.TwapMinPrice,
		&i.TwapMaxPrice,
		&i.Deadline,
		&i.PartialFilledAt,
		&i.FilledAt,
		&i.RejectedAt,
		&i.CancelledAt,
		&i.CreatedAt,
	)
	return i, err
}

const getMatchedOrder = `-- name: GetMatchedOrder :one
SELECT id, pool_ids, paths, wallet, status, side, type, price, actual_amount, amount, slippage, nonce, signature, tx_hash, parent_id, twap_interval_seconds, twap_executed_times, twap_current_executed_times, twap_min_price, twap_max_price, deadline, partial_filled_at, filled_at, rejected_at, cancelled_at, created_at FROM orders
WHERE (
        (side = 'BUY' AND type = 'LIMIT' AND price >= $1)
        OR (side = 'SELL' AND type = 'LIMIT' AND price <= $1)
        OR (side = 'BUY' AND type = 'STOP' AND price <= $1)
        OR (side = 'SELL' AND type = 'STOP' AND price >= $1)
        OR (type = 'TWAP' AND price BETWEEN twap_min_price AND twap_max_price)
    )
    AND status IN ('PENDING', 'PARTIAL_FILLED')
    AND (
        type <> 'TWAP'
        OR ( -- Check TWAP condition
            twap_current_executed_times < twap_executed_times
            AND (
                partial_filled_at IS NULL
                OR partial_filled_at + (twap_interval_seconds || ' seconds')::interval > NOW()
            )
        )
    )
    AND (
        deadline IS NULL
        OR deadline > NOW()
    )
ORDER BY created_at ASC
LIMIT 1
`

func (q *Queries) GetMatchedOrder(ctx context.Context, price pgtype.Numeric) (Order, error) {
	row := q.db.QueryRow(ctx, getMatchedOrder, price)
	var i Order
	err := row.Scan(
		&i.ID,
		&i.PoolIds,
		&i.Paths,
		&i.Wallet,
		&i.Status,
		&i.Side,
		&i.Type,
		&i.Price,
		&i.ActualAmount,
		&i.Amount,
		&i.Slippage,
		&i.Nonce,
		&i.Signature,
		&i.TxHash,
		&i.ParentID,
		&i.TwapIntervalSeconds,
		&i.TwapExecutedTimes,
		&i.TwapCurrentExecutedTimes,
		&i.TwapMinPrice,
		&i.TwapMaxPrice,
		&i.Deadline,
		&i.PartialFilledAt,
		&i.FilledAt,
		&i.RejectedAt,
		&i.CancelledAt,
		&i.CreatedAt,
	)
	return i, err
}

const getMatchedTwapOrder = `-- name: GetMatchedTwapOrder :many
SELECT id, pool_ids, paths, wallet, status, side, type, price, actual_amount, amount, slippage, nonce, signature, tx_hash, parent_id, twap_interval_seconds, twap_executed_times, twap_current_executed_times, twap_min_price, twap_max_price, deadline, partial_filled_at, filled_at, rejected_at, cancelled_at, created_at FROM orders
WHERE type = 'TWAP'
  AND twap_min_price is NULL
  AND status IN ('PENDING', 'PARTIAL_FILLED')
  AND twap_current_executed_times < twap_executed_times
  AND (
        partial_filled_at IS NULL
        OR partial_filled_at + (twap_interval_seconds || ' seconds')::interval > NOW()
  )
`

func (q *Queries) GetMatchedTwapOrder(ctx context.Context) ([]Order, error) {
	rows, err := q.db.Query(ctx, getMatchedTwapOrder)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Order{}
	for rows.Next() {
		var i Order
		if err := rows.Scan(
			&i.ID,
			&i.PoolIds,
			&i.Paths,
			&i.Wallet,
			&i.Status,
			&i.Side,
			&i.Type,
			&i.Price,
			&i.ActualAmount,
			&i.Amount,
			&i.Slippage,
			&i.Nonce,
			&i.Signature,
			&i.TxHash,
			&i.ParentID,
			&i.TwapIntervalSeconds,
			&i.TwapExecutedTimes,
			&i.TwapCurrentExecutedTimes,
			&i.TwapMinPrice,
			&i.TwapMaxPrice,
			&i.Deadline,
			&i.PartialFilledAt,
			&i.FilledAt,
			&i.RejectedAt,
			&i.CancelledAt,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getOrderByID = `-- name: GetOrderByID :one
SELECT id, pool_ids, parent_id, wallet, status, side, type,
       price, amount, slippage, twap_interval_seconds,
       twap_executed_times, twap_current_executed_times,
       twap_min_price, twap_max_price, deadline, nonce,
       paths, tx_hash, partial_filled_at, filled_at, rejected_at,
       cancelled_at, created_at
FROM orders
WHERE wallet = $1 AND id = $2
`

type GetOrderByIDParams struct {
	Wallet string `json:"wallet"`
	ID     int64  `json:"id"`
}

type GetOrderByIDRow struct {
	ID                       int64            `json:"id"`
	PoolIds                  []string         `json:"poolIds"`
	ParentID                 pgtype.Int8      `json:"parentId"`
	Wallet                   string           `json:"wallet"`
	Status                   OrderStatus      `json:"status"`
	Side                     OrderSide        `json:"side"`
	Type                     OrderType        `json:"type"`
	Price                    pgtype.Numeric   `json:"price"`
	Amount                   pgtype.Numeric   `json:"amount"`
	Slippage                 pgtype.Float8    `json:"slippage"`
	TwapIntervalSeconds      pgtype.Int4      `json:"twapIntervalSeconds"`
	TwapExecutedTimes        pgtype.Int4      `json:"twapExecutedTimes"`
	TwapCurrentExecutedTimes pgtype.Int4      `json:"twapCurrentExecutedTimes"`
	TwapMinPrice             pgtype.Numeric   `json:"twapMinPrice"`
	TwapMaxPrice             pgtype.Numeric   `json:"twapMaxPrice"`
	Deadline                 pgtype.Timestamp `json:"deadline"`
	Nonce                    int64            `json:"nonce"`
	Paths                    string           `json:"paths"`
	TxHash                   pgtype.Text      `json:"txHash"`
	PartialFilledAt          pgtype.Timestamp `json:"partialFilledAt"`
	FilledAt                 pgtype.Timestamp `json:"filledAt"`
	RejectedAt               pgtype.Timestamp `json:"rejectedAt"`
	CancelledAt              pgtype.Timestamp `json:"cancelledAt"`
	CreatedAt                pgtype.Timestamp `json:"createdAt"`
}

func (q *Queries) GetOrderByID(ctx context.Context, arg GetOrderByIDParams) (GetOrderByIDRow, error) {
	row := q.db.QueryRow(ctx, getOrderByID, arg.Wallet, arg.ID)
	var i GetOrderByIDRow
	err := row.Scan(
		&i.ID,
		&i.PoolIds,
		&i.ParentID,
		&i.Wallet,
		&i.Status,
		&i.Side,
		&i.Type,
		&i.Price,
		&i.Amount,
		&i.Slippage,
		&i.TwapIntervalSeconds,
		&i.TwapExecutedTimes,
		&i.TwapCurrentExecutedTimes,
		&i.TwapMinPrice,
		&i.TwapMaxPrice,
		&i.Deadline,
		&i.Nonce,
		&i.Paths,
		&i.TxHash,
		&i.PartialFilledAt,
		&i.FilledAt,
		&i.RejectedAt,
		&i.CancelledAt,
		&i.CreatedAt,
	)
	return i, err
}

const getOrdersByWallet = `-- name: GetOrdersByWallet :many
SELECT id, pool_ids, parent_id, wallet, status, side, type,
       price, amount, actual_amount, slippage, twap_interval_seconds,
       twap_executed_times, twap_current_executed_times,
       twap_min_price, twap_max_price, deadline, nonce,
       paths, tx_hash, partial_filled_at, filled_at, rejected_at,
       cancelled_at, created_at
FROM orders
WHERE wallet = $1
    AND (
        ARRAY_LENGTH($4::order_status[], 1) IS NULL
        OR (
            status = ANY($4)
            AND (
                status <> 'PENDING'
                OR deadline IS NULL
                OR deadline > NOW() --Skip expired orders
            )
        )
    )
    AND (
        ARRAY_LENGTH($5::order_status[], 1) IS NULL
        OR (
        status <> ANY($5)
            AND (
                status <> 'PENDING'
                OR deadline IS NULL
                OR (status = 'PENDING' AND deadline <= NOW())
            )
        )
    )
    AND (
        ARRAY_LENGTH($6::order_type[], 1) IS NULL
        OR type = ANY($6)
    )
    AND (
        $7::order_side IS NULL
        OR side = $7
    )
ORDER BY created_at DESC
LIMIT $2 OFFSET $3
`

type GetOrdersByWalletParams struct {
	Wallet    string        `json:"wallet"`
	Limit     int32         `json:"limit"`
	Offset    int32         `json:"offset"`
	Status    []OrderStatus `json:"status"`
	NotStatus []OrderStatus `json:"notStatus"`
	Types     []OrderType   `json:"types"`
	Side      NullOrderSide `json:"side"`
}

type GetOrdersByWalletRow struct {
	ID                       int64            `json:"id"`
	PoolIds                  []string         `json:"poolIds"`
	ParentID                 pgtype.Int8      `json:"parentId"`
	Wallet                   string           `json:"wallet"`
	Status                   OrderStatus      `json:"status"`
	Side                     OrderSide        `json:"side"`
	Type                     OrderType        `json:"type"`
	Price                    pgtype.Numeric   `json:"price"`
	Amount                   pgtype.Numeric   `json:"amount"`
	ActualAmount             pgtype.Numeric   `json:"actualAmount"`
	Slippage                 pgtype.Float8    `json:"slippage"`
	TwapIntervalSeconds      pgtype.Int4      `json:"twapIntervalSeconds"`
	TwapExecutedTimes        pgtype.Int4      `json:"twapExecutedTimes"`
	TwapCurrentExecutedTimes pgtype.Int4      `json:"twapCurrentExecutedTimes"`
	TwapMinPrice             pgtype.Numeric   `json:"twapMinPrice"`
	TwapMaxPrice             pgtype.Numeric   `json:"twapMaxPrice"`
	Deadline                 pgtype.Timestamp `json:"deadline"`
	Nonce                    int64            `json:"nonce"`
	Paths                    string           `json:"paths"`
	TxHash                   pgtype.Text      `json:"txHash"`
	PartialFilledAt          pgtype.Timestamp `json:"partialFilledAt"`
	FilledAt                 pgtype.Timestamp `json:"filledAt"`
	RejectedAt               pgtype.Timestamp `json:"rejectedAt"`
	CancelledAt              pgtype.Timestamp `json:"cancelledAt"`
	CreatedAt                pgtype.Timestamp `json:"createdAt"`
}

func (q *Queries) GetOrdersByWallet(ctx context.Context, arg GetOrdersByWalletParams) ([]GetOrdersByWalletRow, error) {
	rows, err := q.db.Query(ctx, getOrdersByWallet,
		arg.Wallet,
		arg.Limit,
		arg.Offset,
		arg.Status,
		arg.NotStatus,
		arg.Types,
		arg.Side,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetOrdersByWalletRow{}
	for rows.Next() {
		var i GetOrdersByWalletRow
		if err := rows.Scan(
			&i.ID,
			&i.PoolIds,
			&i.ParentID,
			&i.Wallet,
			&i.Status,
			&i.Side,
			&i.Type,
			&i.Price,
			&i.Amount,
			&i.ActualAmount,
			&i.Slippage,
			&i.TwapIntervalSeconds,
			&i.TwapExecutedTimes,
			&i.TwapCurrentExecutedTimes,
			&i.TwapMinPrice,
			&i.TwapMaxPrice,
			&i.Deadline,
			&i.Nonce,
			&i.Paths,
			&i.TxHash,
			&i.PartialFilledAt,
			&i.FilledAt,
			&i.RejectedAt,
			&i.CancelledAt,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const insertOrder = `-- name: InsertOrder :one
INSERT INTO orders (
    pool_ids, parent_id, wallet, status, side, type,
    price, amount, slippage, twap_interval_seconds,
    twap_executed_times, twap_current_executed_times,
    twap_min_price, twap_max_price, deadline,
    signature, paths, nonce, tx_hash,
    partial_filled_at, filled_at, rejected_at,
    cancelled_at, created_at)
VALUES ($1, $2, $3, $4, $5, $6,
        $7, $8, $9, $10,
        $11, $12, $13,
        $14, $15, $16,
        $17, $18, $19, $20,
        $21, $22, $23, $24)
RETURNING
    id, pool_ids, parent_id, wallet, status, side, type,
    price, amount, slippage, twap_interval_seconds,
    twap_executed_times, twap_current_executed_times,
    twap_min_price, twap_max_price, deadline, nonce,
    paths, tx_hash, partial_filled_at, filled_at, rejected_at,
    cancelled_at, created_at
`

type InsertOrderParams struct {
	PoolIds                  []string         `json:"poolIds"`
	ParentID                 pgtype.Int8      `json:"parentId"`
	Wallet                   string           `json:"wallet"`
	Status                   OrderStatus      `json:"status"`
	Side                     OrderSide        `json:"side"`
	Type                     OrderType        `json:"type"`
	Price                    pgtype.Numeric   `json:"price"`
	Amount                   pgtype.Numeric   `json:"amount"`
	Slippage                 pgtype.Float8    `json:"slippage"`
	TwapIntervalSeconds      pgtype.Int4      `json:"twapIntervalSeconds"`
	TwapExecutedTimes        pgtype.Int4      `json:"twapExecutedTimes"`
	TwapCurrentExecutedTimes pgtype.Int4      `json:"twapCurrentExecutedTimes"`
	TwapMinPrice             pgtype.Numeric   `json:"twapMinPrice"`
	TwapMaxPrice             pgtype.Numeric   `json:"twapMaxPrice"`
	Deadline                 pgtype.Timestamp `json:"deadline"`
	Signature                string           `json:"signature"`
	Paths                    string           `json:"paths"`
	Nonce                    int64            `json:"nonce"`
	TxHash                   pgtype.Text      `json:"txHash"`
	PartialFilledAt          pgtype.Timestamp `json:"partialFilledAt"`
	FilledAt                 pgtype.Timestamp `json:"filledAt"`
	RejectedAt               pgtype.Timestamp `json:"rejectedAt"`
	CancelledAt              pgtype.Timestamp `json:"cancelledAt"`
	CreatedAt                pgtype.Timestamp `json:"createdAt"`
}

type InsertOrderRow struct {
	ID                       int64            `json:"id"`
	PoolIds                  []string         `json:"poolIds"`
	ParentID                 pgtype.Int8      `json:"parentId"`
	Wallet                   string           `json:"wallet"`
	Status                   OrderStatus      `json:"status"`
	Side                     OrderSide        `json:"side"`
	Type                     OrderType        `json:"type"`
	Price                    pgtype.Numeric   `json:"price"`
	Amount                   pgtype.Numeric   `json:"amount"`
	Slippage                 pgtype.Float8    `json:"slippage"`
	TwapIntervalSeconds      pgtype.Int4      `json:"twapIntervalSeconds"`
	TwapExecutedTimes        pgtype.Int4      `json:"twapExecutedTimes"`
	TwapCurrentExecutedTimes pgtype.Int4      `json:"twapCurrentExecutedTimes"`
	TwapMinPrice             pgtype.Numeric   `json:"twapMinPrice"`
	TwapMaxPrice             pgtype.Numeric   `json:"twapMaxPrice"`
	Deadline                 pgtype.Timestamp `json:"deadline"`
	Nonce                    int64            `json:"nonce"`
	Paths                    string           `json:"paths"`
	TxHash                   pgtype.Text      `json:"txHash"`
	PartialFilledAt          pgtype.Timestamp `json:"partialFilledAt"`
	FilledAt                 pgtype.Timestamp `json:"filledAt"`
	RejectedAt               pgtype.Timestamp `json:"rejectedAt"`
	CancelledAt              pgtype.Timestamp `json:"cancelledAt"`
	CreatedAt                pgtype.Timestamp `json:"createdAt"`
}

func (q *Queries) InsertOrder(ctx context.Context, arg InsertOrderParams) (InsertOrderRow, error) {
	row := q.db.QueryRow(ctx, insertOrder,
		arg.PoolIds,
		arg.ParentID,
		arg.Wallet,
		arg.Status,
		arg.Side,
		arg.Type,
		arg.Price,
		arg.Amount,
		arg.Slippage,
		arg.TwapIntervalSeconds,
		arg.TwapExecutedTimes,
		arg.TwapCurrentExecutedTimes,
		arg.TwapMinPrice,
		arg.TwapMaxPrice,
		arg.Deadline,
		arg.Signature,
		arg.Paths,
		arg.Nonce,
		arg.TxHash,
		arg.PartialFilledAt,
		arg.FilledAt,
		arg.RejectedAt,
		arg.CancelledAt,
		arg.CreatedAt,
	)
	var i InsertOrderRow
	err := row.Scan(
		&i.ID,
		&i.PoolIds,
		&i.ParentID,
		&i.Wallet,
		&i.Status,
		&i.Side,
		&i.Type,
		&i.Price,
		&i.Amount,
		&i.Slippage,
		&i.TwapIntervalSeconds,
		&i.TwapExecutedTimes,
		&i.TwapCurrentExecutedTimes,
		&i.TwapMinPrice,
		&i.TwapMaxPrice,
		&i.Deadline,
		&i.Nonce,
		&i.Paths,
		&i.TxHash,
		&i.PartialFilledAt,
		&i.FilledAt,
		&i.RejectedAt,
		&i.CancelledAt,
		&i.CreatedAt,
	)
	return i, err
}

const rejectOrder = `-- name: RejectOrder :one
UPDATE orders
SET
    status = 'REJECTED',
    rejected_at = $1
WHERE id = $2
RETURNING id, pool_ids, paths, wallet, status, side, type, price, actual_amount, amount, slippage, nonce, signature, tx_hash, parent_id, twap_interval_seconds, twap_executed_times, twap_current_executed_times, twap_min_price, twap_max_price, deadline, partial_filled_at, filled_at, rejected_at, cancelled_at, created_at
`

type RejectOrderParams struct {
	RejectedAt pgtype.Timestamp `json:"rejectedAt"`
	ID         int64            `json:"id"`
}

func (q *Queries) RejectOrder(ctx context.Context, arg RejectOrderParams) (Order, error) {
	row := q.db.QueryRow(ctx, rejectOrder, arg.RejectedAt, arg.ID)
	var i Order
	err := row.Scan(
		&i.ID,
		&i.PoolIds,
		&i.Paths,
		&i.Wallet,
		&i.Status,
		&i.Side,
		&i.Type,
		&i.Price,
		&i.ActualAmount,
		&i.Amount,
		&i.Slippage,
		&i.Nonce,
		&i.Signature,
		&i.TxHash,
		&i.ParentID,
		&i.TwapIntervalSeconds,
		&i.TwapExecutedTimes,
		&i.TwapCurrentExecutedTimes,
		&i.TwapMinPrice,
		&i.TwapMaxPrice,
		&i.Deadline,
		&i.PartialFilledAt,
		&i.FilledAt,
		&i.RejectedAt,
		&i.CancelledAt,
		&i.CreatedAt,
	)
	return i, err
}
