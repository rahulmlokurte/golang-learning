package main

import (
  "fmt"
  "os"
)

type foobar interface {
  foo()
  bar()
}

type itemA struct{}

func (a *itemA) foo() {
  fmt.Println("foo on A")
}

func (a *itemA) bar() {
  fmt.Println("foo on C")
}

type itemB struct{}

func (a *itemB) foo() {
  fmt.Println("foo on B")
}

func (b *itemB) bar() {
  fmt.Println("foo on B")
}

func doFoo(item foobar) {
  item.foo()
  item.bar()
}

func main() {
  doFoo(&itemA{})
  doFoo(&itemB{})
}



