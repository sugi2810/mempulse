package rpc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Client struct {
	endpoint string
	http     *http.Client
}

func NewClient(endpoint string) *Client {
	return &Client{
		endpoint: endpoint,
		http:     &http.Client{},
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

func (c *Client) call(method string, params []any) (json.RawMessage, error) {
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

	resp, err := c.http.Post(c.endpoint, "application/json", bytes.NewReader(body))
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
		return nil, fmt.Errorf("rpc decode: %s — body was: %s", err.Error(), string(rawBody))
	}

	if rpcResp.Error != nil {
		return nil, fmt.Errorf("rpc error %d: %s", rpcResp.Error.Code, rpcResp.Error.Message)
	}

	return rpcResp.Result, nil
}

func (c *Client) GetBlockNumber() (string, error) {
	result, err := c.call("eth_blockNumber", []any{})
	if err != nil {
		return "", fmt.Errorf("GetBlockNumber: %w", err)
	}

	var blockNumber string
	if err := json.Unmarshal(result, &blockNumber); err != nil {
		return "", fmt.Errorf("GetBlockNumber unmarshal: %w", err)
	}

	return blockNumber, nil
}

func (c *Client) GetGasPrice() (string, error) {
	result, err := c.call("eth_gasPrice", []any{})
	if err != nil {
		return "", fmt.Errorf("GetGasPrice: %w", err)
	}

	var gasPrice string
	if err := json.Unmarshal(result, &gasPrice); err != nil {
		return "", fmt.Errorf("GetGasPrice unmarshal: %w", err)
	}

	return gasPrice, nil
}
