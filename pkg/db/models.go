// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0

package db

import (
	"database/sql/driver"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
)

type OrderSide string

const (
	OrderSideBUY  OrderSide = "BUY"
	OrderSideSELL OrderSide = "SELL"
)

func (e *OrderSide) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = OrderSide(s)
	case string:
		*e = OrderSide(s)
	default:
		return fmt.Errorf("unsupported scan type for OrderSide: %T", src)
	}
	return nil
}

type NullOrderSide struct {
	OrderSide OrderSide `json:"orderSide"`
	Valid     bool      `json:"valid"` // Valid is true if OrderSide is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullOrderSide) Scan(value interface{}) error {
	if value == nil {
		ns.OrderSide, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.OrderSide.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullOrderSide) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.OrderSide), nil
}

type OrderStatus string

const (
	OrderStatusPENDING       OrderStatus = "PENDING"
	OrderStatusPARTIALFILLED OrderStatus = "PARTIAL_FILLED"
	OrderStatusFILLED        OrderStatus = "FILLED"
	OrderStatusREJECTED      OrderStatus = "REJECTED"
	OrderStatusCANCELLED     OrderStatus = "CANCELLED"
)

func (e *OrderStatus) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = OrderStatus(s)
	case string:
		*e = OrderStatus(s)
	default:
		return fmt.Errorf("unsupported scan type for OrderStatus: %T", src)
	}
	return nil
}

type NullOrderStatus struct {
	OrderStatus OrderStatus `json:"orderStatus"`
	Valid       bool        `json:"valid"` // Valid is true if OrderStatus is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullOrderStatus) Scan(value interface{}) error {
	if value == nil {
		ns.OrderStatus, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.OrderStatus.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullOrderStatus) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.OrderStatus), nil
}

type OrderType string

const (
	OrderTypeMARKET OrderType = "MARKET"
	OrderTypeLIMIT  OrderType = "LIMIT"
	OrderTypeSTOP   OrderType = "STOP"
	OrderTypeTWAP   OrderType = "TWAP"
)

func (e *OrderType) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = OrderType(s)
	case string:
		*e = OrderType(s)
	default:
		return fmt.Errorf("unsupported scan type for OrderType: %T", src)
	}
	return nil
}

type NullOrderType struct {
	OrderType OrderType `json:"orderType"`
	Valid     bool      `json:"valid"` // Valid is true if OrderType is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullOrderType) Scan(value interface{}) error {
	if value == nil {
		ns.OrderType, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.OrderType.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullOrderType) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.OrderType), nil
}

type BlockProcessingState struct {
	PoolAddress        string           `json:"poolAddress"`
	LastProcessedBlock int64            `json:"lastProcessedBlock"`
	IsBackfill         bool             `json:"isBackfill"`
	CreatedAt          pgtype.Timestamp `json:"createdAt"`
	UpdatedAt          pgtype.Timestamp `json:"updatedAt"`
}

type ChatThread struct {
	ID          int64            `json:"id"`
	ThreadID    string           `json:"threadId"`
	UserAddress string           `json:"userAddress"`
	ThreadName  string           `json:"threadName"`
	CreatedAt   pgtype.Timestamp `json:"createdAt"`
	UpdatedAt   pgtype.Timestamp `json:"updatedAt"`
	IsDeleted   bool             `json:"isDeleted"`
}

type Order struct {
	ID                       int64            `json:"id"`
	PoolIds                  []string         `json:"poolIds"`
	Paths                    string           `json:"paths"`
	Wallet                   pgtype.Text      `json:"wallet"`
	Status                   OrderStatus      `json:"status"`
	Side                     OrderSide        `json:"side"`
	Type                     OrderType        `json:"type"`
	Price                    pgtype.Numeric   `json:"price"`
	Amount                   pgtype.Numeric   `json:"amount"`
	Slippage                 pgtype.Float8    `json:"slippage"`
	Signature                pgtype.Text      `json:"signature"`
	ParentID                 pgtype.Int8      `json:"parentId"`
	TwapIntervalSeconds      pgtype.Int4      `json:"twapIntervalSeconds"`
	TwapExecutedTimes        pgtype.Int4      `json:"twapExecutedTimes"`
	TwapCurrentExecutedTimes pgtype.Int4      `json:"twapCurrentExecutedTimes"`
	TwapMinPrice             pgtype.Numeric   `json:"twapMinPrice"`
	TwapMaxPrice             pgtype.Numeric   `json:"twapMaxPrice"`
	Deadline                 pgtype.Timestamp `json:"deadline"`
	PartialFilledAt          pgtype.Timestamp `json:"partialFilledAt"`
	FilledAt                 pgtype.Timestamp `json:"filledAt"`
	RejectedAt               pgtype.Timestamp `json:"rejectedAt"`
	CancelledAt              pgtype.Timestamp `json:"cancelledAt"`
	CreatedAt                pgtype.Timestamp `json:"createdAt"`
}

type OrderNonce struct {
	Wallet    string           `json:"wallet"`
	Nonce     int64            `json:"nonce"`
	UpdatedAt pgtype.Timestamp `json:"updatedAt"`
	CreatedAt pgtype.Timestamp `json:"createdAt"`
}

type Pool struct {
	ID        string             `json:"id"`
	Token0ID  string             `json:"token0Id"`
	Token1ID  string             `json:"token1Id"`
	CreatedAt pgtype.Timestamptz `json:"createdAt"`
}

type Price struct {
	ID             int64            `json:"id"`
	PoolID         string           `json:"poolId"`
	BlockNumber    int64            `json:"blockNumber"`
	BlockTimestamp int64            `json:"blockTimestamp"`
	Sender         string           `json:"sender"`
	Recipient      string           `json:"recipient"`
	Amount0        int64            `json:"amount0"`
	Amount1        int64            `json:"amount1"`
	SqrtPriceX96   int64            `json:"sqrtPriceX96"`
	Liquidity      int64            `json:"liquidity"`
	Tick           int32            `json:"tick"`
	PriceUsd       pgtype.Numeric   `json:"priceUsd"`
	TimestampUtc   pgtype.Timestamp `json:"timestampUtc"`
	CreatedAt      pgtype.Timestamp `json:"createdAt"`
}

type Token struct {
	ID        string             `json:"id"`
	Name      string             `json:"name"`
	Symbol    string             `json:"symbol"`
	Decimals  int32              `json:"decimals"`
	IsStable  bool               `json:"isStable"`
	CreatedAt pgtype.Timestamptz `json:"createdAt"`
}
