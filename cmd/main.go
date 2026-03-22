package main

import (
	"fmt"
	"github.com/sugi2810/mempulse/chains"
	"github.com/sugi2810/mempulse/monitor"
)

func checkHealth(m monitor.MempoolMonitor, ch chan string) {
	ch <- m.HealthStatus()
}

func main() {
	clients := []monitor.MempoolMonitor{
		&chains.EthereumClient{URL: "https://eth-rpc.example", Connected: true},
	}

	ch := make(chan string, len(clients))

	for _, c := range clients {
		go checkHealth(c, ch)
	}

	for range clients {
		fmt.Println(<-ch)
	}
}