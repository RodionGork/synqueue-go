package queue

type Queue[E any] interface {
	Send(elem *E)
	Receive() *E
}

type ChanWrap[E any] interface {
	Send(elem *E)
	Chan() <-chan *E
	Close()
}
