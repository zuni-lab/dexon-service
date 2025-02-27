package services

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jinzhu/copier"
	"github.com/zuni-lab/dexon-service/pkg/db"
	"time"
)

type ListOrdersByWalletQuery struct {
	Wallet string `query:"wallet" validate:"eth_addr"`
	Limit  int32  `query:"limit" validate:"gt=0"`
	Offset int32  `query:"offset" validate:"gte=0"`
}

func ListOrderByWallet(ctx context.Context, query ListOrdersByWalletQuery) ([]db.Order, error) {
	var params db.GetOrdersByWalletParams
	if err := copier.Copy(&params, &query); err != nil {
		return nil, err
	}

	orders, err := db.DB.GetOrdersByWallet(ctx, params)
	if err != nil {
		return nil, err
	}

	return orders, nil
}

type CreateOrderBody struct {
	Wallet        string       `json:"wallet" validate:"eth_addr"`
	PoolIDs       []string     `json:"poolIds" validate:"min=1,dive,eth_addr"`
	Side          db.OrderSide `json:"side" validate:"oneof=BUY SELL"`
	Type          db.OrderType `json:"type" validate:"oneof=MARKET LIMIT STOP TWAP"`
	Price         string       `json:"price" validate:"numeric,gt=0"`
	Amount        string       `json:"amount" validate:"numeric,gt=0"`
	TwapTotalTime *int32       `json:"twapTotalTime" validate:"omitempty,gt=0"`
}

func CreateOrder(ctx context.Context, body CreateOrderBody) (*db.Order, error) {
	var params db.InsertOrderParams
	if err := copier.Copy(&params, &body); err != nil {
		return nil, err
	}

	if params.Type == db.OrderTypeMARKET {
		params.Status = db.OrderStatusFILLED
		_ = params.FilledAt.Scan(time.Now())
	} else {
		params.Status = db.OrderStatusPENDING
	}

	pools, err := db.DB.GetPoolsByIDs(ctx, body.PoolIDs)
	if err != nil {
		return nil, err
	} else if len(pools) != len(body.PoolIDs) {
		return nil, errors.New("invalid pool ids")
	}

	order, err := db.DB.InsertOrder(ctx, params)
	if err != nil {
		return nil, err
	}

	return &order, nil
}

func FillPartialOrder(ctx context.Context, parent db.Order, price, amount string) (*db.Order, error) {
	params := db.InsertOrderParams{
		PoolIds:                  parent.PoolIds,
		Wallet:                   parent.Wallet,
		Status:                   db.OrderStatusFILLED,
		Side:                     parent.Side,
		Type:                     db.OrderTypeTWAP,
		Amount:                   parent.Amount,
		Slippage:                 parent.Slippage,
		TwapIntervalSeconds:      parent.TwapIntervalSeconds,
		TwapExecutedTimes:        parent.TwapExecutedTimes,
		TwapCurrentExecutedTimes: parent.TwapCurrentExecutedTimes,
		TwapMinPrice:             parent.TwapMinPrice,
		TwapMaxPrice:             parent.TwapMaxPrice,
	}
	_ = params.ParentID.Scan(parent.ID)
	_ = params.Price.Scan(price)
	_ = params.Amount.Scan(amount)
	_ = params.FilledAt.Scan(time.Now())

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
		order.Type == db.OrderTypeTWAP && order.TwapCurrentExecutedTimes.Int32+1 == order.TwapExecutedTimes.Int32 {
		_ = params.FilledAt.Scan(time.Now())
		params.Status = db.OrderStatusFILLED

		if order.Type == db.OrderTypeTWAP {
			order.TwapCurrentExecutedTimes.Int32 = order.TwapExecutedTimes.Int32
		}
	} else if order.Type == db.OrderTypeTWAP {
		params.Status = db.OrderStatusPARTIALFILLED
		// TODO: create child order
	} else {
		return nil, errors.New("invalid order")
	}

	order, err = db.DB.UpdateOrder(ctx, params)
	if err != nil {
		return nil, err
	}

	return &order, nil
}
