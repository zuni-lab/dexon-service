package services

import (
	"context"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jinzhu/copier"
	"github.com/zuni-lab/dexon-service/pkg/custom"
	"github.com/zuni-lab/dexon-service/pkg/db"
	"slices"
	"time"
)

type ListOrdersByWalletQuery struct {
	Wallet string `json:"wallet" validate:"eth_addr"`
	Limit  int32  `json:"limit"`
	Offset int32  `json:"offset"`
}

type ListOrdersByWalletResponseItem struct {
	db.Order
	Children []db.Order `json:"children"`
}

func ListOrderByWallet(ctx context.Context, query ListOrdersByWalletQuery) ([]ListOrdersByWalletResponseItem, error) {
	var params db.GetOrdersByWalletParams
	if err := copier.Copy(&params, &query); err != nil {
		return nil, err
	}

	orders, err := db.DB.GetOrdersByWallet(ctx, params)
	if err != nil {
		return nil, err
	}

	var (
		item ListOrdersByWalletResponseItem
		res  []ListOrdersByWalletResponseItem
	)
	for _, order := range orders {
		if idx := slices.IndexFunc(res, func(item ListOrdersByWalletResponseItem) bool {
			return item.ID == order.ID
		}); idx != -1 {
			res[idx].Children = append(res[idx].Children, order.Order)
		}

		err = copier.Copy(&item, &order)
		if err != nil {
			return nil, err
		}
		res = append(res, item)
	}

	return res, nil
}

type CreateOrderBody struct {
	Wallet        string            `json:"wallet" validate:"eth_addr"`
	FromToken     string            `json:"from_token" validate:"eth_addr"`
	ToToken       string            `json:"to_token" validate:"eth_addr"`
	Side          db.OrderSide      `json:"side" validate:"oneof=BUY SELL"`
	Condition     db.OrderCondition `json:"condition" validate:"oneof=LIMIT STOP TWAP"`
	Price         string            `json:"price" validate:"numeric,gt=0"`
	Amount        string            `json:"amount" validate:"numeric,gt=0"`
	TwapTotalTime *int32            `json:"twap_total_time" validate:"omitempty,gt=0"`
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
