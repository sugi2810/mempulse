package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/sugi2810/mempulse/api"
	"github.com/sugi2810/mempulse/chains"
	"github.com/sugi2810/mempulse/rpc"
)

func main() {
	alchemyKey := os.Getenv("ALCHEMY_KEY")
	if alchemyKey == "" {
		log.Fatal("ALCHEMY_KEY environment variable not set")
	}

	httpEndpoint := "https://eth-mainnet.g.alchemy.com/v2/" + alchemyKey
	wssEndpoint := "wss://eth-mainnet.g.alchemy.com/v2/" + alchemyKey

	ethClient := rpc.NewClient(httpEndpoint)
	listener := chains.NewMempoolListener(wssEndpoint, httpEndpoint)

	ctx := context.Background()
	listener.Start(ctx)

	http.HandleFunc("/health", api.HealthHandler)
	http.HandleFunc("/chains", api.ChainsHandler)
	http.HandleFunc("/mempool", api.MempoolHandler)
	http.HandleFunc("/ethereum", api.EthereumHandler(ethClient))
	http.HandleFunc("/live", api.LiveMempoolHandler(listener))

	fmt.Println("Mempulse API running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
