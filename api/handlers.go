package api

import (
	"encoding/json"
	"net/http"
	"time"
)

type HealthResponse struct {
	Status    string `json:"status"`
	Version   string `json:"version"`
	Timestamp string `json:"timestamp"`
}

type ChainStatus struct {
	Chain  string `json:"chain"`
	Status string `json:"status"`
	Nodes  int    `json:"nodes"`
}

type ChainsResponse struct {
	Chains []ChainStatus `json:"chains"`
}

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	resp := HealthResponse{
		Status:    "ok",
		Version:   "0.1.0",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}
	json.NewEncoder(w).Encode(resp)
}

func ChainsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	resp := ChainsResponse{
		Chains: []ChainStatus{
			{Chain: "ethereum", Status: "active", Nodes: 5},
			{Chain: "solana", Status: "active", Nodes: 2},
			{Chain: "polygon", Status: "standby", Nodes: 3},
		},
	}
	json.NewEncoder(w).Encode(resp)
}