package engine


type Side string

const (
	SideBuy  Side = "BUY"
	SideSell Side = "SELL"
)

type OrderType string

const (
	OrderTypeLimit  OrderType = "LIMIT"
	OrderTypeMarket OrderType = "MARKET"
)

type OrderStatus string

const (
	OrderStatusAccepted    OrderStatus = "ACCEPTED"
	OrderStatusPartialFill OrderStatus = "PARTIAL_FILL"
	OrderStatusFilled      OrderStatus = "FILLED"
	OrderStatusCancelled   OrderStatus = "CANCELLED"
	OrderStatusRejected    OrderStatus = "REJECTED"
)

type Order struct {
	ID        string      `json:"id"`
	Symbol    string      `json:"symbol"`
	Side      Side        `json:"side"`
	Type      OrderType   `json:"type"`
	Price     int64       `json:"price"` // Price in cents
	Quantity  int64       `json:"quantity"`
	Timestamp int64       `json:"timestamp"` // Unix milliseconds
	Filled    int64       `json:"filled_quantity"`
	Status    OrderStatus `json:"status"`
	HeapIndex int `json:"-"`
}

type Trade struct {
	ID           string `json:"trade_id"`
	Price        int64  `json:"price"`
	Quantity     int64  `json:"quantity"`
	Timestamp    int64  `json:"timestamp"`
	MakerOrderID string `json:"maker_order_id"`
	TakerOrderID string `json:"taker_order_id"`
}
