package storage

import (
	// "context"
  "encoding/json"
	"log"
  "fmt"
	"os"
	"time"

  "collector/pkg/ktprom"

	// "go.mongodb.org/mongo-driver/mongo"
	// "go.mongodb.org/mongo-driver/mongo/options"
)

// var client *mongo.Client
var logger = log.New(os.Stdout, "[storage] ", log.Ltime)

type IdMetrics struct {
  Id      string
  Addr    string
  Service string
  Host    string
  Metrics *ktprom.TopologyMetrics
}


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

// func init() {
// 	var err error
// 
// 	mongoURI := os.Getenv("MONGO_URI")
// 	if mongoURI == "" {
// 		mongoURI = "mongodb://localhost:27017"
// 	}
// 
// 	clientOptions := options.Client().ApplyURI(mongoURI)
// 	client, err = mongo.Connect(context.TODO(), clientOptions)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 
// 	err = client.Ping(context.TODO(), nil)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// }

func buildAddrMap(metrics []IdMetrics) map[string]string {
  addrMap := make(map[string]string)
  for _, idMetrics := range metrics {
    addrMap[idMetrics.Service] = idMetrics.Id
  }
  return addrMap
}

func buildGraphFromMetrics(metrics []IdMetrics) *Graph {
  addrMap := buildAddrMap(metrics)
  nodes := make([]Node, 0, len(metrics))
  edges := make([]Edge, 0, len(metrics))
  for _, podMetrics := range metrics {
    node := Node{
      ID:    podMetrics.Id,
      Metadata: map[string]interface{}{
        "timestamp":    time.Now(),
        "cpu_usage":    podMetrics.Metrics.CPUUsage,
        "mem_usage":    podMetrics.Metrics.MemUsage,
        "queue_size":   podMetrics.Metrics.QueueSize,
        "num_rejected": podMetrics.Metrics.NumRejected,
      },
    }
    nodes = append(nodes, node)
    for addr, _ := range podMetrics.Metrics.SentPkgs {
      edge := Edge{
        Source: podMetrics.Id,
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
	// collection := client.Database("metricsdb").Collection("graphs")
	// ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// defer cancel()
  graph := buildGraphFromMetrics(metrics)
  graphData, err := json.Marshal(*graph)
  if err != nil {
    return fmt.Errorf("Could not marshal graphData: %v", err)
  }
  logger.Print(string(graphData))
	// _, err := collection.InsertOne(ctx, graph)
	// if err != nil {
	// 	return err
	// }
	logger.Println("Successfully stored graph in database")
  f, err := os.OpenFile("./data.jsonl", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
  if err != nil {
    return fmt.Errorf("Could not open data file: %v", err)
  }
  defer f.Close()
  if _, err := f.Write(graphData);
  err != nil {
    return fmt.Errorf("Could not write data file: %v", err)
  }
	return nil
}
