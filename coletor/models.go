package main

type Metrics struct {
	CPU        float64 `json:"cpu"`
	RAM        float64 `json:"ram"`
	QueueSize  int     `json:"queue_size"`
	ErrorCount int     `json:"error_count"`
}
