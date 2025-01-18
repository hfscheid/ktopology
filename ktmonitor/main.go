package main
import (
  "log"
  "os"
  "strconv"
  "time"
  "sync"
  "fmt"
  "io"
  "net/http"

  "github.com/hfscheid/ktopology/ktmonitor/poddiscovery"
  "github.com/hfscheid/ktopology/ktmonitor/storage"
  "github.com/hfscheid/ktopology/ktmodel"
)

var logger = log.New(os.Stdout, "[ktmonitor] ", log.Ltime)
var httpClient = http.Client {
  Timeout: 1 * time.Second,
}

func getPollInterval() time.Duration {
  var pollInterval time.Duration
  pollIntervalStr := os.Getenv("POLL_INTERVAL")
  if pollIntervalStr == "" {
    pollInterval = 30 * time.Second
  } else {
    pollIntervalInt, err := strconv.Atoi(pollIntervalStr)
    if err != nil {
      log.Fatalf("Could not convert POLL_INTERVAL to int: %v", err)
    }
    pollInterval = time.Duration(pollIntervalInt) * time.Second
  }
  log.Printf("Poll interval set to: %v\n", pollInterval)
  return pollInterval
}

func collectMetrics(url string) (*ktmodel.TopologyMetrics, error) {
  resp, err := httpClient.Get(url)
  if err != nil {
    return nil, err
  }
  defer resp.Body.Close()
  body, err := io.ReadAll(resp.Body)
  if err != nil {
    return nil, err
  }
  metricsText := string(body)
  metrics := ktmodel.FromPromStr(metricsText)
  return metrics, nil
}

func poll(pod poddiscovery.PodInfo,
          topDataChan chan ktmodel.TopologyData,
          wg *sync.WaitGroup) {
    defer wg.Done()
    logger.Printf("Polling metrics from pod %v...\n", pod.IP)
    podURL := fmt.Sprintf("http://%v/metrics", pod.IP)
    metrics, err := collectMetrics(podURL)
    if err != nil {
      logger.Printf("Could not read metrics from %s: %v", podURL, err)
      // wg.Done()
      return
    }
    logger.Printf("Sending metrics %v to channel...\n", metrics)
    topDataChan <- ktmodel.TopologyData {
      ID: pod.Name,
      Addr: pod.IP,
      Host: pod.HostIP,
      Service: pod.Service,
      Metrics: *metrics,
      Deployment: pod.Deployment,
    }
    logger.Println("Metrics sent to channel")
    // wg.Done()
}

func recurrentPoll(pollInterval time.Duration) {
  var wg sync.WaitGroup
  for {
    pods, err:= poddiscovery.ListPods()
    logger.Println(pods)
    numPods := len(pods)
    if err != nil {
      logger.Fatalf("Could not discover pods: %v", err)
    }
    wg.Add(numPods)
    topDataChan := make(chan ktmodel.TopologyData, numPods)
    for _, pod := range pods {
      go poll(pod, topDataChan, &wg)
    }
    wg.Wait()
    close(topDataChan)
    logger.Println("Received all pod metrics. Sending to storage...")
    topData := make([]ktmodel.TopologyData, 0, numPods)
    for pms, ok := <- topDataChan; ok; pms, ok = <-topDataChan {
      topData = append(topData, pms)
    }
    if err := storage.StoreMetrics(topData);
    err != nil {
      logger.Printf("Could not store topology: %v\n", err)
    }
    time.Sleep(pollInterval)
  }
}

func main() {
  pollInterval := getPollInterval()
  recurrentPoll(pollInterval)
}
