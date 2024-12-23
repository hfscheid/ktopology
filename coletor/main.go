package main

import (
	"log"
	"time"
)

const PollInterval = 10 * time.Second

func main() {
	nodes := []string{
		"http://node1/metrics",
		"http://node2/metrics",
	}

	for {
		for _, node := range nodes {
			go func(node string) {
				metrics, err := CollectMetrics(node)
				if err != nil {
					log.Printf("Erro ao coletar métricas de %s: %v", node, err)
					return
				}

				err = StoreMetrics(metrics)
				if err != nil {
					log.Printf("Erro ao armazenar métricas de %s: %v", node, err)
				}
			}(node)
		}
		time.Sleep(PollInterval)
	}
}
