package apis

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/Rishabhsingh78/orderMatchingEngine/internals/engine"
	"github.com/Rishabhsingh78/orderMatchingEngine/pkg/utils"
	"github.com/gorilla/mux"
)

type Handler struct {
	Engine *engine.Engine
}

func NewHandler(e *engine.Engine) *Handler {
	return &Handler{Engine: e}
}

func (h *Handler) SubmitOrder(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Symbol   string           `json:"symbol"`
		Side     engine.Side      `json:"side"`
		Type     engine.OrderType `json:"type"`
		Price    int64            `json:"price"`
		Quantity int64            `json:"quantity"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Malformed JSON")
		return
	}

	if req.Quantity <= 0 {
		writeError(w, http.StatusBadRequest, "Invalid order: quantity must be positive")
		return
	}
	if req.Type == engine.OrderTypeLimit && req.Price <= 0 {
		writeError(w, http.StatusBadRequest, "Invalid order: price must be positive")
		return
	}

	order := &engine.Order{
		ID:        utils.GenerateUUID(),
		Symbol:    req.Symbol,
		Side:      req.Side,
		Type:      req.Type,
		Price:     req.Price,
		Quantity:  req.Quantity,
		Timestamp: time.Now().UnixMilli(),
		Status:    engine.OrderStatusAccepted,
	}

	trades, err := h.Engine.SubmitOrder(order)
	if err != nil {
		if err == utils.ErrInsufficientLiquidity {
			writeError(w, http.StatusBadRequest, "Insufficient liquidity")
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	resp := OrderResponse{
		OrderID:           order.ID,
		Status:            order.Status,
		FilledQuantity:    order.Filled,
		RemainingQuantity: order.Quantity - order.Filled,
		Trades:            trades,
	}

	if order.Status == engine.OrderStatusAccepted {
		resp.Message = "Order added to book"
		writeJSON(w, http.StatusCreated, resp)
	} else if order.Status == engine.OrderStatusFilled {
		writeJSON(w, http.StatusOK, resp)
	} else {
		writeJSON(w, http.StatusAccepted, resp)
	}
}

func (h *Handler) CancelOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orderID := vars["order_id"]

	err := h.Engine.CancelOrder(orderID)
	if err != nil {
		if err == utils.ErrOrderNotFound {
			writeError(w, http.StatusNotFound, "Order not found")
			return
		}
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"order_id": orderID,
		"status":   "CANCELLED",
	})
}

func (h *Handler) GetOrderBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	symbol := vars["symbol"]
	depthStr := r.URL.Query().Get("depth")
	depth := 10
	if depthStr != "" {
		if d, err := strconv.Atoi(depthStr); err == nil && d > 0 {
			depth = d
		}
	}

	ob := h.Engine.GetOrderBook(symbol)

	snapshot := ob.GetSnapshot(depth)

	resp := OrderBookResponse{
		Symbol:    snapshot.Symbol,
		Timestamp: snapshot.Timestamp,
		Bids:      make([]PriceLevel, len(snapshot.Bids)),
		Asks:      make([]PriceLevel, len(snapshot.Asks)),
	}
	for i, b := range snapshot.Bids {
		resp.Bids[i] = PriceLevel{Price: b.Price, Quantity: b.Quantity}
	}
	for i, a := range snapshot.Asks {
		resp.Asks[i] = PriceLevel{Price: a.Price, Quantity: a.Quantity}
	}

	writeJSON(w, http.StatusOK, resp)
}

func (h *Handler) GetOrderStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orderID := vars["order_id"]

	order, err := h.Engine.GetOrder(orderID)
	if err != nil {
		writeError(w, http.StatusNotFound, "Order not found")
		return
	}

	resp := OrderStatusResponse{
		OrderID:        order.ID,
		Symbol:         order.Symbol,
		Side:           order.Side,
		Type:           order.Type,
		Price:          order.Price,
		Quantity:       order.Quantity,
		FilledQuantity: order.Filled,
		Status:         order.Status,
		Timestamp:      order.Timestamp,
	}
	writeJSON(w, http.StatusOK, resp)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, ErrorResponse{Error: message})
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
