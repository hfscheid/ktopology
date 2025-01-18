package ktmodel
import (
  "os"
  "testing"
)

func TestFromJSON(t *testing.T) {
  data, err := os.ReadFile("sample.jsonl")
  if err != nil {
    t.Fatalf("Could not parse sample file: %v\n", err)
  }
  series, err := FromJSON(data)
  if err != nil {
    t.Fatalf("Could not parse line: %v\n", err)
  }
  for _, g := range series {
    t.Log(*g)
  }
}

func TestDrawing(t *testing.T) {
  data, err := os.ReadFile("sample.jsonl")
  if err != nil {
    t.Fatalf("Could not parse sample file: %v\n", err)
  }
  series, err := FromJSON(data)
  if err != nil {
    t.Fatalf("Could not parse line: %v\n", err)
  }
  series[len(series)-1].DrawPods()
}
