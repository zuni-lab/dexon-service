package custom

import (
	"github.com/google/btree"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/zuni-lab/dexon-service/pkg/db"
	"math/big"
	"sync"
)

var orderBook *OrderBook

type OrderBook struct {
	limitBuy  *OrderTree
	limitSell *OrderTree
	stopBuy   *OrderTree
	stopSell  *OrderTree
}

type OrderTree struct {
	btree.BTreeG[db.Order]
	mux sync.RWMutex
}

func GetOrderBook() *OrderBook {
	sync.OnceFunc(func() {
		orderBook = &OrderBook{
			limitBuy:  newOrderTree(),
			limitSell: newOrderTree(),
			stopBuy:   newOrderTree(),
			stopSell:  newOrderTree(),
		}
	})()

	return orderBook
}

func compareOrder(a, b db.Order) bool {
	sub := new(big.Int).Sub(a.Price.Int, b.Price.Int).Sign()
	if sub < 0 {
		return true
	} else if sub > 0 {
		return false
	} else {
		return a.CreatedAt.Time.After(b.CreatedAt.Time)
	}
}

func newOrderTree() *OrderTree {
	return &OrderTree{
		BTreeG: *btree.NewG[db.Order](32, compareOrder),
	}
}

func (o *OrderBook) Sub(side db.OrderSide, oType db.OrderType) *OrderTree {
	if side == db.OrderSideBUY && oType == db.OrderTypeLIMIT {
		return o.limitBuy
	} else if side == db.OrderSideBUY && oType == db.OrderTypeSTOP {
		return o.stopBuy
	} else if side == db.OrderSideSELL && oType == db.OrderTypeLIMIT {
		return o.limitSell
	} else if side == db.OrderSideSELL && oType == db.OrderTypeSTOP {
		return o.stopSell
	} else {
		return nil
	}
}

func (o *OrderTree) ReplaceOrInsert(order db.Order) {
	o.mux.Lock()
	defer o.mux.Unlock()
	o.ReplaceOrInsert(order)
}

func (o *OrderTree) Get(order db.Order) db.Order {
	o.mux.RLock()
	defer o.mux.RUnlock()
	return o.Get(order)
}

func (o *OrderTree) Delete(order db.Order) {
	o.mux.Lock()
	defer o.mux.Unlock()
	o.Delete(order)
}

func (o *OrderBook) Match(price pgtype.Numeric, priceTime pgtype.Timestamptz) []*db.Order {
	o.stopBuy.mux.RLock()
	o.limitBuy.mux.RLock()
	o.stopSell.mux.RLock()
	o.limitSell.mux.RLock()
	defer func() {
		o.stopBuy.mux.RUnlock()
		o.limitBuy.mux.RUnlock()
		o.stopSell.mux.RUnlock()
		o.limitSell.mux.RUnlock()
	}()

	var (
		stopBuyOrder   *db.Order
		limitSellOrder *db.Order
		limitBuyOrder  *db.Order
		stopSellOrder  *db.Order
		matchedOrder   *db.Order

		expected = db.Order{
			Price:     price,
			CreatedAt: priceTime,
		}
	)

	orderBook.Sub(db.OrderSideBUY, db.OrderTypeSTOP).AscendGreaterOrEqual(expected, func(order db.Order) bool {
		stopBuyOrder = &order
		return false
	})

	orderBook.Sub(db.OrderSideSELL, db.OrderTypeLIMIT).AscendGreaterOrEqual(expected, func(order db.Order) bool {
		limitSellOrder = &order
		return false
	})

	orderBook.Sub(db.OrderSideBUY, db.OrderTypeLIMIT).DescendLessOrEqual(expected, func(order db.Order) bool {
		limitBuyOrder = &order
		return false
	})

	orderBook.Sub(db.OrderSideSELL, db.OrderTypeSTOP).DescendLessOrEqual(expected, func(order db.Order) bool {
		stopSellOrder = &order
		return false
	})

	matchedOrder = stopBuyOrder
	if limitSellOrder != nil && limitSellOrder.CreatedAt.Time.Before(matchedOrder.CreatedAt.Time) {
		matchedOrder = limitSellOrder
	}
	if limitBuyOrder != nil && limitBuyOrder.CreatedAt.Time.Before(matchedOrder.CreatedAt.Time) {
		matchedOrder = limitBuyOrder
	}
	if stopSellOrder != nil && stopSellOrder.CreatedAt.Time.Before(matchedOrder.CreatedAt.Time) {
		matchedOrder = stopSellOrder
	}

	return []*db.Order{matchedOrder}
}
