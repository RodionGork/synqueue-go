package queue

import (
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
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
	if q.Receive() != nil {
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
	q.Send(&User{2})
	if q.Receive().Id != 2 {
		t.Errorf("Failure on taking the only element from queue")
	}
	checkLoop(t, q.(*synqueue[User]))
	if q.Receive() != nil {
		t.Errorf("Non-empty queue after taking the only element")
	}
}

func TestSeveralElems(t *testing.T) {
	ids := []int{2, 3, 5, 7, 11}
	q := NewSynqueue[User]()
	for _, v := range ids {
		q.Send(&User{v})
	}
	checkLoop(t, q.(*synqueue[User]))
	for i, v := range ids {
		elem := q.Receive()
		if elem == nil {
			t.Errorf("Unexpected emptiness in queue (step %v)", i)
		}
		if elem.Id != v {
			t.Errorf("Wrong element in queue (step %v, value %v)", i, v)
		}
	}
	if elem := q.Receive(); elem != nil {
		t.Errorf("Non-empty queue after taking all elements (%v)", elem)
	}
}

var atomicCnt int64

func fillQueue(q Queue[User], step, start, count int) {
	id := start
	defer func() {
		println("atomicnt", atomicCnt)
	}()
	for i := 0; i < count; i++ {
		q.Send(&User{id})
		atomic.AddInt64(&atomicCnt, 1)
		id += step
	}
}

func TestManyInParallel(t *testing.T) {
	q := NewSynqueue[User]()
	wg := sync.WaitGroup{}
	par := 5
	num := 1000000
	wg.Add(par)
	atomic.StoreInt64(&atomicCnt, 0)
	for i := 0; i < par; i++ {
		go func() {
			fillQueue(q, par, i, num)
			wg.Done()
		}()
	}
	wg.Wait()
	cnt := 0
	for elem := q.(*synqueue[User]).head; elem != nil; elem = elem.next {
		cnt++
	}
	if cnt != num*par {
		t.Errorf("Incorrect elems count instead of %v: %v", num*par, cnt)
	}
}
