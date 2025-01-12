package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

  "collector/ktprom"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

type Node struct {
	ID       string                 `json:"id"`
	Label    string                 `json:"label"`
	Metadata map[string]interface{} `json:"metadata"`
}

type Edge struct {
	Source string `json:"source"`
	Target string `json:"target"`
}

type Graph struct {
	Nodes []Node `json:"nodes"`
	Edges []Edge `json:"edges"`
}

func init() {
	var err error

	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}

	clientOptions := options.Client().ApplyURI(mongoURI)
	client, err = mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
}

func buildGraphFromMerics(metrics []*ktprom.TopologyMetrics) Graph {
  for _, podMetrics := range metrics {
    node := Node{
      ID:    podID,
      Label: "Node " + podID,
      Metadata: map[string]interface{}{
        "timestamp":    time.Now(),
        "cpu_usage":    metrics.CPUUsage,
        "mem_usage":    metrics.MemUsage,
        "queue_size":   metrics.QueueSize,
        "num_rejected": metrics.NumRejected,
      },
    }
  }
	edge := Edge{
		Source: podID,
		Target: fmt.Sprintf("%d", time.Now().UnixNano()),
	}
	graph := Graph{
		Nodes: []Node{node},
		Edges: []Edge{edge},
	}
}

func StoreMetrics(metrics []*ktprom.TopologyMetrics) error {
	collection := client.Database("metricsdb").Collection("graphs")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := collection.InsertOne(ctx, graph)
	if err != nil {
		return err
	}
	log.Println("Gráfico guardado com sucesso no banco de dados.")
	return nil
}
