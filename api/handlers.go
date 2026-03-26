package api

import (
	"encoding/json"
	"net/http"
	"time"
)

// --- Existing types ---

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

// --- New types for Day 12 ---

type Transaction struct {
	Hash     string `json:"hash"`
	From     string `json:"from"`
	To       string `json:"to,omitempty"`
	Gas      uint64 `json:"gas"`
	GasPrice string `json:"gasPrice,omitempty"`
}

type MempoolSnapshot struct {
	Chain        string        `json:"chain"`
	Timestamp    string        `json:"timestamp"`
	PendingCount int           `json:"pendingCount"`
	Transactions []Transaction `json:"transactions"`
}

// --- Existing handlers ---

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

// --- New handler for Day 12 ---

func MempoolHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	snapshot := MempoolSnapshot{
		Chain:        "ethereum",
		Timestamp:    time.Now().UTC().Format(time.RFC3339),
		PendingCount: 3,
		Transactions: []Transaction{
			{Hash: "0xabc123", From: "0xSender1", Gas: 21000, GasPrice: "50 gwei"},
			{Hash: "0xdef456", From: "0xSender2", Gas: 42000, GasPrice: "45 gwei"},
			{Hash: "0xghi789", From: "0xSender3", Gas: 21000, GasPrice: "55 gwei"},
		},
	}
	json.NewEncoder(w).Encode(snapshot)
}