package services

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jinzhu/copier"
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
	Wallet        string       `json:"wallet" validate:"eth_addr"`
	Token0        string       `json:"token0" validate:"eth_addr"`
	Token1        string       `json:"token1" validate:"eth_addr"`
	Side          db.OrderSide `json:"side" validate:"oneof=BUY SELL"`
	Type          db.OrderType `json:"type" validate:"oneof=MARKET LIMIT STOP TWAP"`
	Price         string       `json:"price" validate:"numeric,gt=0"`
	Amount        string       `json:"amount" validate:"numeric,gt=0"`
	TwapTotalTime *int32       `json:"twap_total_time" validate:"omitempty,gt=0"`
}

func CreateOrder(ctx context.Context, body CreateOrderBody) (*db.Order, error) {
	var (
		pool   db.Pool
		params db.InsertOrderParams
	)

	if err := copier.Copy(&params, &body); err != nil {
		return nil, err
	}

	pool, err := db.DB.GetPoolByToken(ctx, db.GetPoolByTokenParams{
		Token0ID: body.Token0,
		Token1ID: body.Token1,
	})
	if err != nil {
		return nil, err
	}

	params.PoolID = pool.ID
	if params.Type == db.OrderTypeMARKET {
		params.Status = db.OrderStatusFILLED
		_ = params.FilledAt.Scan(time.Now())
	}

	order, err := db.DB.InsertOrder(ctx, params)
	if err != nil {
		return nil, err
	}

	return &order, nil
}

func FillPartialOrder(ctx context.Context, parent db.Order, price, amount string) (*db.Order, error) {
	var params db.InsertOrderParams
	if err := copier.Copy(&params, &parent); err != nil {
		return nil, err
	}

	_ = params.ParentID.Scan(parent.ID)
	_ = params.FilledAt.Scan(time.Now())
	_ = params.Price.Scan(price)
	_ = params.Amount.Scan(amount)
	params.Status = db.OrderStatusFILLED

	order, err := db.DB.InsertOrder(ctx, params)
	if err != nil {
		return nil, err
	}

	return &order, nil
}

func MatchOrder(ctx context.Context, price string) (*db.Order, error) {
	var numericPrice pgtype.Numeric
	err := numericPrice.Scan(price)
	if err != nil {
		return nil, err
	}

	order, err := db.DB.GetMatchedOrder(ctx, numericPrice)
	if err != nil {
		return nil, err
	}

	// TODO: Call to contract

	params := db.UpdateOrderParams{
		ID: order.ID,
	}
	if order.Type == db.OrderTypeLIMIT ||
		order.Type == db.OrderTypeSTOP ||
		order.Type == db.OrderTypeTWAP && order.TwapAmount.Int.Cmp(order.Amount.Int) == 0 {
		_ = params.FilledAt.Scan(time.Now())
		params.Status = db.OrderStatusFILLED
	} else if order.Type == db.OrderTypeTWAP {
		params.Status = db.OrderStatusPARTIALFILLED
	} else {
		return nil, errors.New("invalid order")
	}

	order, err = db.DB.UpdateOrder(ctx, params)
	if err != nil {
		return nil, err
	}

	return &order, nil
}
