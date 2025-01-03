package synctypes
import (
  "sync"
)

type Int struct {
  val int
  m sync.Mutex
}
func (i *Int) Set(x int) {
  i.val = x
}
func (i *Int) Sum(x int) {
  i.val += x
}
func (i *Int) Get() int {
  return i.val
}
func (i *Int) Lock() {
  i.m.Lock()
}
func (i *Int) Unlock() {
  i.m.Unlock()
}

type Float struct {
  val float64
  m sync.Mutex
}
func (f *Float) Set(x float64) {
  f.val = x
}
func (f *Float) Get() float64 {
  return f.val
}
func (f *Float) Lock() {
  f.m.Lock()
}
func (f *Float) Unlock() {
  f.m.Unlock()
}

type Array[T any] struct {
  val []T
  m sync.Mutex
}
func (a *Array[T]) Pop() T {
  head := a.val[0]
  a.val = a.val[1:]
  return head
}
func (a *Array[T]) Push(tail T) {
  a.val = append(a.val, tail)
}
func (a *Array[T]) Size() int {
  size := len(a.val)
  return size
}
func (a *Array[T]) Lock() {
  a.m.Lock()
}
func (a *Array[T]) Unlock() {
  a.m.Unlock()
}
func NewArray[T any]() Array[T] {
  val := make([]T, 0)
  return Array[T]{
    val: val,
  }
}
