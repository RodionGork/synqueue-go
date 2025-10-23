package main

type Queue[E any] interface {
	Add(elem *E)
	Take() *E
}

type elemwrap[E any] struct {
	elem *E
	next *elemwrap[E]
}

type synqueue[E any] struct {
	head *elemwrap[E]
	last *elemwrap[E]
}

func NewSynqueue[E any]() *synqueue[E] {
	return &synqueue[E]{nil, nil}
}

func (q *synqueue[E]) Add(elem *E) {
	if q.last == nil {
		q.last = &elemwrap[E]{elem, nil}
		q.head = q.last
	} else {
		q.last.next = &elemwrap[E]{elem, nil}
		q.last = q.last.next
	}
}

func (q *synqueue[E]) Take() *E {
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
