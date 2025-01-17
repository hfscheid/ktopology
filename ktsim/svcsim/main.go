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

  "github.com/hfscheid/ktopology/ktmodel"
)

type Packet struct {
  id int
  value float64
}

var globalMetrics *ktmodel.TopologyMetrics
var client = &http.Client{
  Transport: &http.Transport{DisableKeepAlives: true},
}

func exportMetrics(w http.ResponseWriter, req *http.Request) {
  globalMetrics.UpdateCPU()
  globalMetrics.UpdateMem()
  globalServiceInfo.queue.Lock()
  qSize := globalServiceInfo.queue.Size()
  globalServiceInfo.queue.Unlock()
  globalMetrics.SetQueueSize(qSize)
  w.Write([]byte(globalMetrics.ToPromStr()))
}

func pushToQueue(w http.ResponseWriter, req *http.Request) {
  globalServiceInfo.queue.Lock()
  if globalServiceInfo.queue.Size() == globalServiceInfo.queueSize {
    w.Write([]byte("REJECTED"))
    globalMetrics.IncNumRejected()
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
      globalMetrics.IncSentPkgs(target)
      res, err := client.Post(
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
  globalMetrics = ktmodel.NewTopologyMetrics()
  http.HandleFunc("/", pushToQueue)
  http.HandleFunc("/metrics", exportMetrics)
  go forward()
  http.ListenAndServe(":80", nil)
}
