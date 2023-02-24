# event

An Event represents a collection of discrete occurrences with associated values.

```go
type Event[A any] func(context.Context, chan<- A)
```

## Example

```go
func double(n int) int {
	return n * 2
}

func add(b, a int) int {
	return b + a
}

dequeue := event.Fold(
  event.Map(event.From([]int{1, 2, 3}), double),
  0,
  add,
)

r := make(chan int)

go dequeue(context.TODO(), r)

for i := range r {
  fmt.Println(i)
}
```

Output:

```
2
6
12
```
