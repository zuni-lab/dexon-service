// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0

package db

import (
	"context"
)

type Querier interface {
	CreatePool(ctx context.Context, arg CreatePoolParams) (Pool, error)
	CreatePrice(ctx context.Context, arg CreatePriceParams) (Price, error)
	CreateToken(ctx context.Context, arg CreateTokenParams) (Token, error)
	GetMarketData(ctx context.Context, arg GetMarketDataParams) ([]GetMarketDataRow, error)
	GetOrdersByStatus(ctx context.Context, status []string) ([]Order, error)
	GetOrdersByWallet(ctx context.Context, arg GetOrdersByWalletParams) ([]GetOrdersByWalletRow, error)
	GetPool(ctx context.Context, id string) (Pool, error)
	GetPoolByToken(ctx context.Context, arg GetPoolByTokenParams) (Pool, error)
	GetPools(ctx context.Context) ([]Pool, error)
	GetPriceByPoolID(ctx context.Context, poolID string) (Price, error)
	GetPrices(ctx context.Context, arg GetPricesParams) ([]Price, error)
	InsertOrder(ctx context.Context, arg InsertOrderParams) (Order, error)
	PoolDetails(ctx context.Context, id string) (PoolDetailsRow, error)
}

var _ Querier = (*Queries)(nil)
