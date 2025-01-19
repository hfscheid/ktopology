package main
import (
  "os"
  "log"
  "strconv"
  "strings"
  "sync"

  // "svcsim/synctypes"
)
var globalServiceInfo struct {
  id          int
  queue       []Packet
  queueSize   int
  transforms  map[string]float64
  delays      map[string]int
  numRejected int
  sentMsgs    map[string]int
  queuemx     sync.Mutex
  nrejmx      sync.Mutex
  sentmx      sync.Mutex
  // reg         *prometheus.Registry
}

func configureGlobalServiceInfo() {
  globalServiceInfo.queue = make([]Packet, 0)
  globalServiceInfo.transforms = make(map[string]float64)
  globalServiceInfo.sentMsgs = make(map[string]int)
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
  
  // globalMetrics.cpu = prometheus.NewGauge(prometheus.GaugeOpts{
  //   Name: "cpu_use",
  //   Help: "Percentage of CPU power used by process",
  // })
  // globalMetrics.mem = prometheus.NewGauge(prometheus.GaugeOpts{
  //   Name: "memory_use",
  //   Help: "Bytes of memory used by process",
  // })
  // globalMetrics.rejected = prometheus.NewGauge(prometheus.GaugeOpts{
  //   Name: "num_rejected",
  //   Help: "Number of rejected messages",
  // })

  // globalServiceInfo.reg = prometheus.NewRegistry()
  // globalServiceInfo.reg.MustRegister(globalMetrics.cpu)
  // globalServiceInfo.reg.MustRegister(globalMetrics.mem)
  // globalServiceInfo.reg.MustRegister(globalMetrics.rejected)
}

