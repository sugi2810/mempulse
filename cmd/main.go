package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/sugi2810/mempulse/api"
	"github.com/sugi2810/mempulse/rpc"
)

func main() {
	infuraKey := os.Getenv("INFURA_KEY")
	if infuraKey == "" {
		log.Fatal("INFURA_KEY environment variable not set")
	}

	ethClient := rpc.NewClient("https://mainnet.infura.io/v3/" + infuraKey)

	http.HandleFunc("/health", api.HealthHandler)
	http.HandleFunc("/chains", api.ChainsHandler)
	http.HandleFunc("/mempool", api.MempoolHandler)
	http.HandleFunc("/ethereum", api.EthereumHandler(ethClient))

	fmt.Println("Mempulse API running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
