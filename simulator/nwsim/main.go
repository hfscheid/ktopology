package main
import (
  "nwsim/randomflux"
)

func main() {
  g := randomflux.New(20)
  g.Print()
  g.Draw()
}
