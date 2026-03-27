package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/sugi2810/mempulse/chains"
)

type LiveMempoolResponse struct {
	Chain        string           `json:"chain"`
	Timestamp    string           `json:"timestamp"`
	PendingCount int              `json:"pendingCount"`
	Transactions []chains.PendingTx `json:"transactions"`
}

func LiveMempoolHandler(listener *chains.MempoolListener) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		txs := listener.GetTransactions()

		json.NewEncoder(w).Encode(LiveMempoolResponse{
			Chain:        "ethereum",
			Timestamp:    time.Now().UTC().Format(time.RFC3339),
			PendingCount: len(txs),
			Transactions: txs,
		})
	}
}
