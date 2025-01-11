package main

type Metrics struct {
	CPUUsage    float64 `json:"cpu_usage"`
	MemUsage    float64 `json:"mem_usage"`
	QueueSize   int     `json:"queue_size"`
	NumRejected int     `json:"num_rejected"`
}
