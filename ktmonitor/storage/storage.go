package storage

import (
  // "context"
  "encoding/json"
  "log"
  "fmt"
  "os"
  "time"

  "github.com/hfscheid/ktopology/ktprom"
  "github.com/hfscheid/ktopology/ktgraph"
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
  Deployment string
  Metrics *ktprom.TopologyMetrics
}

// type Node struct {
//  ID          string                 `json:"id"`
//   Addr        string                 `json:"podIp"`
//   Host        string                 `json:"host"`
//   Service     string                 `json:"service"`
//   Deployment  string                 `json:"deployment"`
//  Metadata    map[string]interface{} `json:"metadata"`
// }
// 
// type Edge struct {
//  Source string `json:"source"`
//  Target string `json:"target"`
// }
// 
// type Graph struct {
//  Nodes []Node `json:"nodes"`
//  Edges []Edge `json:"edges"`
// }

// func init() {
//  var err error
// 
//  mongoURI := os.Getenv("MONGO_URI")
//  if mongoURI == "" {
//    mongoURI = "mongodb://localhost:27017"
//  }
// 
//  clientOptions := options.Client().ApplyURI(mongoURI)
//  client, err = mongo.Connect(context.TODO(), clientOptions)
//  if err != nil {
//    log.Fatal(err)
//  }
// 
//  err = client.Ping(context.TODO(), nil)
//  if err != nil {
//    log.Fatal(err)
//  }
// }

func buildAddrMap(metrics []IdMetrics) map[string][]string {
  addrMap := make(map[string][]string)
  for _, idMetrics := range metrics {
    if  _, ok := addrMap[idMetrics.Service];
        !ok {
          addrMap[idMetrics.Service] = make([]string, 0)
    }
    addrMap[idMetrics.Service] = append(addrMap[idMetrics.Service], idMetrics.Id)
  }
  return addrMap
}

func buildGraphFromMetrics(metrics []IdMetrics) *ktgraph.Graph {
  addrMap := buildAddrMap(metrics)
  nodes := make([]ktgraph.Node, 0, len(metrics))
  edges := make([]ktgraph.Edge, 0, len(metrics))
  for _, podMetrics := range metrics {
    node := ktgraph.Node{
      ID:         podMetrics.Id,
      Host:       podMetrics.Host,
      Service:    podMetrics.Service,
      Deployment: podMetrics.Deployment,
      Addr:       podMetrics.Addr,
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
      for _, target := range addrMap[addr] {
        edge := ktgraph.Edge{
          Source: podMetrics.Id,
          Target: target,
        }
        edges = append(edges, edge)
      }
    }
  }
  graph := &ktgraph.Graph{
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
  //  return err
  // }
  f, err := os.OpenFile("./data.jsonl", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
  if err != nil {
    return fmt.Errorf("Could not open data file: %v", err)
  }
  defer f.Close()
  graphData = append(graphData, '\n')
  if _, err := f.Write(graphData);
  err != nil {
    return fmt.Errorf("Could not write data file: %v", err)
  }
  logger.Println("Successfully stored graph in database")
  return nil
}
