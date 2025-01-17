package main
import (
  "testpoddiscovery/poddiscovery"
  "fmt"
)

func main() {
  fmt.Println("Listing pods...")
  fmt.Println(poddiscovery.ListPods())
}
