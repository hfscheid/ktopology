package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type Metrics struct {
	CPU        float64 `json:"cpu"`
	RAM        float64 `json:"ram"`
	QueueSize  int     `json:"queue_size"`
	ErrorCount int     `json:"error_count"`
}

func CollectMetrics(url string) (*Metrics, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var metrics Metrics

	err = json.Unmarshal(body, &metrics)
	if err != nil {
		return nil, err
	}

	return &metrics, nil
}
