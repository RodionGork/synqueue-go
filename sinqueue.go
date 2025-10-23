package main

type Queue[E any] interface {
	Add(elem *E)
	Take() *E
}

type elemwrap[E any] struct {
	elem *E
	next *elemwrap[E]
}

type sinqueue[E any] struct {
	head *elemwrap[E]
	last *elemwrap[E]
}

func NewSinqueue[E any]() *sinqueue[E] {
	return &sinqueue[E]{nil, nil}
}

func (q *sinqueue[E]) Add(elem *E) {
	if q.last == nil {
		q.last = &elemwrap[E]{elem, nil}
		q.head = q.last
	} else {
		q.last.next = &elemwrap[E]{elem, q.head}
		q.last = q.last.next
	}
}

func (q *sinqueue[E]) Take() *E {
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

type User struct {
	Id int
}

func main() {
	q := NewSinqueue[User]()
	q.Add(&User{2})
	println(q.Take().Id)
	q.Add(&User{3})
	q.Add(&User{5})
	println(q.Take().Id)
	println(q.Take().Id)
}
