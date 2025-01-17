package main

import (
  "io"
  "os"
  "log"
  "fmt"
  "net/http"
  "strings"
  "math/rand"
  "math"
  "strconv"
  "time"
)

func logBody(w http.ResponseWriter, req *http.Request) {
  data, _ := io.ReadAll(req.Body)
  log.Printf("received %v\n", string(data))
  w.Write([]byte("OK"))
}

func genRequests(target string, interval float64) {
  for counter := 0;; counter++ {
    time.Sleep(time.Duration(math.Pow(10.0,9.0)*interval))
    pkg := fmt.Sprintf("%v,%.2f", counter, rand.Float64()+1)
    log.Printf("sending %v\n", pkg)
      _, err := http.Post(
        target,
        "text/plain",
        strings.NewReader(pkg),
      )
      if err != nil {
        log.Print(err)
      }
  }
}

func main() {
  http.HandleFunc("/", logBody)
  interval, _ := strconv.ParseFloat(os.Getenv("INTERVAL"), 64)
  go genRequests(os.Getenv("SOURCE"), interval)
  http.ListenAndServe(":80", nil)
}
