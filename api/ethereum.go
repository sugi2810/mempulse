package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/sugi2810/mempulse/rpc"
)

type EthereumStatus struct {
	Chain       string `json:"chain"`
	Timestamp   string `json:"timestamp"`
	BlockNumber string `json:"blockNumber"`
	GasPrice    string `json:"gasPrice"`
	Status      string `json:"status"`
}

func EthereumHandler(client *rpc.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		blockNumber, err := client.GetBlockNumber()
		if err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			json.NewEncoder(w).Encode(map[string]string{
				"error": err.Error(),
			})
			return
		}

		gasPrice, err := client.GetGasPrice()
		if err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			json.NewEncoder(w).Encode(map[string]string{
				"error": err.Error(),
			})
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(EthereumStatus{
			Chain:       "ethereum",
			Timestamp:   time.Now().UTC().Format(time.RFC3339),
			BlockNumber: blockNumber,
			GasPrice:    gasPrice,
			Status:      "active",
		})
	}
}