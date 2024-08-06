package warp_test

import (
	"context"
	"errors"
	"fmt"

	"github.com/onur1/warp"
	"github.com/onur1/warp/event"
	"github.com/onur1/warp/result"
)

func ExampleResult() {
	errEmptySlice := errors.New("empty slice")

	// head returns the first value contained in a slice, or fails with
	// an empty slice error.
	head := func(as []float64) warp.Result[float64] {
		if len(as) == 0 {
			return result.Error[float64](errEmptySlice)
		}
		return result.Ok(as[0])
	}

	errDivByZero := errors.New("cannot divide by zero")

	// inverse returns an inverse of a number, or fails with
	// division by zero error if 0 is encountered.
	inverse := func(n float64) warp.Result[float64] {
		if n == 0 {
			return result.Error[float64](errDivByZero)
		}
		return result.Ok(1 / n)
	}

	double := func(n float64) float64 {
		return n * 2
	}

	check := func(nums []float64) string {
		return result.Reduce(
			context.Background(),
			result.Chain(
				// double first number
				result.Map(head(nums), double),
				// take its inverse
				inverse,
			),
			// error handler
			func(err error) string {
				return fmt.Sprintf("Error is %v", err)
			},
			// success handler
			func(head float64) string {
				return fmt.Sprintf("Result is %.3f", head)
			},
		)
	}

	fmt.Println(check([]float64{24, 25, 26}))

	fmt.Println(check([]float64{0, 1, 2}))

	fmt.Println(check([]float64{}))

	// Output:
	// Result is 0.021
	// Error is cannot divide by zero
	// Error is empty slice
}

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

func double(n int) int {
	return n * 2
}

func add(b, a int) int {
	return b + a
}
