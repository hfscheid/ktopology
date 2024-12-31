package main

import (
  "fmt"
  "math/rand"
  "os"
  "strconv"

  "nwsim/randomflux"
  "nwsim/services"
)

func servicesFromGraph(g *randomflux.Graph) map[int]*services.Service {
  sMap := make(map[int]*services.Service)
  for _, id := range g.Ids() {
    s := services.New(id, rand.Intn(10)+5)
    sMap[id] = s
  }
  for _, id := range g.Ids() {
    nbs, _ := g.Neighbours(id)
    for _, nb := range nbs {
      fmt.Printf("Adding target %v to service %v\n", nb, id)
      sMap[id].AddTarget(
        sMap[nb],
        rand.Intn(5)+1,
        rand.Float64()+0.5,
      )
    }
  }
  return sMap
}

func writeManifests(sMap map[int]*services.Service) {
  os.Mkdir("manifests", 0750)
  for id, s := range sMap {
    f, _ := os.Create(fmt.Sprintf("manifests/manifest-%v.yaml", id))
    f.Write(s.ConfigMap())
    f.Write([]byte("\n---\n"))
    f.Write(s.Deployment())
    f.Write([]byte("\n---\n"))
    f.Write(s.Service())
    f.Write([]byte("\n---\n"))
  }
}

func main() {
  size, _ := strconv.Atoi(os.Args[1])
  g := randomflux.New(size)
  g.Draw()
  writeManifests(servicesFromGraph(g))
}
