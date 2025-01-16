package ktgraph

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/dominikbraun/graph"
	"github.com/dominikbraun/graph/draw"
)

type Node struct {
  ID          string          `json:"id"`
  Addr        string          `json:"podIp"`
  Host        string          `json:"host"`
  Service     string          `json:"service"`
  Deployment  string          `json:"deployment"`
  Metadata    map[string]any  `json:"metadata"`
}

type Edge struct {
  Source string `json:"source"`
  Target string `json:"target"`
}

type Graph struct {
  Timestamp time.Time `json:"timestamp"`
  Nodes     []Node    `json:"nodes"`
  Edges     []Edge    `json:"edges"`
}

func (g *Graph) ToJSON() ([]byte, error) {
  return json.Marshal(*g)
}

func FromJSON(data []byte) ([]*Graph, error) {
  lines := make([][]byte, 0)
  start := 0
  for i := range data {
    if data[i] == '\n' {
      line := make([]byte, (i-start)) // clip last '\n' char
      copy(line, data[start:i])
      lines = append(lines, line)
      start = i+1
    }
  }
  series := make([]*Graph, 0, len(lines))
  for i := range lines {
    if len(lines[i]) == 0 {
      continue
    }
    g := &Graph{}
    err := json.Unmarshal(lines[i], g)
    if err != nil {
      return nil, fmt.Errorf(
        "Could not unmarshal line:\n%v\n, error: %v",
        lines[i],
        err,
      )
    }
    series = append(series, g)
  }
  return series, nil
}

func (g *Graph) DrawPods() {
  drawGraph := graph.New(graph.StringHash, graph.Directed())
  for _, n := range g.Nodes {
    strAttributes := make(map[string]string)
    for k, v := range n.Metadata {
      if  t, ok := v.(int);
          ok {
        strAttributes[k] = fmt.Sprintf("%v", t)
      } else if t, ok := v.(string);
                ok {
        strAttributes[k] = t
      } else if t, ok := v.(float64);
                ok {
        strAttributes[k] = fmt.Sprintf("%.2f", t)
      }
    }
    _ = drawGraph.AddVertex(n.ID, graph.VertexAttributes(strAttributes))
  }
  for _, e := range g.Edges {
    _ = drawGraph.AddEdge(e.Source, e.Target)
  }
  file, _ := os.Create("topology.dot")
  _ = draw.DOT(drawGraph, file)
}

func (g *Graph) DrawDeployments() {
}
