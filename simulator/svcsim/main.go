package main

import (
  "fmt"
  "io"
  "log"
  "net/http"
  "os"
  "strconv"
  "strings"
  "sync"
  "time"

  "github.com/prometheus/client_golang/prometheus"
)

type Packet struct {
  id int
  value float64
}

var globalMetrics struct {
  cpu prometheus.Gauge
  mem prometheus.Gauge
  rejected prometheus.Counter
}

var mutex sync.Mutex

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

  // TODO: sink svc has empty targets, must have condition

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

func pushToQueue(w http.ResponseWriter, req *http.Request) {
  mutex.Lock()
  if len(globalServiceInfo.queue) == globalServiceInfo.queueSize {
    w.Write([]byte("REJECTED"))
    mutex.Unlock()
    return
  }
  mutex.Unlock()

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
  mutex.Lock()
  globalServiceInfo.queue = append(
    globalServiceInfo.queue,
    Packet{pId, pVal},
  )
  mutex.Unlock()
  w.Write([]byte("OK"))
}

func metrics(w http.ResponseWriter, req *http.Request) {
  fmt.Println(globalServiceInfo)
}

func sendToTargets(p Packet) {
  for target, delay := range globalServiceInfo.delays {
    time.Sleep(10^9*time.Duration(delay))
    newValue := globalServiceInfo.transforms[target]*p.value
    for {
      res, err := http.Post(
        target,
        "text/plain",
        strings.NewReader(
          fmt.Sprintf("%v,%.2f", p.id, newValue),
        ),
      )
      resBody, _ := io.ReadAll(res.Body)
      if string(resBody) == "OK" {
        log.Printf("Message %v to %v OK\n", p.id, target)
        break
      } else if string(resBody) == "REJECTED" {
        log.Printf("Message %v to %v REJECTED\n", p.id, target)
        continue
      }
      if err != nil {
        log.Printf("Error sending message to target: %v\n", err)
      }
    }
  }
}

func forward() {
  for {
    mutex.Lock()
    if len(globalServiceInfo.queue) > 0 {
      p := globalServiceInfo.queue[0]
      globalServiceInfo.queue = globalServiceInfo.queue[1:]
      mutex.Unlock()
      sendToTargets(p)
    } else {
      mutex.Unlock()
    }
  }
}


func main() {
  configureGlobalServiceInfo()
  http.HandleFunc("/", pushToQueue)
  http.HandleFunc("/metrics", metrics)
  // go forward()
  http.ListenAndServe(":80", nil)
}
