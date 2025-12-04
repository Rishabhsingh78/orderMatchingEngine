package apis

import (
	"net/http"

	"github.com/gorilla/mux"
)

func NewRouter(h *Handler) *mux.Router {
	r := mux.NewRouter()
	api := r.PathPrefix("/api/v1").Subrouter()

	api.HandleFunc("/orders", h.SubmitOrder).Methods(http.MethodPost)
	api.HandleFunc("/orders/{order_id}", h.CancelOrder).Methods(http.MethodDelete)
	api.HandleFunc("/orders/{order_id}", h.GetOrderStatus).Methods(http.MethodGet)
	api.HandleFunc("/orderbook/{symbol}", h.GetOrderBook).Methods(http.MethodGet)

	// Health check
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "healthy"})
	}).Methods(http.MethodGet)

	return r
}
