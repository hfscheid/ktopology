package main
import (
  "net/http"
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

type Service struct {
  id int
  queue []Packet
  queueSize int
  transforms map[int]float64
  delays map[int]int
}

func forward(w http.ResponseWriter, req *http.Request) {
}

func metrics(w http.ResponseWriter, req *http.Request) {
}

func main() {
  http.HandleFunc("/", forward)
  http.HandleFunc("/metrics", forward)
}
