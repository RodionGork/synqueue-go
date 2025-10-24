# SynQueue for Go language

Channels in Go languages are powerful but they always want us to specify buffer limit.

This is not convenient and violates principle that applications (usually) should not impose
artificial limits (especially when it is hard to think what sensible amount should be
set as limit - and it makes no sense to extract this to user settings).

There are various approaches to workaround the problem, here we'll try to implement
thread-safe single-linked queue of elements. We'll start with basic implementation based
on mutex and later probably enrich the solution with other approaches and utility methods.

### Public Types and Methods

```go
package queue 

type Queue[E any] interface {
    Send(elem *E)
    Receive() *E
}

// creates new Queue synchronized internally using mutex
func NewSynqueue[E any]() Queue[E]

type ChanWrap[E any] interface {
    Send(elem *E)
    Chan() <-chan *E
    Close()
}

// creates new Channel Wrapper for the given queue
func NewChanWrap[E any](queue Queue[E]) ChanWrap[E]
```

### Usage

First we'll want to add library to the project, e.g.

    go get github.com/rodiongork/synqueue-go

By now here is only one implementation, `synqueue`, using mutex internally.

Create it like this, for certain type of contained elements, e.g. `User`:

```go
type User struct {
    Id int // example field
    // ... some more fields
}

queue := NewSynqueue[User]()
```

Then add elements to queue with `Send` method:

```go
for id := 0; id < 1000; id++ {
    queue.Send(&User{Id: id})
}
```

Fetching elements from the queue is done with `Receive` method:

```go
for {
    user := queue.Receive()
    if user == nil {        // queue is currently empty
        time.Sleep(100 * time.Millisecond)
        continue
    }
    fmt.Printf("user#%d\n", user.Id)
}
```

No need to "close" the queue, just discard it when it is not needed.

**Wrapping with channel**

The `ChanWrap` is here as a ready wrapper for the `Queue` if you want to use it in
the context where go `chan` is needed (e.g. in `select` statements, `for ... range`
or the right side of `<-` operator).

The working part of the example above may be rewritten the following way, using this feature:

```go
// User[E] and queue are declared just the same

w := NewChanWrap(queue)

for id := 0; id < 1000; id++ {
    w.Send(&User{Id: id})
}

for user := range w.Chan() {
    fmt.Printf("user#%d\n", user.Id)
}
```

Note that `for range` on the channel loops until the channel is closed, so in the given
example it won't exit by itself. You'll need to call `w.Close()` from some different
goroutine (or use `select` in some way).
