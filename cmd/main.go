package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/sugi2810/mempulse/api"
)

func main() {
	http.HandleFunc("/health", api.HealthHandler)
	http.HandleFunc("/chains", api.ChainsHandler)
	http.HandleFunc("/mempool", api.MempoolHandler)

	fmt.Println("Mempulse API running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
