package queue

import "sync"

type Queue[E any] interface {
	Send(elem *E)
	Receive() *E
}

type elemwrap[E any] struct {
	elem *E
	next *elemwrap[E]
}

type synqueue[E any] struct {
	head *elemwrap[E]
	last *elemwrap[E]
	mu   *sync.Mutex
}

func NewSynqueue[E any]() Queue[E] {
	return &synqueue[E]{nil, nil, &sync.Mutex{}}
}

func (q *synqueue[E]) Send(elem *E) {
	q.mu.Lock()
	defer q.mu.Unlock()
	if q.last == nil {
		q.last = &elemwrap[E]{elem, nil}
		q.head = q.last
	} else {
		q.last.next = &elemwrap[E]{elem, nil}
		q.last = q.last.next
	}
}

func (q *synqueue[E]) Receive() *E {
	q.mu.Lock()
	defer q.mu.Unlock()
	if q.head == nil {
		return nil
	}
	elem := q.head.elem
	q.head = q.head.next
	if q.head == nil {
		q.last = nil
	}
	return elem
}
