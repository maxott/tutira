
package queue

import (
  "container/list"
)

func NewFoo() *Foo {
  return &Foo{list.New()}
}

type Foo struct {
  list *list.List
}

func (q *Foo) Len() int {
  return q.list.Len()
}

func (q *Foo) Enqueue(i string) {
  q.list.PushBack(i)
}

func (q *Foo) Dequeue() string {
  if q.list.Len() == 0 {
    panic(ErrEmptyQueue)
  }
  raw := q.list.Remove(q.list.Front())
  if typed, ok := raw.(string); ok {
    return typed
  }
  panic(ErrInvalidType)
}
