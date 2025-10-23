# SynQueue for Go language

**currently it is a work in progress, being transferred/cleaned from private project**

Channels in Go languages are powerful but they always want us to specify buffer limit.

This is not convenient and violates principle that applications should not impose some
artificial limits (especially when it is hard to think what sensible amount should be
set as limit - and it makes no sense to extract this to user settings).

_Of course there sometimes are cases when we do want to control the buffer size -
and moreover one special case of size equal to 0 for unbuffered channel._

There are various approaches to workaround the problem, here we'll try to implement
thread-safe single-linked queue of elements. We'll start with basic implementation based
on mutex and later probably enrich the solution with other approaches and utility methods.

### Usage

By now here is only one implementation, `synqueue`, using mutex internally.

Create it like this, for certain type of contained elements, e.g. `User`:

    type User struct {
        Id int // example field
        // ... some more fields
    }

    queue := NewSynqueue[User]()

Then add elements to queue with `Send` method:

    for id := 0; id < 1000; id++ {
        queue.Send(&User{Id: id})
    }

Fetching elements from the queue is done with `Receive` method:

    for {
        user := queue.Receive()
        if user == nil {        // queue is currently empty
            time.Sleep(100 * time.Millisecond)
            continue
        }
        fmt.Printf("user#%d\n", user.Id)
    }

No need to "close" the queue, just discard it when it is not needed.
