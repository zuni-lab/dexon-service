package services

import (
	"context"
	"database/sql"
	"errors"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"
	"github.com/zuni-lab/dexon-service/config"
	"github.com/zuni-lab/dexon-service/pkg/db"
	"github.com/zuni-lab/dexon-service/pkg/evm"
	"github.com/zuni-lab/dexon-service/pkg/utils"
)

func MatchOrder(ctx context.Context, price *big.Float) (*db.Order, error) {
	numericPrice, err := utils.BigFloatToNumeric(price)
	if err != nil {
		return nil, err
	}

	order, err := db.DB.GetMatchedOrder(ctx, numericPrice)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) || errors.Is(err, pgx.ErrNoRows) {
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
	contract, err := evmManager().DexonInstance(ctx)
	if err != nil {
		return nil, err
	}

	txManager, err := evm.NewTxManager(evmManager().Client())
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	auth, err := bind.NewKeyedTransactorWithChainID(config.Env.RawPrivKey, txManager.ChainID())
	if err != nil {
		return nil, err
	}

	mappedOrder, err := mapOrderToEvmOrder(order)
	if err != nil {
		return nil, err
	}

	data, err := evm.ExecuteOrderData(&contract.DexonTransactor, mappedOrder)
	if err != nil {
		return nil, err
	}

	receipt, err := txManager.SendAndWaitForTx(
		ctx,
		auth,
		common.HexToAddress("contract_address"),
		data,
	)

	if err != nil {
		params := db.RejectOrderParams{
			ID: order.ID,
		}

		_ = params.RejectedAt.Scan(time.Now().UTC())
		_ = params.TxHash.Scan(receipt.TxHash.String())

		rejectedOrder, err := db.DB.RejectOrder(ctx, params)
		if err != nil {
			return nil, err
		}

		return &rejectedOrder, nil
	}

	params := db.FillOrderParams{
		ID: order.ID,
	}
	_ = params.FilledAt.Scan(time.Now().UTC())
	_ = params.TxHash.Scan(receipt.TxHash.String())

	filledOrder, err := db.DB.FillOrder(ctx, params)
	if err != nil {
		return nil, err
	}

	return &filledOrder, nil
}

func mapOrderToEvmOrder(order *db.Order) (*evm.Order, error) {
	userAddress, err := evm.NormalizeAddress(order.Wallet.String)
	if err != nil {
		return nil, err
	}

	nonce := new(big.Int).SetUint64(uint64(order.Nonce))

	path, err := evm.NormalizeHex(order.Paths)
	if err != nil {
		return nil, err
	}

	amount := order.Amount.Int

	triggerPrice := new(big.Int).Mul(order.Price.Int, new(big.Int).Exp(new(big.Int).SetInt64(10), new(big.Int).SetInt64(18), nil))

	slippage := new(big.Int).SetInt64(int64(order.Slippage.Float64 * 10e5))

	deadline := new(big.Int).SetInt64(order.Deadline.Time.Unix())

	signature, err := evm.NormalizeHex(order.Signature.String)
	if err != nil {
		return nil, err
	}

	orderType, err := convertOrderTypeToEvmType(order.Type)
	if err != nil {
		return nil, err
	}
	orderSide, err := convertOrderSideToEvmType(order.Side)
	if err != nil {
		return nil, err
	}

	return &evm.Order{
		Account:      userAddress,
		Nonce:        nonce,
		Path:         path,
		Amount:       amount,
		TriggerPrice: triggerPrice,
		Slippage:     slippage,
		OrderType:    orderType,
		OrderSide:    orderSide,
		Deadline:     deadline,
		Signature:    signature,
	}, nil
}

func convertOrderTypeToEvmType(orderType db.OrderType) (uint8, error) {
	switch orderType {
	case db.OrderTypeLIMIT:
		return 0, nil
	case db.OrderTypeSTOP:
		return 1, nil
	default:
		return 0, errors.New("invalid order type")
	}
}

func convertOrderSideToEvmType(side db.OrderSide) (uint8, error) {
	switch side {
	case db.OrderSideBUY:
		return 0, nil
	case db.OrderSideSELL:
		return 1, nil
	default:
		return 0, errors.New("invalid order side")
	}
}

func fillTwapOrder(ctx context.Context, order *db.Order, price *big.Float) (*db.Order, error) {
	params := db.FillTwapOrderParams{
		ID:                       order.ID,
		TwapCurrentExecutedTimes: order.TwapExecutedTimes,
	}
	_ = params.FilledAt.Scan(time.Now().UTC())

	var err error
	if order.TwapCurrentExecutedTimes.Int32+1 == order.TwapExecutedTimes.Int32 {
		_ = params.FilledAt.Scan(time.Now().UTC())
		params.Status = db.OrderStatusFILLED
	} else {
		params.Status = db.OrderStatusPARTIALFILLED
	}

	amount := calculateTwapAmount(order)
	err = fillPartialOrder(ctx, order, price, amount)
	if err != nil {
		return nil, err
	}

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
	_ = params.FilledAt.Scan(time.Now().UTC())

	_, err := db.DB.InsertOrder(ctx, params)
	if err != nil {
		return err
	}

	return nil
}
