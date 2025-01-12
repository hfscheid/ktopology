package main

import (
	"context"
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

func buildAddrMap(metrics []IdMetrics) map[string]string {
  addrMap := make(map[string]string)
  for _, idMetrics := range metrics {
    addrMap[idMetrics.addr] = idMetrics.id
  }
  return addrMap
}

func buildGraphFromMetrics(metrics []IdMetrics) *Graph {
  addrMap := buildAddrMap(metrics)
  nodes := make([]Node, 0, len(metrics))
  edges := make([]Edge, 0, len(metrics))
  for _, podMetrics := range metrics {
    node := Node{
      ID:    podMetrics.id,
      Metadata: map[string]interface{}{
        "timestamp":    time.Now(),
        "cpu_usage":    podMetrics.metrics.CPUUsage,
        "mem_usage":    podMetrics.metrics.MemUsage,
        "queue_size":   podMetrics.metrics.QueueSize,
        "num_rejected": podMetrics.metrics.NumRejected,
      },
    }
    nodes = append(nodes, node)
    for addr, _ := range podMetrics.metrics.SentPkgs {
      edge := Edge{
        Source: podMetrics.id,
        Target: addrMap[addr],
      }
      edges = append(edges, edge)
    }
  }
	graph := &Graph{
		Nodes: nodes,
		Edges: edges,
	}
  return graph
}

func StoreMetrics(metrics []IdMetrics) error {
	collection := client.Database("metricsdb").Collection("graphs")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
  graph := buildGraphFromMetrics(metrics)
	_, err := collection.InsertOne(ctx, graph)
	if err != nil {
		return err
	}
	log.Println("Successfully stored graph in database")
	return nil
}
