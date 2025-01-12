package main

import (
  "context"
  "fmt"
  "log"
  "os"
  "strconv"
  "time"
  "sync"

  "collector/ktprom"

  metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
  "k8s.io/client-go/kubernetes"
  "k8s.io/client-go/rest"
  // "k8s.io/client-go/tools/clientcmd"
)

type IdMetrics struct {
  id      string
  addr    string
  metrics *ktprom.TopologyMetrics
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

func getK8sClient() *kubernetes.Clientset {
  config, err := rest.InClusterConfig()
  if err != nil {
    log.Fatalf("Could not get in-cluster Kubernetes config: %v", err)
  }
  clientset, err := kubernetes.NewForConfig(config)
  if err != nil {
    log.Fatalf("Could not create Kubernetes clientset: %v", err)
  }
  return clientset
}

func poll(k8sClient *kubernetes.Clientset, pollInterval time.Duration) {
  var wg sync.WaitGroup
  for {
    log.Println("Polling Kubernetes API for Pod IPs...")
    pods, err := k8sClient.CoreV1().
      Pods("default").
      List(context.TODO(), metav1.ListOptions{})
    if err != nil {
      log.Printf("Could not list nodes: %v", err)
      time.Sleep(pollInterval)
      continue
    }
    numPods := pods.Size()
    log.Printf("Found %v pods in namespace \"default\"\n", numPods)
    log.Println("Polling metrics from pods...")
    wg.Add(numPods)
    podMetricsChan := make(chan IdMetrics, numPods)
    for _, pod := range pods.Items {
      podIP := pod.Status.PodIP
      podURL := fmt.Sprintf("http://%v:8080/metrics", podIP)
      go func() {
        metrics, err := CollectMetrics(podURL)
        if err != nil {
          log.Printf("Could not read metrics from %s: %v", podURL, err)
          return
        }
        podMetricsChan <- IdMetrics{id: pod.Name, addr: podIP, metrics: metrics}
        wg.Done()
      }()
    }
    wg.Wait()
    podMetrics := make([]IdMetrics, 0, numPods)
    for pms, ok := <- podMetricsChan; ok; pms, ok = <-podMetricsChan {
      podMetrics = append(podMetrics, pms)
    }
    if err := StoreMetrics(podMetrics);
    err != nil {
      log.Printf("Could not store topology: %v\n", err)
    }
    time.Sleep(pollInterval)
  }
}

func main() {
  pollInterval := getPollInterval()
  k8sClient := getK8sClient()
  poll(k8sClient, pollInterval)
}
