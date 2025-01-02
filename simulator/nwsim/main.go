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
  ids := g.Ids()
  for _, id := range ids {
    s := services.New(id, rand.Intn(10)+5)
    sMap[id] = s
  }
  for _, id := range ids {
    nbs, _ := g.Neighbours(id)
    for _, nb := range nbs {
      sMap[id].AddTarget(
        sMap[nb],
        rand.Intn(5)+1,
        rand.Float64()+0.5,
        fmt.Sprintf("service-%v", nb),
      )
    }
  }

  // Tie simulated network to sink
  sink := services.New(-1, 0)
  sMap[ids[len(ids)-1]].AddTarget(
    sink,
    rand.Intn(5)+1,
    rand.Float64()+0.5,
    "sink",
  )
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
