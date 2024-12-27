package randomflux

import (
  "math/rand"
  "os"
  "fmt"

  "github.com/dominikbraun/graph"
  "github.com/dominikbraun/graph/draw"
)

type Graph struct {
  am [][]int
  size int
}

func (g *Graph) mitosis() {
  id := g.size-2

  for i := range g.am {
    if g.am[i][id] == 1 {
      r := rand.Float64()
      if r > 0.5 {
        // connect new vertex with its parent
        g.am[id][id+1] = 1
        g.am[id+1][id] = -1
      } else {
        // connect new vertex with parent's head
        g.am[i][id+1] = 1
        g.am[id+1][i] = -1
      }
    }
  }

  for i := range g.am[id] {
    if i == id+1 {
      continue
    }
    if g.am[id][i] == 1 {
      r := rand.Float64()
      if r > 0.9 && g.am[id][id+1] > 0 {
        // take over parent's foot
        g.am[id+1][i] = 1
        g.am[i][id+1] = -1
        g.am[id][i] = 0
        g.am[i][id+1] = 0
      } else {
        // do not take over parent's foot
        g.am[id+1][i] = 1
        g.am[i][id+1] = -1
      }
    }
  }
  g.size++
}

func New(size int) *Graph {
  _g := make([][]int, size)
  for i := range _g {
    _g[i] = make([]int, size)
  }
  _g[0][1] = 1
  _g[1][0] = -1
  _g[size-1][1] = -1
  _g[1][size-1] = 1

  g := Graph{_g, 3}
  for g.size < size {
    g.mitosis()
  }
  return &g
}

func (g *Graph) Print() {
  for _, row := range g.am {
    fmt.Println(row)
  }
}

func (g *Graph) Draw() {
  gv := graph.New(graph.IntHash, graph.Directed())
  for i := range g.am {
    _ = gv.AddVertex(i)
  }
  for i := range g.am {
    for j := range g.am[i] {
      if g.am[i][j] > 0 {
        _ = gv.AddEdge(i,j)
      }
    }
  }
  file, _ := os.Create("draw.dot")
  _ = draw.DOT(gv, file)
}

func (g *Graph) Ids() []int {
  ids := make([]int, g.size)
  for i := range g.am {
    ids[i] = i
  }
  return ids
}

func (g *Graph) Neighbours(id int) ([]int, []int) {
  pNbs := make([]int, 0)
  nNbs := make([]int, 0)
  for i, v := range g.am[id] {
    switch v {
    case 1:
      pNbs = append(pNbs, i)
    case -1:
      pNbs = append(nNbs, i)
    }
  }
  return pNbs, nNbs
}
