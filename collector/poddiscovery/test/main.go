package main
import (
  "collector/poddiscovery"
  "fmt"
)

func main() {
  fmt.Println("Listing pods...")
  fmt.Println(poddiscovery.ListPods())
}
