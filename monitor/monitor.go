package monitor

type MempoolMonitor interface {
	HealthStatus() string
	NodeCount()    int
}