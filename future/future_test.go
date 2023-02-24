package future_test

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/onur1/datatypes/event"
	"github.com/onur1/datatypes/future"
	"github.com/onur1/datatypes/result"
	"github.com/stretchr/testify/assert"
)

func fatalerror(x any) error {
	return fmt.Errorf("fatal: %v", x)
}

func double(n int) int {
	return n * 2
}

var (
	errFailed = errors.New("failed")
	errFirst  = errors.New("first")
	errSecond = errors.New("second")
)

func TestFuture(t *testing.T) {
	testCases := []struct {
		desc      string
		future    future.Future[int]
		expected  []result.Result[int]
		unordered bool
	}{
		{
			desc:     "Succeed",
			future:   future.Succeed(42),
			expected: []result.Result[int]{result.Succeed(42)},
		},
		{
			desc:     "Fail",
			future:   future.Fail[int](errFailed),
			expected: []result.Result[int]{result.Fail[int](errFailed)},
		},
		{
			desc:   "Success",
			future: future.Success(event.From([]int{1, 2, 3})),
			expected: []result.Result[int]{
				result.Succeed(1),
				result.Succeed(2),
				result.Succeed(3),
			},
		},
		{
			desc: "Failure",
			future: future.Failure[int](event.From([]error{
				errFirst,
				errSecond,
			})),
			expected: []result.Result[int]{
				result.Fail[int](errFirst),
				result.Fail[int](errSecond),
			},
		},
		{
			desc:     "After",
			future:   future.After(time.Millisecond*1, 42),
			expected: []result.Result[int]{result.Succeed(42)},
		},
		{
			desc:     "FailAfter",
			future:   future.FailAfter[int](time.Millisecond*1, errFailed),
			expected: []result.Result[int]{result.Fail[int](errFailed)},
		},
		{
			desc: "Attempt (succeed)",
			future: future.Attempt(func() (int, error) {
				return 42, nil
			}, fatalerror),
			expected: []result.Result[int]{result.Succeed(42)},
		},
		{
			desc: "Attempt (fail)",
			future: future.Attempt(func() (int, error) {
				return 0, errFailed
			}, fatalerror),
			expected: []result.Result[int]{result.Fail[int](errFailed)},
		},
		{
			desc: "Attempt (panic)",
			future: future.Attempt(func() (int, error) {
				panic("barbaz")
			}, fatalerror),
			expected: []result.Result[int]{result.Fail[int](errors.New("fatal: barbaz"))},
		},
		{
			desc:     "Map (succeed)",
			future:   future.Map(future.Succeed(42), double),
			expected: []result.Result[int]{result.Succeed(84)},
		},
		{
			desc:     "Map (fail)",
			future:   future.Map(future.Fail[int](errFailed), double),
			expected: []result.Result[int]{result.Fail[int](errFailed)},
		},
		{
			desc:     "Ap (succeed)",
			future:   future.Ap(future.Success(event.From([]func(int) int{double})), future.Succeed(42)),
			expected: []result.Result[int]{result.Succeed(84)},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			assertEq(t, tC.future, tC.expected, tC.unordered)
		})
	}
}

func assertEq(t *testing.T, dequeue future.Future[int], expected []result.Result[int], unordered bool) {
	r := make(chan result.Result[int])

	go dequeue(context.TODO(), r)

	i := 0
	l := len(expected)

	for fn := range r {
		actualValue, actualErr := fn()
		expectedValue, expectedErr := expected[i]()
		if expectedErr != nil {
			assert.Equal(t, expectedErr.Error(), actualErr.Error())
		} else {
			assert.Equal(t, expectedValue, actualValue)
		}
		i += 1
	}

	assert.Equal(t, l, i)
}
