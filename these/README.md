# these

A These represents a value pair which contains two values that may have non-zero values at the same time. This is "inclusive-or" (as opposed to "exclusive-or" provided by Either), both values can have non-zero values, or only one of them.

```go
type These[E, A any] func() (E, A)
```
