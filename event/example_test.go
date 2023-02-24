package event_test

import (
	"context"
	"fmt"

	"github.com/onur1/datatypes/event"
)

func ExampleEvent() {
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

	// Output:
	// 2
	// 6
	// 12
}
