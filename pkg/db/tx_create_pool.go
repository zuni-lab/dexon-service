package db

import (
	"context"
)

// type CreatePoolTxParams struct {
// 	ID     string
// 	Token0 Token
// 	Token1 Token
// }

type CreatePoolTxResult struct {
	Pool Pool
}

func (store *SqlStore) CreatePoolTx(ctx context.Context, arg PoolDetailsRow) (CreatePoolTxResult, error) {
	var result CreatePoolTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		_, err = q.CreateToken(ctx, CreateTokenParams{
			ID:       arg.Token0ID,
			Name:     arg.Token0Name,
			Symbol:   arg.Token0Symbol,
			Decimals: arg.Token0Decimals,
			IsStable: arg.Token0IsStable,
		})

		if err != nil {
			return err
		}

		_, err = q.CreateToken(ctx, CreateTokenParams{
			ID:       arg.Token1ID,
			Name:     arg.Token1Name,
			Symbol:   arg.Token1Symbol,
			Decimals: arg.Token1Decimals,
			IsStable: arg.Token1IsStable,
		})

		if err != nil {
			return err
		}

		result.Pool, err = q.CreatePool(ctx, CreatePoolParams{
			ID:       arg.ID,
			Token0ID: arg.Token0ID,
			Token1ID: arg.Token1ID,
		})

		if err != nil {
			return err
		}

		return nil
	})

	return result, err
}

type CreateBatchPoolsTxParams struct {
	Pools []PoolDetailsRow
}

type CreateBatchPoolsTxResult struct {
	Pools []Pool
}

func (store *SqlStore) CreateBatchPoolsTx(ctx context.Context, arg CreateBatchPoolsTxParams) (CreateBatchPoolsTxResult, error) {
	var result CreateBatchPoolsTxResult
	result.Pools = make([]Pool, 0, len(arg.Pools))

	err := store.execTx(ctx, func(q *Queries) error {
		tokenMap := make(map[string]struct{})

		for _, poolParam := range arg.Pools {
			if _, exists := tokenMap[poolParam.Token0ID]; !exists {
				_, err := q.CreateToken(ctx, CreateTokenParams{
					ID:       poolParam.Token0ID,
					Name:     poolParam.Token0Name,
					Symbol:   poolParam.Token0Symbol,
					Decimals: poolParam.Token0Decimals,
				})
				if err != nil {
					return err
				}
				tokenMap[poolParam.Token0ID] = struct{}{}
			}

			if _, exists := tokenMap[poolParam.Token1ID]; !exists {
				_, err := q.CreateToken(ctx, CreateTokenParams{
					ID:       poolParam.Token1ID,
					Name:     poolParam.Token1Name,
					Symbol:   poolParam.Token1Symbol,
					Decimals: poolParam.Token1Decimals,
				})
				if err != nil {
					return err
				}
				tokenMap[poolParam.Token1ID] = struct{}{}
			}
		}

		for _, poolParam := range arg.Pools {
			pool, err := q.CreatePool(ctx, CreatePoolParams{
				ID:       poolParam.ID,
				Token0ID: poolParam.Token0ID,
				Token1ID: poolParam.Token1ID,
			})
			if err != nil {
				return err
			}
			result.Pools = append(result.Pools, pool)
		}

		return nil
	})

	return result, err
}
