package main
import (
  "fmt"
  "math/rand"
)

type SimService struct {
  targets []*SimService  
  transforms []int
  waitTime int
  id int
}

func newSimService(id int, wts ...int) *SimService {
  var wt int
  if len(wts) == 0 {
    wt = rand.Intn(5)
  } else {
    wt = wts[0]
  }
  s := SimService{
    id: id,
    waitTime: wt,
  }
  return &s
}

func generateManifests(svcs []*SimService) {
  for _, svc := range svcs {
    fmt.Printf("Generating manifest for svc id %v\n", (*svc).id)
  }
}

func main() {
  svcs := make([]*SimService, 0)
  // generate source
  svcs = append(svcs, newSimService(0, 0))
  // generate middle
  // generate sink
  svcs = append(svcs, newSimService(1, 0))
  generateManifests(svcs)
}
