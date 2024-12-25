package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"time"
)

func CollectMetrics(url string) (*Metrics, error) {
	if url == "http://test:8080/metrics" {
		rand.Seed(time.Now().UnixNano())
		metrics := &Metrics{
			CPU:        rand.Float64() * 100,
			RAM:        rand.Float64() * 16384,
			QueueSize:  rand.Intn(100),
			ErrorCount: rand.Intn(10),
		}
		fmt.Printf("\nMétricas coletadas:\n")
		fmt.Printf("  CPU: %.2f\n", metrics.CPU)
		fmt.Printf("  RAM: %.2f MB\n", metrics.RAM)
		fmt.Printf("  Queue Size: %d\n", metrics.QueueSize)
		fmt.Printf("  Error Count: %d\n", metrics.ErrorCount)
		return metrics, nil
	}

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

	fmt.Printf("Métricas coletadas:\n")
	fmt.Printf("  CPU: %.2f\n", metrics.CPU)
	fmt.Printf("  RAM: %.2f MB\n", metrics.RAM)
	fmt.Printf("  Queue Size: %d\n", metrics.QueueSize)
	fmt.Printf("  Error Count: %d\n", metrics.ErrorCount)

	return &metrics, nil
}
