package engine

import (
	"container/heap"
	"sort"
	"sync"
	"time"

	"github.com/Rishabhsingh78/orderMatchingEngine/pkg/utils"
)

type BidHeap []*Order

func (h BidHeap) Len() int { return len(h) }
func (h BidHeap) Less(i, j int) bool {
	if h[i].Price == h[j].Price {
		return h[i].Timestamp < h[j].Timestamp
	}
	return h[i].Price > h[j].Price
}
func (h BidHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
	h[i].HeapIndex = i
	h[j].HeapIndex = j
}
func (h *BidHeap) Push(x interface{}) {
	n := len(*h)
	item := x.(*Order)
	item.HeapIndex = n
	*h = append(*h, item)
}
func (h *BidHeap) Pop() interface{} {
	old := *h
	n := len(old)
	item := old[n-1]
	old[n-1] = nil // avoid memory leak
	item.HeapIndex = -1
	*h = old[0 : n-1]
	return item
}

type AskHeap []*Order

func (h AskHeap) Len() int { return len(h) }
func (h AskHeap) Less(i, j int) bool {
	// Lower price first. If prices equal, earlier timestamp first
	if h[i].Price == h[j].Price {
		return h[i].Timestamp < h[j].Timestamp
	}
	return h[i].Price < h[j].Price
}
func (h AskHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
	h[i].HeapIndex = i
	h[j].HeapIndex = j
}
func (h *AskHeap) Push(x interface{}) {
	n := len(*h)
	item := x.(*Order)
	item.HeapIndex = n
	*h = append(*h, item)
}
func (h *AskHeap) Pop() interface{} {
	old := *h
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	item.HeapIndex = -1
	*h = old[0 : n-1]
	return item
}

type OrderBook struct {
	Symbol            string
	Bids              BidHeap
	Asks              AskHeap
	Orders            map[string]*Order 
	TotalBidLiquidity int64
	TotalAskLiquidity int64
	mu                sync.RWMutex
}

func NewOrderBook(symbol string) *OrderBook {
	ob := &OrderBook{
		Symbol: symbol,
		Bids:   make(BidHeap, 0),
		Asks:   make(AskHeap, 0),
		Orders: make(map[string]*Order),
	}
	heap.Init(&ob.Bids)
	heap.Init(&ob.Asks)
	return ob
}

func (ob *OrderBook) ProcessOrder(order *Order) ([]Trade, error) {
	ob.mu.Lock()
	defer ob.mu.Unlock()

	if order.Quantity <= 0 {
		return nil, utils.ErrInvalidQuantity
	}
	if order.Type == OrderTypeLimit && order.Price <= 0 {
		return nil, utils.ErrInvalidPrice
	}

	var trades []Trade
	var err error

	if order.Side == SideBuy {
		trades, err = ob.matchBuyOrder(order)
	} else {
		trades, err = ob.matchSellOrder(order)
	}

	if err != nil {
		return nil, err
	}

	if order.Type == OrderTypeLimit && order.Quantity > order.Filled {
		order.Status = OrderStatusAccepted
		if order.Filled > 0 {
			order.Status = OrderStatusPartialFill
		}
		ob.addOrder(order)
	} else if order.Type == OrderTypeMarket && order.Quantity > order.Filled {
		
	}

	return trades, nil
}

func (ob *OrderBook) matchBuyOrder(order *Order) ([]Trade, error) {
	trades := []Trade{}
	if order.Type == OrderTypeMarket {
		if !ob.hasLiquidity(SideSell, order.Quantity) {
			return nil, utils.ErrInsufficientLiquidity
		}
	}

	for ob.Asks.Len() > 0 && order.Filled < order.Quantity {
		bestAsk := ob.Asks[0]

		if order.Type == OrderTypeLimit && order.Price < bestAsk.Price {
			break
		}

		matchQty := order.Quantity - order.Filled
		if matchQty > bestAsk.Quantity-bestAsk.Filled {
			matchQty = bestAsk.Quantity - bestAsk.Filled
		}

		trade := Trade{
			ID:           utils.GenerateUUID(), // We need a UUID generator
			Price:        bestAsk.Price,
			Quantity:     matchQty,
			Timestamp:    time.Now().UnixMilli(),
			MakerOrderID: bestAsk.ID,
			TakerOrderID: order.ID,
		}
		trades = append(trades, trade)

		// Update orders
		order.Filled += matchQty
		bestAsk.Filled += matchQty
		ob.TotalAskLiquidity -= matchQty

		// If bestAsk filled, remove it
		if bestAsk.Filled >= bestAsk.Quantity {
			bestAsk.Status = OrderStatusFilled
			heap.Pop(&ob.Asks)
			// delete(ob.Orders, bestAsk.ID) // Keep for history
		} else {
			bestAsk.Status = OrderStatusPartialFill
		}
	}

	if order.Filled >= order.Quantity {
		order.Status = OrderStatusFilled
	}

	return trades, nil
}

func (ob *OrderBook) matchSellOrder(order *Order) ([]Trade, error) {
	trades := []Trade{}

	// For Market Order, check liquidity first
	if order.Type == OrderTypeMarket {
		if !ob.hasLiquidity(SideBuy, order.Quantity) {
			return nil, utils.ErrInsufficientLiquidity
		}
	}

	for ob.Bids.Len() > 0 && order.Filled < order.Quantity {
		bestBid := ob.Bids[0]

		// Price check for Limit orders
		if order.Type == OrderTypeLimit && order.Price > bestBid.Price {
			break
		}

		// Match
		matchQty := order.Quantity - order.Filled
		if matchQty > bestBid.Quantity-bestBid.Filled {
			matchQty = bestBid.Quantity - bestBid.Filled
		}

		trade := Trade{
			ID:           utils.GenerateUUID(),
			Price:        bestBid.Price,
			Quantity:     matchQty,
			Timestamp:    time.Now().UnixMilli(),
			MakerOrderID: bestBid.ID,
			TakerOrderID: order.ID,
		}
		trades = append(trades, trade)

		order.Filled += matchQty
		bestBid.Filled += matchQty
		ob.TotalBidLiquidity -= matchQty

		if bestBid.Filled >= bestBid.Quantity {
			bestBid.Status = OrderStatusFilled
			heap.Pop(&ob.Bids)
		} else {
			bestBid.Status = OrderStatusPartialFill
		}
	}

	if order.Filled >= order.Quantity {
		order.Status = OrderStatusFilled
	}

	return trades, nil
}

func (ob *OrderBook) addOrder(order *Order) {
	ob.Orders[order.ID] = order
	remaining := order.Quantity - order.Filled
	if order.Side == SideBuy {
		ob.TotalBidLiquidity += remaining
		heap.Push(&ob.Bids, order)
	} else {
		ob.TotalAskLiquidity += remaining
		heap.Push(&ob.Asks, order)
	}
}

func (ob *OrderBook) CancelOrder(orderID string) error {
	ob.mu.Lock()
	defer ob.mu.Unlock()

	order, ok := ob.Orders[orderID]
	if !ok {
		return utils.ErrOrderNotFound
	}

	remaining := order.Quantity - order.Filled
	if order.Side == SideBuy {
		ob.TotalBidLiquidity -= remaining
		heap.Remove(&ob.Bids, order.HeapIndex)
	} else {
		ob.TotalAskLiquidity -= remaining
		heap.Remove(&ob.Asks, order.HeapIndex)
	}

	order.Status = OrderStatusCancelled
	return nil
}

func (ob *OrderBook) hasLiquidity(side Side, quantity int64) bool {
	if side == SideBuy {
		return ob.TotalBidLiquidity >= quantity
	}
	return ob.TotalAskLiquidity >= quantity
}

type PriceLevel struct {
	Price    int64 `json:"price"`
	Quantity int64 `json:"quantity"`
}

type OrderBookSnapshot struct {
	Symbol    string       `json:"symbol"`
	Timestamp int64        `json:"timestamp"`
	Bids      []PriceLevel `json:"bids"`
	Asks      []PriceLevel `json:"asks"`
}

func (ob *OrderBook) GetSnapshot(depth int) OrderBookSnapshot {
	ob.mu.RLock()
	defer ob.mu.RUnlock()

	snapshot := OrderBookSnapshot{
		Symbol:    ob.Symbol,
		Timestamp: time.Now().UnixMilli(),
		Bids:      make([]PriceLevel, 0),
		Asks:      make([]PriceLevel, 0),
	}

	bidMap := make(map[int64]int64)
	for _, order := range ob.Bids {
		bidMap[order.Price] += (order.Quantity - order.Filled)
	}
	for price, qty := range bidMap {
		snapshot.Bids = append(snapshot.Bids, PriceLevel{Price: price, Quantity: qty})
	}
	sortPriceLevels(snapshot.Bids, true)

	askMap := make(map[int64]int64)
	for _, order := range ob.Asks {
		askMap[order.Price] += (order.Quantity - order.Filled)
	}
	for price, qty := range askMap {
		snapshot.Asks = append(snapshot.Asks, PriceLevel{Price: price, Quantity: qty})
	}
	sortPriceLevels(snapshot.Asks, false)

	if len(snapshot.Bids) > depth {
		snapshot.Bids = snapshot.Bids[:depth]
	}
	if len(snapshot.Asks) > depth {
		snapshot.Asks = snapshot.Asks[:depth]
	}

	return snapshot
}

func sortPriceLevels(levels []PriceLevel, descending bool) {
	sort.Slice(levels, func(i, j int) bool {
		if descending {
			return levels[i].Price > levels[j].Price
		}
		return levels[i].Price < levels[j].Price
	})
}
