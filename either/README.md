# either

An Either represents a value pair which contains two values that can never co-exist, either the left one the right one will have [zero value](https://go.dev/ref/spec#The_zero_value).

```go
type Either[E, A any] func() (E, A)
```
