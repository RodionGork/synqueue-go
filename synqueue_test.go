package main

import (
	"fmt"
	"strings"
	"testing"
)

type User struct {
	Id int
}

func (u User) String() string {
	return fmt.Sprintf("User{%d}", u.Id)
}

func (q *synqueue[E]) String() string {
	el := q.head
	res := []string{}
	for el != nil {
		res = append(res, fmt.Sprintf("%v", el.elem))
		el = el.next
	}
	return strings.Join(res, ", ")
}

func TestEmpty(t *testing.T) {
	q := NewSynqueue[User]()
	if q.Take() != nil {
		t.Errorf("Non-empty queue when nothing was put into it")
	}
}

func checkLoop(t *testing.T, q *synqueue[User]) {
	seen := map[string]bool{}
	elem := q.head
	seen[fmt.Sprintf("%p", elem)] = true
	for elem != nil && len(seen) < 10 {
		np := fmt.Sprintf("%p", elem.next)
		if seen[np] {
			t.Fatalf("loop detected (%v): %v", np, seen)
		}
		seen[np] = true
		elem = elem.next
	}
}

func TestSingleElem(t *testing.T) {
	q := NewSynqueue[User]()
	q.Add(&User{2})
	if q.Take().Id != 2 {
		t.Errorf("Failure on taking the only element from queue")
	}
	checkLoop(t, q)
	if q.Take() != nil {
		t.Errorf("Non-empty queue after taking the only element")
	}
}

func TestSeveralElems(t *testing.T) {
	ids := []int{2, 3, 5, 7, 11}
	q := NewSynqueue[User]()
	for _, v := range ids {
		q.Add(&User{v})
	}
	checkLoop(t, q)
	for i, v := range ids {
		elem := q.Take()
		if elem == nil {
			t.Errorf("Unexpected emptiness in queue (step %v)", i)
		}
		if elem.Id != v {
			t.Errorf("Wrong element in queue (step %v, value %v)", i, v)
		}
	}
	if elem := q.Take(); elem != nil {
		t.Errorf("Non-empty queue after taking all elements (%v)", elem)
	}
}
