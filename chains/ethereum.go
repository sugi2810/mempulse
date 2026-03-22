package chains

type EthereumClient struct {
	URL       string
	Connected bool
}

func (e *EthereumClient) HealthStatus() string {
	return "Ethereum: Active"
}

func (e *EthereumClient) NodeCount() int {
	return 5
}