package apis

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Rishabhsingh78/orderMatchingEngine/internals/engine"
	"github.com/gorilla/mux"
)

func TestGetOrderBook(t *testing.T) {
	e := engine.NewEngine()
	h := NewHandler(e)
	router := mux.NewRouter()
	router.HandleFunc("/orderbook/{symbol}", h.GetOrderBook).Methods("GET")

	e.GetOrderBook("BTCUSD")
	req, _ := http.NewRequest("GET", "/orderbook/BTCUSD", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	if ctype := rr.Header().Get("Content-Type"); ctype != "application/json" {
		t.Errorf("content type header does not match: got %v want %v",
			ctype, "application/json")
	}

	var resp OrderBookResponse
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Errorf("handler returned invalid body: %v", err)
	}

	if resp.Symbol != "BTCUSD" {
		t.Errorf("handler returned wrong symbol: got %v want %v",
			resp.Symbol, "BTCUSD")
	}
}
