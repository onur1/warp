package nilable_test

import (
	"context"
	"errors"
	"testing"

	"github.com/onur1/fpgo"
	"github.com/onur1/fpgo/nilable"
	"github.com/onur1/fpgo/result"
	"github.com/stretchr/testify/assert"
)

func TestNilable(t *testing.T) {
	testCases := []struct {
		desc     string
		nilable  fpgo.Nilable[int]
		expected int
	}{
		{
			desc:    "Nil",
			nilable: nilable.Nil[int](),
		},
		{
			desc:     "Some",
			nilable:  nilable.Some(42),
			expected: 42,
		},
		{
			desc:     "Map",
			nilable:  nilable.Map(nilable.Some(42), double),
			expected: 84,
		},
		{
			desc:    "Map (nil)",
			nilable: nilable.Map(nilable.Nil[int](), double),
		},
		{
			desc:     "Ap",
			nilable:  nilable.Ap(nilable.Some(double), nilable.Some(42)),
			expected: 84,
		},
		{
			desc:    "Ap (nil)",
			nilable: nilable.Ap(nilable.Some(double), nilable.Nil[int]()),
		},
		{
			desc:     "FromResult (succeed)",
			nilable:  nilable.FromResult(context.TODO(), result.Ok(42)),
			expected: 42,
		},
		{
			desc:    "FromResult (fail)",
			nilable: nilable.FromResult(context.TODO(), result.Error[int](errors.New("some error"))),
		},
		{
			desc: "Attempt",
			nilable: nilable.Attempt(func() int {
				return 42
			}),
			expected: 42,
		},
		{
			desc: "Attempt (nil)",
			nilable: nilable.Attempt(func() int {
				panic("")
			}),
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			assertEq(t, tC.nilable, tC.expected)
		})
	}
}

func assertEq(t *testing.T, v fpgo.Nilable[int], expected int) {
	if v == nil {
		assert.Equal(t, expected, 0)
	} else {
		assert.Equal(t, expected, *v)
	}
}

func double(n int) int {
	return n * 2
}
