package apis

import "github.com/Rishabhsingh78/orderMatchingEngine/internals/engine"

type OrderResponse struct {
	OrderID           string             `json:"order_id"`
	Status            engine.OrderStatus `json:"status"`
	Message           string             `json:"message,omitempty"`
	FilledQuantity    int64              `json:"filled_quantity,omitempty"`
	RemainingQuantity int64              `json:"remaining_quantity,omitempty"`
	Trades            []engine.Trade     `json:"trades,omitempty"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type OrderBookResponse struct {
	Symbol    string       `json:"symbol"`
	Timestamp int64        `json:"timestamp"`
	Bids      []PriceLevel `json:"bids"`
	Asks      []PriceLevel `json:"asks"`
}

type PriceLevel struct {
	Price    int64 `json:"price"`
	Quantity int64 `json:"quantity"`
}

type OrderStatusResponse struct {
	OrderID        string             `json:"order_id"`
	Symbol         string             `json:"symbol"`
	Side           engine.Side        `json:"side"`
	Type           engine.OrderType   `json:"type"`
	Price          int64              `json:"price"`
	Quantity       int64              `json:"quantity"`
	FilledQuantity int64              `json:"filled_quantity"`
	Status         engine.OrderStatus `json:"status"`
	Timestamp      int64              `json:"timestamp"`
}
