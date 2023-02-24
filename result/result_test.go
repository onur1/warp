package result_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/onur1/datatypes/result"
	"github.com/stretchr/testify/assert"
)

var (
	errFailed  = errors.New("failed")
	errWrapped = fmt.Errorf("wrapped: %w", errFailed)
)

func TestResult(t *testing.T) {
	testCases := []struct {
		desc        string
		result      result.Result[int]
		expected    int
		expectedErr error
	}{
		{
			desc:     "Succeed",
			result:   result.Succeed(42),
			expected: 42,
		},
		{
			desc:        "Fail",
			result:      result.Fail[int](errFailed),
			expectedErr: errFailed,
		},
		{
			desc:     "Map (succeed)",
			result:   result.Map(result.Succeed(1), double),
			expected: 2,
		},
		{
			desc:        "Map (fail)",
			result:      result.Map(result.Fail[int](errFailed), double),
			expectedErr: errFailed,
		},
		{
			desc:     "MapError (succeed)",
			result:   result.MapError(result.Succeed(42), wrappedError),
			expected: 42,
		},
		{
			desc:        "MapError (fail)",
			result:      result.MapError(result.Fail[int](errFailed), wrappedError),
			expectedErr: errWrapped,
		},
		{
			desc: "Bimap (succeed)",
			result: result.Map(
				result.Bimap(result.Succeed(-1), wrappedError, isPositive),
				func(n bool) int {
					if n {
						return 1
					} else {
						return 2
					}
				},
			),
			expected: 2,
		},
		{
			desc: "Bimap (fail)",
			result: result.Map(
				result.Bimap(result.Fail[int](errFailed), wrappedError, isPositive),
				func(n bool) int {
					if n {
						return 1
					} else {
						return 2
					}
				},
			),
			expectedErr: errWrapped,
		},
		{
			desc:     "Ap (succeed)",
			result:   result.Ap(result.Succeed(double), result.Succeed(42)),
			expected: 84,
		},
		{
			desc:        "Ap (fail)",
			result:      result.Ap(result.Succeed(double), result.Fail[int](errFailed)),
			expectedErr: errFailed,
		},
		{
			desc:     "ApFirst (succeed)",
			result:   result.ApFirst(result.Succeed(1), result.Succeed(2)),
			expected: 1,
		},
		{
			desc:        "ApFirst (fail)",
			result:      result.ApFirst(result.Fail[int](errFailed), result.Succeed(2)),
			expectedErr: errFailed,
		},
		{
			desc:     "ApSecond (succeed)",
			result:   result.ApSecond(result.Succeed(1), result.Succeed(2)),
			expected: 2,
		},
		{
			desc:        "ApSecond (fail)",
			result:      result.ApSecond(result.Succeed(1), result.Fail[int](errFailed)),
			expectedErr: errFailed,
		},
		{
			desc: "Chain (succeed)",
			result: result.Chain(result.Succeed(42), func(a int) result.Result[int] {
				return result.Succeed(a + 1)
			}),
			expected: 43,
		},
		{
			desc: "Chain (fail)",
			result: result.Chain(result.Fail[int](errFailed), func(a int) result.Result[int] {
				return result.Succeed(a + 1)
			}),
			expectedErr: errFailed,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			assertEq(t, tC.result, tC.expected, tC.expectedErr)
		})
	}
}

func assertEq(t *testing.T, res result.Result[int], expected int, expectedErr error) {
	x, err := res()
	if err != nil {
		assert.Equal(t, expectedErr, err)
	} else {
		assert.Equal(t, expected, x)
	}
}

func double(n int) int {
	return n * 2
}

func isPositive(n int) bool {
	return n > 0
}

func wrappedError(err error) error {
	return fmt.Errorf("wrapped: %w", err)
}
