package storage

import (
  "encoding/json"
  "log"
  "fmt"
  "os"

  "github.com/hfscheid/ktopology/ktmodel"
)

var logger = log.New(os.Stdout, "[storage] ", log.Ltime)

func processGraph(tData []ktmodel.TopologyData) ([]byte, error) {
  graph := ktmodel.BuildNwTopologyFromData(tData)
  qos := func(tData []ktmodel.TopologyData) (string, float64) {
    totalRejected := 0.0
    for _, tDatum := range tData {
      totalRejected += tDatum.Metrics.Metrics["num_rejected"]
    }
    return "total_rejected", float64(totalRejected)
  }
  graph.AddNwCalculus(qos)
  graph.DoNwCalculus()
  graphData, err := json.Marshal(*graph)
  return graphData, err
}

func StoreMetrics(tData []ktmodel.TopologyData) error {
  graphData, err  := processGraph(tData)
  if err != nil {
    return fmt.Errorf("Could not marshal graphData: %v", err)
  }
  logger.Print(string(graphData))
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
