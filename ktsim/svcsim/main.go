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

  // "github.com/hfscheid/ktopology/ktmodel"
)

type Packet struct {
  id int
  value float64
}

var client = &http.Client{
  Transport: &http.Transport{DisableKeepAlives: true},
}

func exportMetrics(w http.ResponseWriter, req *http.Request) {
  globalServiceInfo.queuemx.Lock()
  qUse := len(globalServiceInfo.queue)
  globalServiceInfo.queuemx.Unlock()
  promStr := fmt.Sprintf(
`# HELP number of rejected messages due to full queue
# TYPE num_rejected counter
num_rejected %v
# HELP Queue size
# TYPE queue_size gauge
queue_size %v
# HELP Queue use
# TYPE queue_use gauge
queue_use %v`,
    globalServiceInfo.numRejected,
    globalServiceInfo.queueSize,
    qUse,
  )
  for k, v := range globalServiceInfo.sentMsgs {
    promStr += fmt.Sprintf(`
# HELP Number of successfully forwarded messages to %v
# TYPE forwarded_count counter
forwarded_count{addr=%v} %v`,
      k, k, v,
    )
  }
  w.Write([]byte(promStr))
  globalServiceInfo.nrejmx.Lock()
  globalServiceInfo.numRejected = 0
  globalServiceInfo.nrejmx.Unlock()
}

func pushToQueue(w http.ResponseWriter, req *http.Request) {
  globalServiceInfo.queuemx.Lock()
  if len(globalServiceInfo.queue) == globalServiceInfo.queueSize {
    globalServiceInfo.queuemx.Unlock()
    globalServiceInfo.nrejmx.Lock()
    globalServiceInfo.numRejected += 1
    globalServiceInfo.nrejmx.Unlock()
    w.Write([]byte("REJECTED"))
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

  globalServiceInfo.queue = append(
    globalServiceInfo.queue, 
    Packet{pId, pVal},
  )
  globalServiceInfo.queuemx.Unlock()
  w.Write([]byte("OK"))
}

func sendToTargets(p Packet) {
  for target, delay := range globalServiceInfo.delays {
    newValue := globalServiceInfo.transforms[target]*p.value
    for {
      time.Sleep(time.Duration(math.Pow(10.0, 9.0)*float64(delay)))
      log.Printf("Sending to %v\n", target)
      globalServiceInfo.sentmx.Lock()
      globalServiceInfo.sentMsgs[target] = 1
      globalServiceInfo.sentmx.Unlock()
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
    globalServiceInfo.queuemx.Lock()
    if len(globalServiceInfo.queue) > 0 {
      p := globalServiceInfo.queue[0]
      globalServiceInfo.queue = globalServiceInfo.queue[1:]
      globalServiceInfo.queuemx.Unlock()
      sendToTargets(p)
    } else {
      globalServiceInfo.queuemx.Unlock()
    }
  }
}


func main() {
  configureGlobalServiceInfo()
  // globalMetrics = ktmodel.NewTopologyMetrics()
  http.HandleFunc("/", pushToQueue)
  http.HandleFunc("/metrics", exportMetrics)
  go forward()
  http.ListenAndServe(":80", nil)
}
