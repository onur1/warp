package state_test

import (
	"testing"

	"github.com/onur1/data"
	"github.com/onur1/data/state"
	"github.com/stretchr/testify/assert"
)

func TestState(t *testing.T) {
	testCases := []struct {
		desc     string
		state    data.State[int, int]
		expected []int
		init     int
	}{
		{
			desc: "Map",
			state: state.Map(func(s int) (int, int) {
				return s - 1, s + 1
			}, double),
			expected: []int{-2, 1},
		},
		{
			desc:     "Ap",
			state:    state.Ap(state.Of[int](double), state.Of[int](1)),
			expected: []int{2, 0},
		},
		{
			desc: "Chain",
			state: state.Chain(func(s int) (int, int) {
				return s - 1, s + 1
			}, func(int) data.State[int, int] {
				return func(s int) (int, int) {
					return s - 1, s + 1
				}
			}),
			expected: []int{0, 2},
		},
		{
			desc:     "ApFirst",
			state:    state.ApFirst(state.Of[int](1), state.Of[int](2)),
			expected: []int{1, 0},
		},
		{
			desc:     "ApSecond",
			state:    state.ApSecond(state.Of[int](1), state.Of[int](2)),
			expected: []int{2, 0},
		},
		{
			desc:     "Put",
			state:    state.Put[int, int](2),
			expected: []int{0, 2},
			init:     1,
		},
		{
			desc:     "Get",
			state:    state.Get[int](),
			expected: []int{1, 1},
			init:     1,
		},
		{
			desc:     "Modify",
			state:    state.Modify[int, int](double),
			expected: []int{0, 2},
			init:     1,
		},
		{
			desc:     "Gets",
			state:    state.Gets(double),
			expected: []int{2, 1},
			init:     1,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			assertEq(t, tC.state, tC.expected, tC.init)
		})
	}
}

func assertEq(t *testing.T, fa data.State[int, int], expected []int, init int) {
	a, s := fa(init)
	assert.Equal(t, expected, []int{a, s})
}

func double(n int) int {
	return n * 2
}
