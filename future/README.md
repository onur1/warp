# future

A Future represents an Event that may fail, in that, it returns a value which is encapsulated in a [Result](../result).

```go
type Future[A any] Event[Result[A]]
```
