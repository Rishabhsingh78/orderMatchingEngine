package engine

import (
	"sync"

	"github.com/Rishabhsingh78/orderMatchingEngine/pkg/utils"
)


type Engine struct {
	OrderBooks       map[string]*OrderBook
	OrderSymbolIndex map[string]string 
	mu               sync.RWMutex
}

func NewEngine() *Engine {
	return &Engine{
		OrderBooks:       make(map[string]*OrderBook),
		OrderSymbolIndex: make(map[string]string),
	}
}

func (e *Engine) GetOrderBook(symbol string) *OrderBook {
	e.mu.Lock()
	defer e.mu.Unlock()

	ob, exists := e.OrderBooks[symbol]
	if !exists {
		ob = NewOrderBook(symbol)
		e.OrderBooks[symbol] = ob
	}
	return ob
}

func (e *Engine) SubmitOrder(order *Order) ([]Trade, error) {
	if order.Symbol == "" {
		return nil, utils.ErrInvalidSymbol
	}

	e.mu.Lock()
	e.OrderSymbolIndex[order.ID] = order.Symbol
	e.mu.Unlock()

	ob := e.GetOrderBook(order.Symbol)
	return ob.ProcessOrder(order)
}

func (e *Engine) CancelOrder(orderID string) error {
	e.mu.RLock()
	symbol, exists := e.OrderSymbolIndex[orderID]
	e.mu.RUnlock()

	if !exists {
		return utils.ErrOrderNotFound
	}

	ob := e.GetOrderBook(symbol)
	return ob.CancelOrder(orderID)
}

func (e *Engine) GetOrder(orderID string) (*Order, error) {
	e.mu.RLock()
	symbol, exists := e.OrderSymbolIndex[orderID]
	e.mu.RUnlock()

	if !exists {
		return nil, utils.ErrOrderNotFound
	}

	ob := e.GetOrderBook(symbol)
	ob.mu.RLock()
	defer ob.mu.RUnlock()

	order, ok := ob.Orders[orderID]
	if !ok {
		return nil, utils.ErrOrderNotFound
	}
	return order, nil
}
