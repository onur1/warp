package event_test

import (
	"context"
	"sort"
	"testing"

	"github.com/onur1/warp"
	"github.com/onur1/warp/event"
	"github.com/onur1/warp/nilable"
	"github.com/stretchr/testify/assert"
)

func TestEvent(t *testing.T) {
	testCases := []struct {
		desc      string
		event     warp.Event[int]
		expected  []int
		unordered bool
	}{
		{
			desc:     "Of",
			event:    event.Of(42),
			expected: []int{42},
		},
		{
			desc:     "From",
			event:    event.From([]int{4, 5, 6}),
			expected: []int{4, 5, 6},
		},
		{
			desc:     "Map",
			event:    event.Map(event.From([]int{1, 2, 3}), double),
			expected: []int{2, 4, 6},
		},
		{
			desc: "Chain",
			event: event.Chain(event.From([]int{1, 2, 3}), func(a int) warp.Event[int] {
				return event.From([]int{a, a + 1})
			}),
			expected: []int{1, 2, 2, 3, 3, 4},
		},
		{
			desc:     "Filter",
			event:    event.Filter(event.From([]int{-3, 4, -1, 5, 0, 6}), isPositive),
			expected: []int{4, 5, 6},
		},
		{
			desc:     "Ap",
			event:    event.Ap(event.From([](func(int) int){double}), event.From([]int{1, 2, 3, 4})),
			expected: []int{2, 4, 6, 8},
		},
		{
			desc:     "SampleOn",
			event:    event.SampleOn(event.From([]int{1}), event.From([](func(int) int){double, triple})),
			expected: []int{2, 3},
		},
		{
			desc:      "Alt",
			event:     event.Alt(event.Of(1), event.Of(2)),
			expected:  []int{1, 2},
			unordered: true,
		},
		{
			desc:     "Fold",
			event:    event.Fold(event.From([]int{1, 2}), 5, add),
			expected: []int{6, 8},
		},
		{
			desc: "WithLast",
			event: event.Map(
				event.WithLast(event.From([]int{1, 41})), func(l event.Last[int]) int {
					if l.Last == 0 {
						return l.Now
					}
					return add(l.Now, l.Last)
				}),
			expected: []int{1, 42},
		},
		{
			desc:     "Reduce",
			event:    event.Of(event.Reduce(context.TODO(), event.From([]int{1, 2, 3}), 36, add)),
			expected: []int{42},
		},
		{
			desc:     "ReduceRight",
			event:    event.Of(event.ReduceRight(context.TODO(), event.From([]int{10, 200, 1000}), 50, div)),
			expected: []int{1},
		},
		{
			desc:     "CountAll",
			event:    event.Of(event.CountAll(context.TODO(), event.From([]int{4, 5, 6}))),
			expected: []int{3},
		},
		{
			desc:     "Take",
			event:    event.Take(event.From([]int{4, 5, 6}), 2),
			expected: []int{4, 5},
		},
		{
			desc:     "FilterMap",
			event:    event.FilterMap(event.From([]int{-3, 4, -1, 5, 0, 6}), doublePositive),
			expected: []int{8, 10, 12},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			assertEq(t, tC.event, tC.expected, tC.unordered)
		})
	}
}

func assertEq(t *testing.T, dequeue warp.Event[int], expected []int, unordered bool) {
	r := make(chan int)

	go dequeue(context.TODO(), r)

	i := 0
	l := len(expected)
	var collected []int

	for v := range r {
		collected = append(collected, v)
		i += 1
	}

	assert.Equal(t, l, i)

	if unordered {
		sort.Ints(collected)
	}

	assert.Equal(t, expected, collected)
}

func double(n int) int {
	return n * 2
}

func triple(n int) int {
	return n * 3
}

func isPositive(n int) bool {
	return n > 0
}

func add(b, a int) int {
	return b + a
}

func div(b, a int) int {
	return b / a
}

func doublePositive(n int) warp.Nilable[int] {
	if !isPositive(n) {
		return nil
	}
	return nilable.Some(n * 2)
}
