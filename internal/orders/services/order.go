package services

import (
	"context"
	"database/sql"
	"errors"
	"math/big"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jinzhu/copier"
	"github.com/rs/zerolog/log"
	"github.com/zuni-lab/dexon-service/pkg/db"
	"github.com/zuni-lab/dexon-service/pkg/utils"
)

type ListOrdersByWalletQuery struct {
	Wallet string           `query:"wallet" validate:"eth_addr"`
	Status []db.OrderStatus `query:"status" validate:"dive,oneof=PENDING PARTIAL_FILLED FILLED REJECTED CANCELLED"`
	Types  []db.OrderType   `query:"types" validate:"dive,oneof=MARKET LIMIT STOP TWAP"`
	Side   *string          `query:"side" validate:"omitempty,oneof=BUY SELL"`
	Limit  int32            `query:"limit" validate:"gt=0"`
	Offset int32            `query:"offset" validate:"gte=0"`
}

func ListOrderByWallet(ctx context.Context, query ListOrdersByWalletQuery) ([]db.Order, error) {
	query.Wallet = utils.NormalizeAddress(query.Wallet)

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

type GetOrderByIDQuery struct {
	ID     int64
	Wallet string `query:"wallet" validate:"eth_addr"`
}

func GetOrderByID(ctx context.Context, query GetOrderByIDQuery) (*db.Order, error) {
	query.Wallet = utils.NormalizeAddress(query.Wallet)

	var params db.GetOrderByIDParams
	if err := copier.Copy(&params, &query); err != nil {
		return nil, err
	}

	order, err := db.DB.GetOrderByID(ctx, params)
	if err != nil {
		return nil, err
	}

	return &order, nil
}

type CreateOrderBody struct {
	Wallet    string       `json:"wallet" validate:"eth_addr"`
	PoolIDs   []string     `json:"poolIds" validate:"min=1,dive,eth_addr"`
	Side      db.OrderSide `json:"side" validate:"oneof=BUY SELL"`
	Type      db.OrderType `json:"type" validate:"oneof=MARKET LIMIT STOP TWAP"`
	Price     *string      `json:"price" validate:"required_unless=Type TWAP,numeric,gt=0"`
	Amount    string       `json:"amount" validate:"numeric,gt=0"`
	Slippage  float64      `json:"slippage" validate:"gte=0"`
	Signature string       `json:"signature" validate:"max=130"`
	Paths     string       `json:"paths" validate:"max=256"`
	Deadline  *time.Time   `json:"deadline" validate:"omitempty,datetime=2006-01-02 15:04:05"`

	TwapIntervalSeconds *int64  `json:"twapIntervalSeconds" validate:"required_if=Type TWAP,gt=59"`
	TwapExecutedTimes   *int64  `json:"twapExecutedTimes" validate:"required_if=Type TWAP,gt=0"`
	TwapMinPrice        *string `json:"twapMinPrice" validate:"omitempty,numeric,gte=0"`
	TwapMaxPrice        *string `json:"twapMaxPrice" validate:"required_with=TwapMinPrice,numeric,gtefield=TwapMinPrice"`
}

func CreateOrder(ctx context.Context, body CreateOrderBody) (*db.Order, error) {
	body.Wallet = utils.NormalizeAddress(body.Wallet)

	var params db.InsertOrderParams
	if err := copier.Copy(&params, &body); err != nil {
		return nil, err
	}

	now := time.Now()
	_ = params.CreatedAt.Scan(now)
	if params.Type == db.OrderTypeMARKET {
		_ = params.FilledAt.Scan(now)
		params.CreatedAt = params.FilledAt
		params.Status = db.OrderStatusFILLED
	} else {
		params.Status = db.OrderStatusPENDING

		if params.Type == db.OrderTypeTWAP {
			_ = params.Price.Scan("0")
			_ = params.Slippage.Scan("0")
		} else {
			params.TwapIntervalSeconds.Valid = false
			params.TwapExecutedTimes.Valid = false
			params.TwapMinPrice.Valid = false
			params.TwapMaxPrice.Valid = false
		}
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

type CancelOrderBody struct {
	ID     int64
	Wallet string `json:"wallet" validate:"eth_addr"`
}

func CancelOrder(ctx context.Context, body CancelOrderBody) (*db.Order, error) {
	body.Wallet = utils.NormalizeAddress(body.Wallet)

	var params db.CancelOrderParams
	if err := copier.Copy(&params, &body); err != nil {
		return nil, err
	}

	_ = params.CancelledAt.Scan(time.Now())
	order, err := db.DB.CancelOrder(ctx, params)
	if err != nil {
		return nil, err
	}

	return &order, nil
}

func MatchOrder(ctx context.Context, price *big.Float) (*db.Order, error) {
	numericPrice, err := utils.BigFloatToNumeric(price)
	if err != nil {
		return nil, err
	}

	order, err := db.DB.GetMatchedOrder(ctx, numericPrice)
	if err != nil {
		if err == sql.ErrNoRows || err == pgx.ErrNoRows {
			return nil, errors.New("no order matched")
		}
		return nil, err
	}

	log.Info().Any("matched orders", order).Msg("Matched order")

	var filledOrder *db.Order
	if order.Type == db.OrderTypeTWAP {
		filledOrder, err = fillTwapOrder(ctx, &order, price)
	} else {
		filledOrder, err = fillOrder(ctx, &order)
	}

	if err != nil {
		return nil, err
	}

	return filledOrder, nil
}

func fillOrder(ctx context.Context, order *db.Order) (*db.Order, error) {
	params := db.FillOrderParams{
		ID: order.ID,
	}
	_ = params.FilledAt.Scan(time.Now())

	filledOrder, err := db.DB.FillOrder(ctx, params)
	if err != nil {
		return nil, err
	}

	// TODO: Call to contract

	return &filledOrder, nil
}

func fillTwapOrder(ctx context.Context, order *db.Order, price *big.Float) (*db.Order, error) {
	params := db.FillTwapOrderParams{
		ID:                       order.ID,
		TwapCurrentExecutedTimes: order.TwapExecutedTimes,
	}
	_ = params.FilledAt.Scan(time.Now())

	var err error
	if order.TwapCurrentExecutedTimes.Int32+1 == order.TwapExecutedTimes.Int32 {
		_ = params.FilledAt.Scan(time.Now())
		params.Status = db.OrderStatusFILLED
	} else {
		params.Status = db.OrderStatusPARTIALFILLED
	}

	amount := calculateTwapAmount(order)
	err = fillPartialOrder(ctx, order, price, amount)
	if err != nil {
		return nil, err
	}

	// TODO: Call to contract

	filledOrder, err := db.DB.FillTwapOrder(ctx, params)
	if err != nil {
		return nil, err
	}

	return &filledOrder, nil
}

func calculateTwapAmount(order *db.Order) *big.Float {
	divisor := big.NewFloat(float64(order.TwapCurrentExecutedTimes.Int32))
	f64Amount, _ := order.Amount.Float64Value()
	bigAmount := big.NewFloat(f64Amount.Float64)

	return new(big.Float).Quo(bigAmount, divisor)
}

func fillPartialOrder(ctx context.Context, parent *db.Order, price, amount *big.Float) error {
	params := db.InsertOrderParams{
		PoolIds: parent.PoolIds,
		Wallet:  parent.Wallet,
		Status:  db.OrderStatusFILLED,
		Side:    parent.Side,
		Type:    db.OrderTypeTWAP,
		Amount:  parent.Amount,
	}
	_ = params.ParentID.Scan(parent.ID)
	_ = params.Price.Scan(price.String())
	_ = params.Amount.Scan(amount.String())
	_ = params.FilledAt.Scan(time.Now())

	_, err := db.DB.InsertOrder(ctx, params)
	if err != nil {
		return err
	}

	return nil
}
