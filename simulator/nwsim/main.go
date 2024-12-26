package main
import (
  "os"
  "fmt"
  "math/rand"
)

const min, max = 1, 5
const minTransform, maxTransform = 0.5, 2.0

type SimService struct {
  targets map[*SimService]float64
  sources map[*SimService]struct{}
  waitTime float64
  id int
}

func newSimService(id int, wts ...float64) *SimService {
  var wt float64
  if len(wts) == 0 {
    wt = rand.Float64()
  } else {
    wt = wts[0]
  }
  r := min + wt * (max-min)
  s := SimService{
    id: id,
    waitTime: r,
    targets: make(map[*SimService]float64),
    sources: make(map[*SimService]struct{}),
  }
  return &s
}

func connect(head, foot *SimService) {
  head.targets[foot] = minTransform + rand.Float64()*(maxTransform - minTransform)
  foot.sources[head] = struct{}{}
}

func generateManifests(svcs []*SimService) {
  for _, svc := range svcs {
    fmt.Printf("Generating manifest for svc id %v\n", (*svc).id)
  }
}

func split(svc *SimService, net []*SimService, parentProb, netProb float64) {
  newSvc := newSimService(svc.id+1) 
  for t := range svc.targets {
    if rand.Float64() < parentProb {
      newSvc.targets[t] = minTransform + rand.Float64()*(maxTransform - minTransform)
    }
  }
  for s := range svc.sources {
    if rand.Float64() < parentProb {
      newSvc.sources[s] = struct{}{}
    }
  }
  for _, s := range net {
    if rand.Float64() < netProb {
      _, ok := newSvc.targets[s]
      if !ok {
        newSvc.targets[s] = minTransform + rand.Float64()*(maxTransform - minTransform)
      }
    }
    if rand.Float64() < netProb {
      _, ok := newSvc.sources[s]
      if !ok {
        newSvc.sources[s] = struct{}{}
      }
    }
  }
}

func randomRecursiveGenerate(seed *SimService) []*SimService {
  for
}

func main() {
  svcs := make([]*SimService, 0)
  // generate source
  src := newSimService(0, 0)
  // generate middle
  middle := newSimService(2)
  connect(src, middle)
  // generate sink
  sink := newSimService(1, 0)
  connect(middle, sink)
  // recursively expand middle
  net := randomRecursiveGenerate(middle)
  svcs = append(svcs, src)
  svcs = append(svcs, net...)
  svcs = append(svcs, sink)
  generateManifests(svcs)
}
