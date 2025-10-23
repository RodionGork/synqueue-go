# SynQueue for Go language

Channels in Go languages are powerful but they always want us to specify buffer limit.

This is not convenient and violates principle that applications should not impose some
artificial limits (especially when it is hard to think what sensible amount should be
set as limit - and it makes no sense to extract this to user settings).

_Of course there sometimes are cases when we do want to control the buffer size -
and moreover one special case of size equal to 0 for unbuffered channel._

There are various approaches to workaround the problem, here we'll try to implement
thread-safe single-linked queue of elements. We'll start with basic implementation based
on mutex and later probably enrich the solution with other approaches and utility methods.
