package services

import (
	"context"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jinzhu/copier"
	"github.com/zuni-lab/dexon-service/pkg/custom"
	"github.com/zuni-lab/dexon-service/pkg/db"
	"time"
)

type ListOrdersByWalletQuery struct {
	Wallet string `json:"wallet" validate:"eth_addr"`
	Limit  int32  `json:"limit"`
	Offset int32  `json:"offset"`
}

func ListOrderByWallet(ctx context.Context, query ListOrdersByWalletQuery) ([]db.Order, error) {
	var params db.GetOrdersByWalletParams
	if err := copier.Copy(&params, &query); err != nil {
		return nil, err
	}

	return db.DB.GetOrdersByWallet(ctx, params)
}

type CreateOrderBody struct {
	Wallet    string            `json:"wallet" validate:"eth_addr"`
	FromToken string            `json:"from_token" validate:"eth_addr"`
	ToToken   string            `json:"to_token" validate:"eth_addr"`
	Side      db.OrderSide      `json:"side" validate:"oneof=BUY SELL"`
	Condition db.OrderCondition `json:"condition" validate:"oneof=LIMIT STOP"`
	Price     string            `json:"price" validate:"numeric,gt=0"`
}

func CreateOrder(ctx context.Context, body CreateOrderBody) (*db.Order, error) {
	params := db.InsertOrderParams{
		CreatedAt: pgtype.Timestamptz{Time: time.Now(), Valid: true},
	}
	if err := copier.Copy(&params, &body); err != nil {
		return nil, err
	}

	order, err := db.DB.InsertOrder(ctx, params)
	if err != nil {
		return nil, err
	}

	orderBook := custom.GetOrderBook()
	orderBook.Type(order.Side, order.Condition).ReplaceOrInsert(order)
	return &order, nil
}
