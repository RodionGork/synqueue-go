package queue

import (
	"sync/atomic"
	"time"
)

type chanWrap[E any] struct {
	q Queue[E]
	c chan *E
	f atomic.Bool
}

func NewChanWrap[E any](queue Queue[E]) ChanWrap[E] {
	res := &chanWrap[E]{queue, make(chan *E, 0), atomic.Bool{}}
	go res.run()
	return res
}

func (w *chanWrap[E]) Send(elem *E) {
	w.q.Send(elem)
}

func (w *chanWrap[E]) Chan() <-chan *E {
	return w.c
}

func (w *chanWrap[E]) Close() {
	close(w.c)
	w.f.Store(true)
}

func (w *chanWrap[E]) run() {
	defer func() {
		recover()
	}()
	delay := time.Millisecond * 30
	for !w.f.Load() {
		elem := w.q.Receive()
		if elem == nil {
			time.Sleep(delay)
			continue
		}
		w.c <- elem
	}
}
