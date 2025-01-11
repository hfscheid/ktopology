package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

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

func StoreMetrics(metrics *Metrics, nodeID string) error {
	collection := client.Database("metricsdb").Collection("graphs")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	node := Node{
		ID:    nodeID,
		Label: "Node " + nodeID,
		Metadata: map[string]interface{}{
			"timestamp":    time.Now(),
			"cpu_usage":    metrics.CPUUsage,
			"mem_usage":    metrics.MemUsage,
			"queue_size":   metrics.QueueSize,
			"num_rejected": metrics.NumRejected,
		},
	}

	edge := Edge{
		Source: nodeID,
		Target: fmt.Sprintf("%d", time.Now().UnixNano()),
	}

	graph := Graph{
		Nodes: []Node{node},
		Edges: []Edge{edge},
	}

	fmt.Printf("Guardando gráfico:\n")
	fmt.Printf("  Node ID: %s\n", node.ID)
	fmt.Printf("  CPU Usage: %.2f\n", node.Metadata["cpu_usage"].(float64))
	fmt.Printf("  Memory Usage: %.2f MB\n", node.Metadata["mem_usage"].(float64))
	fmt.Printf("  Queue Size: %d\n", node.Metadata["queue_size"].(int))
	fmt.Printf("  Number of Rejected Messages: %d\n", node.Metadata["num_rejected"].(int))

	_, err := collection.InsertOne(ctx, graph)
	if err != nil {
		return err
	}

	fmt.Println("Gráfico guardado com sucesso no banco de dados.")

	return nil
}
