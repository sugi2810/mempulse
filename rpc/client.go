package rpc

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const defaultTimeout = 10 * time.Second

type Client struct {
	endpoint string
	http     *http.Client
}

func NewClient(endpoint string) *Client {
	return &Client{
		endpoint: endpoint,
		http:     &http.Client{Timeout: 10 * time.Second},
	}
}

type rpcRequest struct {
	Jsonrpc string `json:"jsonrpc"`
	Method  string `json:"method"`
	Params  []any  `json:"params"`
	ID      int    `json:"id"`
}

type rpcResponse struct {
	Jsonrpc string          `json:"jsonrpc"`
	ID      int             `json:"id"`
	Result  json.RawMessage `json:"result"`
	Error   *rpcError       `json:"error,omitempty"`
}

type rpcError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (c *Client) call(ctx context.Context, method string, params []any) (json.RawMessage, error) {
	req := rpcRequest{
		Jsonrpc: "2.0",
		Method:  method,
		Params:  params,
		ID:      1,
	}

	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("rpc marshal: %w", err)
	}

	// Create the HTTP request with context
	httpReq, err := http.NewRequestWithContext(ctx, "POST", c.endpoint, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("rpc request: %w", err)
	}
	// Set the headers so the endpoint knows how to handle the payload
	httpReq.Header.Set("Content-Type", "application/json")
	// Execute the request
	resp, err := c.http.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("rpc post: %w", err)
	}
	defer resp.Body.Close()

	rawBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("rpc read body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("rpc status %d: %s", resp.StatusCode, string(rawBody))
	}

	var rpcResp rpcResponse
	if err := json.Unmarshal(rawBody, &rpcResp); err != nil {
		return nil, fmt.Errorf("rpc decode: %w", err)
	}

	if rpcResp.Error != nil {
		return nil, fmt.Errorf("rpc error %d: %s", rpcResp.Error.Code, rpcResp.Error.Message)
	}

	return rpcResp.Result, nil
}

func (c *Client) GetBlockNumber(ctx context.Context) (string, error) {
	// Pass ctx as the first argument to the internal call method
	raw, err := c.call(ctx, "eth_blockNumber", []any{})
	if err != nil {
		return "", fmt.Errorf("GetBlockNumber: %w", err)
	}
	var blockNumber string
	if err := json.Unmarshal(raw, &blockNumber); err != nil {
		return "", fmt.Errorf("unmarshal block number: %w", err)
	}
	return blockNumber, nil
}

// GetGasPrice returns the current price per gas in wei.
func (c *Client) GetGasPrice(ctx context.Context) (string, error) {
	// Pass the context into our internal call method
	raw, err := c.call(ctx, "eth_gasPrice", []any{})
	if err != nil {
		return "", fmt.Errorf("get gas price: %w", err)
	}
	var gasPrice string
	// Unmarshal the JSON-RPC result field into our string
	if err := json.Unmarshal(raw, &gasPrice); err != nil {
		return "", fmt.Errorf("unmarshal gas price: %w", err)
	}
	return gasPrice, nil
}