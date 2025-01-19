package ktmodel

import (
  // "time"
  "fmt"
  "strings"
  "strconv"
  //  "sync"

  // "github.com/shirou/gopsutil/v4/cpu"
  // "github.com/shirou/gopsutil/v4/mem"
)

type TopologyMetrics struct {
  // CPUUsage    float64         `json:"cpu_usage"`
  // MemUsage    float64         `json:"mem_usage"`
  // QueueSize   int             `json:"queue_size"`
  // QueueUse    int             `json:"queue_use"`
  // NumRejected int             `json:"num_rejected"`
  SentPkgs  map[string]int      `json:"sent_pkgs"`
  Metrics   map[string]float64  `json:"metrics"`
  // mx          sync.Mutex
}
func NewTopologyMetrics() *TopologyMetrics {
  return &TopologyMetrics{
    SentPkgs: make(map[string]int),
    Metrics: make(map[string]float64),
  }
}

// func (t *TopologyMetrics) UpdateCPU() {
//   cpuUsage, _ := cpu.Percent(time.Duration(0), false)
//   t.CPUUsage = cpuUsage[0]
// }

// func (t *TopologyMetrics) UpdateMem() {
//   mem_stat, _ := mem.VirtualMemory()
//   t.MemUsage = mem_stat.UsedPercent
// }

// func (t *TopologyMetrics) SetQueueSize(qs int) {
//   t.mx.Lock()
//   t.QueueSize = qs
//   t.mx.Unlock()
// }

// func (t *TopologyMetrics) SetQueueUse(qu int) {
//   t.mx.Lock()
//   t.QueueUse = qu
//   t.mx.Unlock()
// }

// func (t *TopologyMetrics) IncNumRejected() {
//   t.mx.Lock()
//   t.NumRejected += 1
//   t.mx.Unlock()
// }

// func (t *TopologyMetrics) ResetNumRejected() {
//   t.mx.Lock()
//   t.NumRejected = 0
//   t.mx.Unlock()
// }

// func (t *TopologyMetrics) IncSentPkgs(addr string) {
//   t.mx.Lock()
//   sentPkgs, ok := t.SentPkgs[addr]
//   if !ok {
//     t.SentPkgs[addr] = 1
//   } else {
//     t.SentPkgs[addr] = sentPkgs+1
//   }
//   t.mx.Unlock()
// }

// func (t *TopologyMetrics) ToPromStr() string {
//   str := fmt.Sprintf(
// `# HELP CPU usage
// # TYPE cpu_usage gauge
// cpu_usage %.3f
// # HELP Memory usage in Megabytes
// # TYPE mem_usage gauge
// mem_usage %.3f
// # HELP number of rejected messages due to full queue
// # TYPE num_rejected counter
// num_rejected %v
// # HELP Queue size
// # TYPE queue_size gauge
// queue_size %v
// # HELP Queue use
// # TYPE queue_use gauge
// queue_use %v`,
//   t.CPUUsage, t.MemUsage, t.NumRejected, t.QueueSize, t.QueueUse)
//   for addr, count := range t.SentPkgs {
//     str += fmt.Sprintf(
// `
// # HELP Number of successfully forwarded messages to %v
// # TYPE forwarded_count counter
// forwarded_count{addr=%v} %v`,
//       addr, addr, count)
//   }
//   return str
// }

func popAddrValue(key string) (string, string) {
  var (
    begin int
    end int
  )
  for i, v := range key {
    if string(v) == "{" {
      begin = i
    } else if string(v) == "}" {
      end = i
    }
  }
  prefix := key[:begin]
  label := key[begin+1:end]
  addr := strings.Split(label,"=")[1]
  return prefix, addr
}

func FromPromStr(s string) *TopologyMetrics {
  lines := strings.Split(s, "\n")
  metrics := NewTopologyMetrics()
  var addr string
  for _, line := range lines {
    if string(line[0]) == "#" {
      continue
    }
    splitLine := strings.Split(line, " ")
    prefix := splitLine[0]
    if strings.Contains(prefix, "{") {
      prefix, addr = popAddrValue(prefix)
    }
    valueStr := strings.TrimSpace(splitLine[1])
    // switch prefix {
    // case "cpu_usage": 
    //   metrics.CPUUsage, _ = strconv.ParseFloat(valueStr, 64)
    // case "mem_usage":
    //   metrics.MemUsage, _ = strconv.ParseFloat(valueStr, 64)
    // case "queue_size":
    //   metrics.QueueSize, _ = strconv.Atoi(valueStr)
    // case "queue_use":
    //   metrics.QueueUse, _ = strconv.Atoi(valueStr)
    // case "num_rejected":
    //   metrics.NumRejected, _ = strconv.Atoi(valueStr)
    if prefix == "forwarded_count" {
      metrics.SentPkgs[addr], _ = strconv.Atoi(valueStr)
    } else if addr != "" {
      metrics.Metrics[fmt.Sprintf("%v-%v", prefix, addr)], _ = 
        strconv.ParseFloat(valueStr, 64)
    } else {
      metrics.Metrics[prefix], _ = 
        strconv.ParseFloat(valueStr, 64)
    }
  }
  return metrics
}
