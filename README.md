# Order Matching Engine
A high-performance order matching engine written in Go.

## Features
- Limit and Market Orders
- Price-Time Priority Matching
- In-memory Order Book
- REST API
- Concurrency Safe

## Requirements
- Go 1.21+

## Running the Server
```bash
go run cmd/server/main.go
```
The server will start on port 8080.

## API Endpoints

### Submit Order
`POST /api/v1/orders`
```json
{
  "symbol": "AAPL",
  "side": "BUY",
  "type": "LIMIT",
  "price": 15050,
  "quantity": 100
}
```

### Cancel Order
`DELETE /api/v1/orders/{order_id}`

### Get Order Book
`GET /api/v1/orderbook/{symbol}?depth=10`

### Get Order Status
`GET /api/v1/orders/{order_id}`

