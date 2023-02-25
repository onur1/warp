# these

A These represents a value pair which contains two values that may have zero values or not.

This is "inclusive-or" (as opposed to "exclusive-or" provided by [Either](../either)), both values can have zero values, or only one of them, or both of them can have non-zero values.

```go
type These[E, A comparable] func() (E, A)
```
