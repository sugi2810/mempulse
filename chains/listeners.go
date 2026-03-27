package chains

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

const MaxTransactions = 100

type PendingTx struct {
	Hash     string `json:"hash"`
	Gas      string `json:"gas,omitempty"`
	GasPrice string `json:"gasPrice,omitempty"`
	Value    string `json:"value,omitempty"`
	To       string `json:"to,omitempty"`
}

type MempoolListener struct {
	mu           sync.RWMutex
	transactions []PendingTx
	wssEndpoint  string
	httpEndpoint string
}

func NewMempoolListener(wssEndpoint string, httpEndpoint string) *MempoolListener {
	return &MempoolListener{
		wssEndpoint:  wssEndpoint,
		httpEndpoint: httpEndpoint,
		transactions: make([]PendingTx, 0, MaxTransactions),
	}
}

type wsRequest struct {
	Jsonrpc string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
	ID      int           `json:"id"`
}

type wsResponse struct {
	Params *wsParams `json:"params,omitempty"`
}

type wsParams struct {
	Result json.RawMessage `json:"result"`
}

func (m *MempoolListener) fetchTx(hash string) (*PendingTx, error) {
	reqBody, _ := json.Marshal(map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "eth_getTransactionByHash",
		"params":  []string{hash},
		"id":      1,
	})

	resp, err := http.Post(m.httpEndpoint, "application/json", bytes.NewReader(reqBody))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var result struct {
		Result *PendingTx `json:"result"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return result.Result, nil
}

func (m *MempoolListener) Start(ctx context.Context) {
	go func() {
		conn, _, err := websocket.DefaultDialer.Dial(m.wssEndpoint, nil)
		if err != nil {
			log.Printf("mempool listener: failed to connect: %v", err)
			return
		}
		defer conn.Close()

		req := wsRequest{
			Jsonrpc: "2.0",
			Method:  "eth_subscribe",
			Params:  []interface{}{"newPendingTransactions"},
			ID:      1,
		}

		if err := conn.WriteJSON(req); err != nil {
			log.Printf("mempool listener: subscribe failed: %v", err)
			return
		}

		fmt.Println("Mempulse: listening to Ethereum mempool...")

		for {
			select {
			case <-ctx.Done():
				return
			default:
				_, msg, err := conn.ReadMessage()
				if err != nil {
					log.Printf("mempool listener: read error: %v", err)
					return
				}

				var resp wsResponse
				if err := json.Unmarshal(msg, &resp); err != nil {
					continue
				}

				if resp.Params == nil {
					continue
				}

				var hash string
				if err := json.Unmarshal(resp.Params.Result, &hash); err != nil {
					continue
				}

				if hash == "" {
					continue
				}

				// fetch full tx in background
				go func(h string) {
					tx, err := m.fetchTx(h)
					if err != nil || tx == nil {
						return
					}
					tx.Hash = h
					m.addTransaction(*tx)
				}(hash)
			}
		}
	}()
}

func (m *MempoolListener) addTransaction(tx PendingTx) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if len(m.transactions) >= MaxTransactions {
		m.transactions = m.transactions[1:]
	}
	m.transactions = append(m.transactions, tx)
}

func (m *MempoolListener) GetTransactions() []PendingTx {
	m.mu.RLock()
	defer m.mu.RUnlock()

	result := make([]PendingTx, len(m.transactions))
	copy(result, m.transactions)
	return result
}