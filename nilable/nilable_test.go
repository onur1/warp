package nilable_test

import (
	"errors"
	"testing"

	"github.com/onur1/data"
	"github.com/onur1/data/nilable"
	"github.com/onur1/data/result"
	"github.com/stretchr/testify/assert"
)

func TestNilable(t *testing.T) {
	testCases := []struct {
		desc     string
		nilable  data.Nilable[int]
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
			nilable:  nilable.FromResult(result.Succeed(42)),
			expected: 42,
		},
		{
			desc:    "FromResult (fail)",
			nilable: nilable.FromResult(result.Fail[int](errors.New("some error"))),
		},
		{
			desc: "FromNullable",
			nilable: func() data.Nilable[int] {
				var twelve = int(12)
				return nilable.FromNullable(&twelve)
			}(),
			expected: 12,
		},
		{
			desc:    "FromNullable (nil)",
			nilable: nilable.FromNullable[int](nil),
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

func assertEq(t *testing.T, v data.Nilable[int], expected int) {
	if v == nil {
		assert.Equal(t, expected, 0)
	} else {
		assert.Equal(t, expected, *v)
	}
}

func double(n int) int {
	return n * 2
}
