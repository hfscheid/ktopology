package ktprom

import (
  "time"
  "fmt"
  "strings"
  "strconv"

  "github.com/shirou/gopsutil/v4/cpu"
  "github.com/shirou/gopsutil/v4/mem"
)

type TopologyMetrics struct {
  CPUUsage    float64         `json:"cpu_usage"`
  MemUsage    float64         `json:"mem_usage"`
  QueueSize   int             `json:"queue_size"`
  NumRejected int             `json:"num_rejected"`
  SentPkgs    map[string]int  `json:"sent_pkgs"`
}
func NewTopologyMetrics() *TopologyMetrics {
  return &TopologyMetrics{
    SentPkgs: make(map[string]int),
  }
}

func (t *TopologyMetrics) UpdateCPU() {
  cpuUsage, _ := cpu.Percent(time.Duration(0), false)
  t.CPUUsage = cpuUsage[0]
}

func (t *TopologyMetrics) UpdateMem() {
  mem_stat, _ := mem.VirtualMemory()
  t.MemUsage = mem_stat.UsedPercent
}

func (t *TopologyMetrics) SetQueueSize(qs int) {
  t.QueueSize = qs
}

func (t *TopologyMetrics) IncNumRejected() {
  t.NumRejected += 1
}

func (t *TopologyMetrics) IncSentPkgs(addr string) {
  sentPkgs, ok := t.SentPkgs[addr]
  if !ok {
    t.SentPkgs[addr] = 1
  } else {
    t.SentPkgs[addr] = sentPkgs+1
  }
}

func (t *TopologyMetrics) ToPromStr() string {
  str := fmt.Sprintf(
`# HELP CPU usage
# TYPE cpu_usage gauge
cpu_usage %.3f
# HELP Memory usage in Megabytes
# TYPE mem_usage gauge
mem_usage %.3f
# HELP number of rejected messages due to full queue
# TYPE num_rejected counter
num_rejected %v
# HELP Queue size
# TYPE queue_size gauge
queue_size %v`,
  t.CPUUsage, t.MemUsage, t.NumRejected, t.QueueSize)
  for addr, count := range t.SentPkgs {
    str += fmt.Sprintf(
`
# HELP Number of successfully forwarded messages to %v
# TYPE forwarded_count counter
forwarded_count{addr=%v} %v`,
      addr, addr, count)
  }
  return str
}

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
  return key[:begin], key[begin:end]
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
    switch prefix {
    case "cpu_usage": 
      metrics.CPUUsage, _ = strconv.ParseFloat(valueStr, 64)
    case "mem_usage":
      metrics.MemUsage, _ = strconv.ParseFloat(valueStr, 64)
    case "queue_size":
      metrics.QueueSize, _ = strconv.Atoi(valueStr)
    case "num_rejected":
      metrics.NumRejected, _ = strconv.Atoi(valueStr)
    case "forwarded_count":
      metrics.SentPkgs[addr], _ = strconv.Atoi(valueStr)
    }
  }
  return metrics
}
