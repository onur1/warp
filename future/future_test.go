package future_test

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/onur1/warp"
	"github.com/onur1/warp/event"
	"github.com/onur1/warp/future"
	"github.com/onur1/warp/result"
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
		future    warp.Future[int]
		expected  []warp.Result[int]
		unordered bool
	}{
		{
			desc:     "Succeed",
			future:   future.Succeed(42),
			expected: []warp.Result[int]{result.Ok(42)},
		},
		{
			desc:     "Fail",
			future:   future.Fail[int](errFailed),
			expected: []warp.Result[int]{result.Error[int](errFailed)},
		},
		{
			desc:   "Success",
			future: future.Success(event.From([]int{1, 2, 3})),
			expected: []warp.Result[int]{
				result.Ok(1),
				result.Ok(2),
				result.Ok(3),
			},
		},
		{
			desc: "Failure",
			future: future.Failure[int](event.From([]error{
				errFirst,
				errSecond,
			})),
			expected: []warp.Result[int]{
				result.Error[int](errFirst),
				result.Error[int](errSecond),
			},
		},
		{
			desc:     "After",
			future:   future.After(time.Millisecond*1, 42),
			expected: []warp.Result[int]{result.Ok(42)},
		},
		{
			desc:     "FailAfter",
			future:   future.FailAfter[int](time.Millisecond*1, errFailed),
			expected: []warp.Result[int]{result.Error[int](errFailed)},
		},
		{
			desc: "Attempt (succeed)",
			future: future.Attempt(func(_ context.Context) (int, error) {
				return 42, nil
			}, fatalerror),
			expected: []warp.Result[int]{result.Ok(42)},
		},
		{
			desc: "Attempt (fail)",
			future: future.Attempt(func(_ context.Context) (int, error) {
				return 0, errFailed
			}, fatalerror),
			expected: []warp.Result[int]{result.Error[int](errFailed)},
		},
		{
			desc: "Attempt (panic)",
			future: future.Attempt(func(_ context.Context) (int, error) {
				panic("barbaz")
			}, fatalerror),
			expected: []warp.Result[int]{result.Error[int](errors.New("fatal: barbaz"))},
		},
		{
			desc:     "Map (succeed)",
			future:   future.Map(future.Succeed(42), double),
			expected: []warp.Result[int]{result.Ok(84)},
		},
		{
			desc:     "Map (fail)",
			future:   future.Map(future.Fail[int](errFailed), double),
			expected: []warp.Result[int]{result.Error[int](errFailed)},
		},
		{
			desc:     "Ap (succeed)",
			future:   future.Ap(future.Success(event.From([]func(int) int{double})), future.Succeed(42)),
			expected: []warp.Result[int]{result.Ok(84)},
		},
		{
			desc:     "FromEvent",
			future:   future.FromEvent(event.Of(42)),
			expected: []warp.Result[int]{result.Ok(42)},
		},
		{
			desc:     "From",
			future:   future.From([]int{42, 43, 44}),
			expected: []warp.Result[int]{result.Ok(42), result.Ok(43), result.Ok(44)},
		},
		{
			desc: "Parallel",
			future: future.Parallel(
				warp.Future[int](event.From([]warp.Result[int]{
					result.After(time.Millisecond*20, 1),
					result.Ok(2),
					result.Ok(3),
					result.After(time.Millisecond*10, 4),
					result.Ok(5),
				})),
				2,
			),
			expected: []warp.Result[int]{result.Ok(1), result.Ok(2), result.Ok(3), result.Ok(4), result.Ok(5)},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			assertEq(t, tC.future, tC.expected, tC.unordered)
		})
	}
}

func assertEq(t *testing.T, dequeue warp.Future[int], expected []warp.Result[int], unordered bool) {
	r := make(chan warp.Result[int])

	go dequeue(context.TODO(), r)

	i := 0
	l := len(expected)

	for fn := range r {
		actualValue, actualErr := fn(context.TODO())
		expectedValue, expectedErr := expected[i](context.TODO())
		if expectedErr != nil {
			assert.Equal(t, expectedErr.Error(), actualErr.Error())
		} else {
			assert.Equal(t, expectedValue, actualValue)
		}
		i += 1
	}

	assert.Equal(t, l, i)
}
