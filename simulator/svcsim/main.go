package main

import (
  "net/http"
  "os"
  "fmt"
  "log"
  "strconv"
  "strings"

  "github.com/prometheus/client_golang/prometheus"
)

type Metrics struct {
  cpu prometheus.Gauge
  mem prometheus.Gauge
  rejected prometheus.Counter
}

type Packet struct {
  id int
  value float64
}

var globalServiceInfo struct {
  id int
  queue []Packet
  queueSize int
  transforms map[string]float64
  delays map[string]int
}

func configureGlobalServiceInfo() {
  globalServiceInfo.queue = make([]Packet, 0)
  globalServiceInfo.transforms = make(map[string]float64)
  globalServiceInfo.delays = make(map[string]int)

  globalServiceInfo.id, _   = strconv.Atoi(os.Getenv("ID"))
  globalServiceInfo.queueSize, _   = strconv.Atoi(os.Getenv("QUEUESIZE"))

  targets     := strings.Split(os.Getenv("TARGETS"), ",")
  delays      := strings.Split(os.Getenv("DELAYS"), ",")
  transforms  := strings.Split(os.Getenv("TRANSFORMS"), ",")
  for i, target := range targets {
    delay, err := strconv.Atoi(delays[i])
    if err != nil {
      log.Fatalf("Bad delay formatting: %v\n", delays[i])
    }
    globalServiceInfo.delays[target] = delay

    transform, err := strconv.ParseFloat(transforms[i], 64)
    if err != nil {
      log.Fatalf("Bad delay formatting: %v\n", transforms[i])
    }
    globalServiceInfo.transforms[target] = transform
  }
}

func forward(w http.ResponseWriter, req *http.Request) {

}

func metrics(w http.ResponseWriter, req *http.Request) {
  fmt.Println(globalServiceInfo)
}


func main() {
  configureGlobalServiceInfo()
  http.HandleFunc("/", forward)
  http.HandleFunc("/metrics", metrics)
  http.ListenAndServe(":8080", nil)
}
