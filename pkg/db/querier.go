// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

type Querier interface {
	CancelAllOrders(ctx context.Context, arg CancelAllOrdersParams) error
	CancelOrder(ctx context.Context, arg CancelOrderParams) (CancelOrderRow, error)
	CountChatThreads(ctx context.Context, userAddress string) (int64, error)
	CountOrdersByWallet(ctx context.Context, arg CountOrdersByWalletParams) (int64, error)
	CreatePool(ctx context.Context, arg CreatePoolParams) (Pool, error)
	CreatePrice(ctx context.Context, arg CreatePriceParams) (Price, error)
	CreateToken(ctx context.Context, arg CreateTokenParams) (Token, error)
	FillOrder(ctx context.Context, arg FillOrderParams) (Order, error)
	FillTwapOrder(ctx context.Context, arg FillTwapOrderParams) (Order, error)
	GetBlockProcessingState(ctx context.Context, arg GetBlockProcessingStateParams) (BlockProcessingState, error)
	GetChatThread(ctx context.Context, arg GetChatThreadParams) (ChatThread, error)
	GetChatThreads(ctx context.Context, arg GetChatThreadsParams) ([]ChatThread, error)
	GetMarketData(ctx context.Context, arg GetMarketDataParams) ([]GetMarketDataRow, error)
	GetMatchedOrder(ctx context.Context, price pgtype.Numeric) (Order, error)
	GetOrderByID(ctx context.Context, arg GetOrderByIDParams) (GetOrderByIDRow, error)
	GetOrdersByWallet(ctx context.Context, arg GetOrdersByWalletParams) ([]GetOrdersByWalletRow, error)
	GetPool(ctx context.Context, id string) (Pool, error)
	GetPoolByToken(ctx context.Context, arg GetPoolByTokenParams) (Pool, error)
	GetPools(ctx context.Context) ([]Pool, error)
	GetPoolsByIDs(ctx context.Context, ids []string) ([]Pool, error)
	GetPriceByPoolID(ctx context.Context, poolID string) (Price, error)
	GetPrices(ctx context.Context, arg GetPricesParams) ([]Price, error)
	InsertOrder(ctx context.Context, arg InsertOrderParams) (InsertOrderRow, error)
	PoolDetails(ctx context.Context, id string) (PoolDetailsRow, error)
	RejectOrder(ctx context.Context, arg RejectOrderParams) (Order, error)
	UpsertBlockProcessingState(ctx context.Context, arg UpsertBlockProcessingStateParams) error
	UpsertChatThread(ctx context.Context, arg UpsertChatThreadParams) (ChatThread, error)
}

var _ Querier = (*Queries)(nil)
