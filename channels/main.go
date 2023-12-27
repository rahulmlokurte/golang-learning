package main

import (
	"container/list"
	"errors"
	"sync"
	"sync/atomic"
	"time"
)

var (
  ErrBufferFull = errors.New("buffer is full")
)

type Channel interface {
  Send(interface{})
  Recv() (interface{}, bool)
  Close()
  Next() bool
}

type BufferedChannel struct {
  buf *buffer
  lock sync.Mutex
  closed bool
  sendCounter int32
  sendQ *list.List
}

type buffer struct {
  q *list.List
  maxLen int
}

func newBuffer(size int) *buffer {
  return &buffer{
    q: new(list.List).Init(),
    maxLen: size,
  }
}

func (c *buffer) IsFull() bool {
  return c.q.Len() >= c.maxLen
}

func (c *buffer) IsEmpty() bool {
  return c.q.Len() == 0
}

func (c *buffer) Enqueue(val interface{}) error {
  if c.IsFull() {
    return ErrBufferFull
  }
  c.q.PushBack(val)
  return nil
}

func (c *buffer) Dequeue() interface{} {
  if c.IsEmpty() {
    return nil
  }
  return c.q.Remove(c.q.Front())
}

func NewChannel(size int) *BufferedChannel {
  if size > 0 {
    return &BufferedChannel{
      buf :newBuffer(size),
    }
  }
  panic("size of channel needs to be greater than 0")
}

func (c *BufferedChannel) close() {
  c.lock.Lock()
  defer c.lock.Unlock()
  c.closed = true
}

func (c *BufferedChannel) Send(val interface{}) {
  if c.closed {
    panic("channel closed")
  }
  c.lock.Lock()
  defer c.lock.Unlock()

  if !c.buf.IsFull() {
    c.buf.Enqueue(val)
    return
  }

  ticket := atomic.AddInt32(&c.sendCounter, 1)
  c.sendQ.PushBack(ticket)
  c.lock.Unlock()

  for {
    c.lock.Lock()
    if !c.buf.IsFull() && ticket == c.sendQ.Front().Value.(int32) {
      break
    }
    c.lock.Unlock()
    time.Sleep(10 * time.Millisecond)
  }
  c.buf.Enqueue(val)
  c.sendQ.Remove(c.sendQ.Front())
}

func (c *BufferedChannel) Recv() (interface{}, bool) {
  c.lock.Lock()
  defer c.lock.Unlock()
 
  if c.buf.IsEmpty() && c.closed {
    return nil, false
  }

  if !c.buf.IsEmpty(){
    return c.buf.Dequeue(), true
  }
  return nil, false
}

func (c *BufferedChannel) Next() bool {
  for {
    c.lock.Lock()

    if c.closed && c.buf.IsEmpty() {
      c.lock.Unlock()
      return false
    }
    if !c.buf.IsEmpty() {
      c.lock.Unlock()
      return true
    }
    c.lock.Unlock()
    time.Sleep(10 * time.Millisecond)
  }
}
