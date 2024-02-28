package server

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

func (server *Server) Route() {
	/**
	 * Build endpoint for exposing intraday data
	 */
	server.router.HandleFunc("/stock/{ticker}", server.getStockInfoHandler).Methods("GET")
}

func (server *Server) getStockInfoHandler(w http.ResponseWriter, r *http.Request) {
	// Extract ticker from the request path
	vars := mux.Vars(r)
	ticker := vars["ticker"]
	if strings.TrimSpace(ticker) == "" {
		server.logger.Error("Empty ticker specified")
		http.Error(w, "Empty stock ticker", http.StatusBadRequest)
		return
	}

	// Query the database for the stock info
	stockInfo := make([]StockInfo, 0)
	err := server.db.SelectContext(context.Background(), &stockInfo, "SELECT * FROM intraday WHERE ticker = $1", ticker)
	if err != nil {
		server.logger.WithError(err).Error("Failed to fetch stock info from database")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// return response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stockInfo)
}
