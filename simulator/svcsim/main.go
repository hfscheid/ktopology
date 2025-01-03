package main

import (
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

  "svcsim/synctypes"
  
  "github.com/shirou/gopsutil/v4/cpu"
  "github.com/shirou/gopsutil/v4/mem"
)

type Packet struct {
  id int
  value float64
}

var globalMetrics struct {
  rejected synctypes.Int
}

func exportMetrics(w http.ResponseWriter, req *http.Request) {
  cpu_usage, _ := cpu.Percent(time.Duration(0), false)
  mem_stat, _ := mem.VirtualMemory()
  mem_usage := mem_stat.UsedPercent
  w.Write([]byte(fmt.Sprintf(
`# HELP CPU usage
# TYPE cpu_usage gauge
cpu_usage %.3f
# HELP Memory usage in Megabytes
# TYPE mem_usage gauge
mem_usage %.3f
# HELP number of rejected messages due to full queue
# TYPE num_rejected counter
num_rejected %v`,
  cpu_usage, mem_usage, globalMetrics.rejected.Get())))
}

func pushToQueue(w http.ResponseWriter, req *http.Request) {
  globalServiceInfo.queue.Lock()
  if globalServiceInfo.queue.Size() == globalServiceInfo.queueSize {
    w.Write([]byte("REJECTED"))
    globalMetrics.rejected.Sum(1)
    globalServiceInfo.queue.Unlock()
    return
  }

  rawBody, err := io.ReadAll(req.Body)
  if err != nil {
    log.Printf("Bad request: %v", err)
  }
  packetStr := strings.Split(string(rawBody), ",")
  pId, err  := strconv.Atoi(packetStr[0])
  if err != nil {
    log.Printf("Bad request id: %v\nError: %v\n", packetStr, err)
    return
  }
  pVal, err := strconv.ParseFloat(packetStr[1], 64)
  if err != nil {
    log.Printf("Bad request value: %v\nError: %v\n", packetStr, err)
    return
  }

  globalServiceInfo.queue.Push(
    Packet{pId, pVal},
  )
  globalServiceInfo.queue.Unlock()
  w.Write([]byte("OK"))
}

func sendToTargets(p Packet) {
  for target, delay := range globalServiceInfo.delays {
    newValue := globalServiceInfo.transforms[target]*p.value
    for {
      time.Sleep(time.Duration(math.Pow(10.0, 9.0)*float64(delay)))
      log.Printf("Sending to %v\n", target)
      res, err := http.Post(
        target,
        "text/plain",
        strings.NewReader(
          fmt.Sprintf("%v,%.2f", p.id, newValue),
        ),
      )
      if res == nil {
        log.Printf("No response from target.\n")
        break
      }
      resBody, err := io.ReadAll(res.Body)
      if err != nil {
        log.Printf("Response error when forwarding: %v\n", err)
        break
      }
      res.Body.Close()
      if string(resBody) == "OK" {
        log.Printf("Message %v to %v OK\n", p.id, target)
        break
      } else if string(resBody) == "REJECTED" {
        log.Printf("Message %v to %v REJECTED\n", p.id, target)
        continue
      }
    }
  }
}

func forward() {
  for {
    globalServiceInfo.queue.Lock()
    if globalServiceInfo.queue.Size() > 0 {
      p := globalServiceInfo.queue.Pop()
      globalServiceInfo.queue.Unlock()
      sendToTargets(p)
    } else {
      globalServiceInfo.queue.Unlock()
    }
  }
}


func main() {
  configureGlobalServiceInfo()
  http.HandleFunc("/", pushToQueue)
  http.HandleFunc("/metrics", exportMetrics)
  go forward()
  http.ListenAndServe(":8081", nil)
}
