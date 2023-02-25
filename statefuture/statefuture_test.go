package statefuture_test

import (
	"context"
	"testing"

	"github.com/onur1/data"
	"github.com/onur1/data/result"
	"github.com/onur1/data/statefuture"
	"github.com/stretchr/testify/assert"
)


func double(n int) int {
	return n * 2
}

func TestStateFuture(t *testing.T) {
	testCases := []struct {
		desc        string
		stateFuture data.StateFuture[int, int]
		expected    []data.Result[[]int]
		expectedErr error
		init        int
		unordered   bool
	}{
		{
			desc:        "Succeed",
			stateFuture: statefuture.Succeed[int](42),
			expected:    []data.Result[[]int]{result.Succeed([]int{42, 0})},
		},
		{
			desc:        "Map",
			stateFuture: statefuture.Map(statefuture.Succeed[int](42), double),
			expected:    []data.Result[[]int]{result.Succeed([]int{84, 0})},
		},
		{
			desc:        "Ap",
			stateFuture: statefuture.Ap(statefuture.Succeed[int](double), statefuture.Succeed[int](2)),
			expected:    []data.Result[[]int]{result.Succeed([]int{4, 0})},
		},
		{
			desc: "Chain",
			stateFuture: statefuture.Chain(statefuture.Succeed[int](42), func(n int) data.StateFuture[int, int] {
				return statefuture.Succeed[int](26)
			}),
			expected: []data.Result[[]int]{result.Succeed([]int{26, 0})},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			assertEq(t, tC.stateFuture, tC.expected, tC.expectedErr, tC.unordered, tC.init)
		})
	}
}

func assertEq(t *testing.T, sf data.StateFuture[int, int], expected []data.Result[[]int], expectedErr error, unordered bool, init int) {
	dequeue := sf(init)

	r := make(chan data.Result[data.These[int, int]])

	go dequeue(context.TODO(), r)

	i := 0
	l := len(expected)

	for next := range r {
		getThese, actualErr := next()
		expectedValue, expectedErr := expected[i]()
		if expectedErr != nil {
			assert.Equal(t, expectedErr.Error(), actualErr.Error())
		} else {
			a, s := getThese()
			assert.Equal(t, expectedValue, []int{a, s})
		}
		i += 1
	}

	assert.Equal(t, l, i)
}
