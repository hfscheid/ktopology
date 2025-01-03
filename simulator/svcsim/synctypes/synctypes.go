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

type IncMap[K comparable] struct {
  val map[K]int
  m sync.Mutex
}
func NewMap[K comparable]() IncMap[K] {
  val := make(map[K]int)
  return IncMap[K]{
    val: val,
  }
}
func (m *IncMap[K]) Set(key K, value int) {
  m.val[key] = value
}
func (m *IncMap[K]) Inc(key K) {
  value, ok := m.val[key]
  if !ok {
    m.val[key] = 1
  } else {
    m.val[key] = value+1
  }
}
func (m *IncMap[K]) Get(key K) (int, bool) {
  value, ok := m.val[key]
  return value, ok
}
func (m *IncMap[K]) Lock() {
  m.m.Lock()
}
func (m *IncMap[K]) Unlock() {
  m.m.Unlock()
}
func (m *IncMap[K]) Keys() []K {
  keys := make([]K, 0, len(m.val))
  for k := range m.val {
    keys = append(keys, k)
  }
  return keys
}
