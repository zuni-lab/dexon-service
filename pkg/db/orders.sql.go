// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: orders.sql

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const getOrdersByWallet = `-- name: GetOrdersByWallet :many
SELECT o1.id, o1.parent_id, o1.wallet, o1.from_token, o1.to_token, o1.status, o1.side, o1.condition, o1.price, o1.amount, o1.twap_total_time, o1.filled_at, o1.cancelled_at, o1.created_at, o2.id, o2.parent_id, o2.wallet, o2.from_token, o2.to_token, o2.status, o2.side, o2.condition, o2.price, o2.amount, o2.twap_total_time, o2.filled_at, o2.cancelled_at, o2.created_at FROM orders AS o1
LEFT JOIN orders AS o2 ON o1.id = o2.parent_id AND o2.parent_id IS NOT NULL
WHERE o1.wallet = $1
ORDER BY o1.created_at DESC
LIMIT $2 OFFSET $3
`

type GetOrdersByWalletParams struct {
	Wallet pgtype.Text `json:"wallet"`
	Limit  int32       `json:"limit"`
	Offset int32       `json:"offset"`
}

type GetOrdersByWalletRow struct {
	ID            int64              `json:"id"`
	ParentID      pgtype.Int8        `json:"parent_id"`
	Wallet        pgtype.Text        `json:"wallet"`
	FromToken     string             `json:"from_token"`
	ToToken       string             `json:"to_token"`
	Status        OrderStatus        `json:"status"`
	Side          OrderSide          `json:"side"`
	Condition     OrderCondition     `json:"condition"`
	Price         pgtype.Numeric     `json:"price"`
	Amount        pgtype.Numeric     `json:"amount"`
	TwapTotalTime pgtype.Int4        `json:"twap_total_time"`
	FilledAt      pgtype.Timestamptz `json:"filled_at"`
	CancelledAt   pgtype.Timestamptz `json:"cancelled_at"`
	CreatedAt     pgtype.Timestamptz `json:"created_at"`
	Order         Order              `json:"order"`
}

func (q *Queries) GetOrdersByWallet(ctx context.Context, arg GetOrdersByWalletParams) ([]GetOrdersByWalletRow, error) {
	rows, err := q.db.Query(ctx, getOrdersByWallet, arg.Wallet, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetOrdersByWalletRow{}
	for rows.Next() {
		var i GetOrdersByWalletRow
		if err := rows.Scan(
			&i.ID,
			&i.ParentID,
			&i.Wallet,
			&i.FromToken,
			&i.ToToken,
			&i.Status,
			&i.Side,
			&i.Condition,
			&i.Price,
			&i.Amount,
			&i.TwapTotalTime,
			&i.FilledAt,
			&i.CancelledAt,
			&i.CreatedAt,
			&i.Order.ID,
			&i.Order.ParentID,
			&i.Order.Wallet,
			&i.Order.FromToken,
			&i.Order.ToToken,
			&i.Order.Status,
			&i.Order.Side,
			&i.Order.Condition,
			&i.Order.Price,
			&i.Order.Amount,
			&i.Order.TwapTotalTime,
			&i.Order.FilledAt,
			&i.Order.CancelledAt,
			&i.Order.CreatedAt,
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
INSERT INTO orders (parent_id, wallet, from_token, to_token, side, condition, price, amount, twap_total_time, created_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
RETURNING id, parent_id, wallet, from_token, to_token, status, side, condition, price, amount, twap_total_time, filled_at, cancelled_at, created_at
`

type InsertOrderParams struct {
	ParentID      pgtype.Int8        `json:"parent_id"`
	Wallet        pgtype.Text        `json:"wallet"`
	FromToken     string             `json:"from_token"`
	ToToken       string             `json:"to_token"`
	Side          OrderSide          `json:"side"`
	Condition     OrderCondition     `json:"condition"`
	Price         pgtype.Numeric     `json:"price"`
	Amount        pgtype.Numeric     `json:"amount"`
	TwapTotalTime pgtype.Int4        `json:"twap_total_time"`
	CreatedAt     pgtype.Timestamptz `json:"created_at"`
}

func (q *Queries) InsertOrder(ctx context.Context, arg InsertOrderParams) (Order, error) {
	row := q.db.QueryRow(ctx, insertOrder,
		arg.ParentID,
		arg.Wallet,
		arg.FromToken,
		arg.ToToken,
		arg.Side,
		arg.Condition,
		arg.Price,
		arg.Amount,
		arg.TwapTotalTime,
		arg.CreatedAt,
	)
	var i Order
	err := row.Scan(
		&i.ID,
		&i.ParentID,
		&i.Wallet,
		&i.FromToken,
		&i.ToToken,
		&i.Status,
		&i.Side,
		&i.Condition,
		&i.Price,
		&i.Amount,
		&i.TwapTotalTime,
		&i.FilledAt,
		&i.CancelledAt,
		&i.CreatedAt,
	)
	return i, err
}
