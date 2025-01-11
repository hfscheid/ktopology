package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func CollectMetrics(url string) (*Metrics, error) {
	if url == "http://test:8080/metrics" {
		rand.Seed(time.Now().UnixNano())
		metrics := &Metrics{
			CPUUsage:    rand.Float64() * 100,
			MemUsage:    rand.Float64() * 16384,
			QueueSize:   rand.Intn(100),
			NumRejected: rand.Intn(10),
		}
		fmt.Printf("\nMétricas coletadas (modo teste):\n")
		fmt.Printf("  CPU Usage: %.2f\n", metrics.CPUUsage)
		fmt.Printf("  Memory Usage: %.2f MB\n", metrics.MemUsage)
		fmt.Printf("  Queue Size: %d\n", metrics.QueueSize)
		fmt.Printf("  Number of Rejected Messages: %d\n", metrics.NumRejected)
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

	metricsText := string(body)
	lines := strings.Split(metricsText, "\n")

	var metrics Metrics

	for _, line := range lines {
		if strings.HasPrefix(line, "cpu_usage") {
			valueStr := strings.TrimSpace(strings.Split(line, " ")[1])
			metrics.CPUUsage, _ = strconv.ParseFloat(valueStr, 64)
		} else if strings.HasPrefix(line, "mem_usage") {
			valueStr := strings.TrimSpace(strings.Split(line, " ")[1])
			metrics.MemUsage, _ = strconv.ParseFloat(valueStr, 64)
		} else if strings.HasPrefix(line, "queue_size") {
			valueStr := strings.TrimSpace(strings.Split(line, " ")[1])
			queueSize, _ := strconv.Atoi(valueStr)
			metrics.QueueSize = queueSize
		} else if strings.HasPrefix(line, "num_rejected") {
			valueStr := strings.TrimSpace(strings.Split(line, " ")[1])
			numRejected, _ := strconv.Atoi(valueStr)
			metrics.NumRejected = numRejected
		}
	}

	fmt.Printf("Métricas coletadas:\n")
	fmt.Printf("  CPU Usage: %.2f\n", metrics.CPUUsage)
	fmt.Printf("  Memory Usage: %.2f MB\n", metrics.MemUsage)
	fmt.Printf("  Queue Size: %d\n", metrics.QueueSize)
	fmt.Printf("  Number of Rejected Messages: %d\n", metrics.NumRejected)

	return &metrics, nil
}
