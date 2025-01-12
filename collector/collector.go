package main

import (
  "io/ioutil"
  "net/http"

  "collector/ktprom"
)

func CollectMetrics(url string) (*ktprom.TopologyMetrics, error) {
  resp, err := http.Get(url)
  if err != nil {
    return nil, err
  }
  defer resp.Body.Close()
  body, err := ioutil.ReadAll(resp.Body)
  if err != nil {
    return nil, err
  }
  metricsText := string(body)
  metrics := ktprom.FromPromStr(metricsText)
  return metrics, nil
}
