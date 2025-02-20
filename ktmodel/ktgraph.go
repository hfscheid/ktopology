package ktmodel
import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/dominikbraun/graph"
	"github.com/dominikbraun/graph/draw"
)

type TopologyData struct {
  ID          string          `json:"id"`
  Addr        string          `json:"podIp"`
  Service     string          `json:"host"`
  Host        string          `json:"service"`
  Deployment  string          `json:"deployment"`
  Metrics     TopologyMetrics `json:"metrics"`
}

func (t *TopologyData) StringMap() map[string]string {
  sm := make(map[string]string)
  for k, v := range t.Metrics.Metrics {
    sm[k] = fmt.Sprintf("%v", v)
  }
  return sm
}

type edge struct {
  Source string `json:"source"`
  Target string `json:"target"`
}

type NwCalculus func([]TopologyData)(string, float64)

type NwTopology struct {
  Timestamp time.Time           `json:"timestamp"`
  Nodes     []TopologyData      `json:"nodes"`
  Edges     []edge              `json:"edges"`
  NwData    map[string]float64  `json:"network_data"`
  nwOps     []NwCalculus        `json:"-"`
}

func (g *NwTopology) AddNwCalculus(nwc NwCalculus) {
  g.nwOps = append(g.nwOps, nwc)
}

func (g *NwTopology) DoNwCalculus() {
  g.NwData = make(map[string]float64)
  for i := range g.nwOps {
    key, value := g.nwOps[i](g.Nodes)
    g.NwData[key] = value
  }
}

func (g *NwTopology) ToJSON() ([]byte, error) {
  return json.Marshal(*g)
}

func FromJSON(data []byte) ([]*NwTopology, error) {
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
  series := make([]*NwTopology, 0, len(lines))
  for i := range lines {
    if len(lines[i]) == 0 {
      continue
    }
    g := &NwTopology{}
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

func (g *NwTopology) DrawPods() {
  drawNwTopology := graph.New(graph.StringHash, graph.Directed())
  for _, n := range g.Nodes {
    strAttributes := n.StringMap()
    _ = drawNwTopology.AddVertex(n.ID, graph.VertexAttributes(strAttributes))
  }
  for _, e := range g.Edges {
    _ = drawNwTopology.AddEdge(e.Source, e.Target)
  }
  file, _ := os.Create("pods.dot")
  _ = draw.DOT(drawNwTopology, file)
}

func (g *NwTopology) DrawDeployments() {
  drawNwTopology := graph.New(graph.StringHash, graph.Directed())
  podMap := make(map[string]string)
  dMap := make(map[string]int)
  vertexMap := make(map[string]string)
  for _, n := range g.Nodes {
    val, ok := dMap[n.Deployment]
    if !ok {
      dMap[n.Deployment] = 1
    } else {
      dMap[n.Deployment] = val+1
    }
    podMap[n.ID] = n.Deployment
  }
  for k, v := range dMap {
    vertexName := fmt.Sprintf("%v replicas: %v", k, v)
    _ = drawNwTopology.AddVertex(vertexName)
    vertexMap[k] = vertexName
  }
  for _, e := range g.Edges {
    _ = drawNwTopology.AddEdge(
      vertexMap[podMap[e.Source]],
      vertexMap[podMap[e.Target]],
    )
  }
  file, _ := os.Create("deployments.dot")
  _ = draw.DOT(drawNwTopology, file)
}

func buildAddrMap(tData []TopologyData) map[string][]string {
  addrMap := make(map[string][]string)
  for _, tDatum := range tData {
    if  _, ok := addrMap[tDatum.Service];
        !ok {
          addrMap[tDatum.Service] = make([]string, 0)
    }
    addrMap[tDatum.Service] = append(addrMap[tDatum.Service], tDatum.ID)
  }
  return addrMap
}

func BuildNwTopologyFromData(tData []TopologyData) *NwTopology {
  addrMap := buildAddrMap(tData)
  edges := make([]edge, 0, len(tData))
  for i := range tData {
    for addr := range tData[i].Metrics.SentPkgs {
      for _, target := range addrMap[addr] {
        edge := edge{
          Source: tData[i].ID,
          Target: target,
        }
        edges = append(edges, edge)
      }
    }
  }
  graph := &NwTopology{
    Timestamp: time.Now(),
    Nodes: tData,
    Edges: edges,
  }
  return graph
}
