package result_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/onur1/fpgo"
	"github.com/onur1/fpgo/nilable"
	"github.com/onur1/fpgo/result"
	"github.com/stretchr/testify/assert"
)

var (
	errFailed  = errors.New("failed")
	errWrapped = fmt.Errorf("wrapped: %w", errFailed)
)

func TestResult(t *testing.T) {
	testCases := []struct {
		desc        string
		result      fpgo.Result[int]
		expected    int
		expectedErr error
	}{
		{
			desc:     "Ok",
			result:   result.Ok(42),
			expected: 42,
		},
		{
			desc:        "Error",
			result:      result.Error[int](errFailed),
			expectedErr: errFailed,
		},
		{
			desc:     "Map",
			result:   result.Map(result.Ok(1), double),
			expected: 2,
		},
		{
			desc:        "Map (error)",
			result:      result.Map(result.Error[int](errFailed), double),
			expectedErr: errFailed,
		},
		{
			desc:     "MapError",
			result:   result.MapError(result.Ok(42), wrappedError),
			expected: 42,
		},
		{
			desc:        "MapError (error)",
			result:      result.MapError(result.Error[int](errFailed), wrappedError),
			expectedErr: errWrapped,
		},
		{
			desc: "Bimap",
			result: result.Map(
				result.Bimap(result.Ok(-1), wrappedError, isPositive),
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
			desc: "Bimap (error)",
			result: result.Map(
				result.Bimap(result.Error[int](errFailed), wrappedError, isPositive),
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
			desc:     "Ap",
			result:   result.Ap(result.Ok(double), result.Ok(42)),
			expected: 84,
		},
		{
			desc:        "Ap (error)",
			result:      result.Ap(result.Ok(double), result.Error[int](errFailed)),
			expectedErr: errFailed,
		},
		{
			desc:     "ApFirst",
			result:   result.ApFirst(result.Ok(1), result.Ok(2)),
			expected: 1,
		},
		{
			desc:        "ApFirst (error)",
			result:      result.ApFirst(result.Error[int](errFailed), result.Ok(2)),
			expectedErr: errFailed,
		},
		{
			desc:     "ApSecond",
			result:   result.ApSecond(result.Ok(1), result.Ok(2)),
			expected: 2,
		},
		{
			desc:        "ApSecond (error)",
			result:      result.ApSecond(result.Ok(1), result.Error[int](errFailed)),
			expectedErr: errFailed,
		},
		{
			desc: "Chain",
			result: result.Chain(result.Ok(42), func(a int) fpgo.Result[int] {
				return result.Ok(a + 1)
			}),
			expected: 43,
		},
		{
			desc: "Chain (error)",
			result: result.Chain(result.Error[int](errFailed), func(a int) fpgo.Result[int] {
				return result.Ok(a + 1)
			}),
			expectedErr: errFailed,
		},
		{
			desc: "FromNilable (some)",
			result: result.FromNilable(nilable.Some(42), func() error {
				return errFailed
			}),
			expected: 42,
		},
		{
			desc: "FromNilable (nil)",
			result: result.FromNilable(nilable.Nil[int](), func() error {
				return errFailed
			}),
			expectedErr: errFailed,
		},
		{
			desc: "OrElse (error)",
			result: result.OrElse(result.Error[int](errFailed), func(err error) fpgo.Result[int] {
				return result.Ok(555)
			}),
			expected: 555,
		},
		{
			desc: "OrElse",
			result: result.OrElse(result.Ok(666), func(err error) fpgo.Result[int] {
				return result.Ok(555)
			}),
			expected: 666,
		},
		{
			desc: "GetOrElse (error)",
			result: result.Ok(
				result.GetOrElse(
					context.TODO(),
					result.Error[int](errFailed), func(err error) int {
						return 444
					}),
			),
			expected: 444,
		},
		{
			desc: "GetOrElse",
			result: result.Ok(
				result.GetOrElse(
					context.TODO(),
					result.Ok(42), func(err error) int {
						return 444
					}),
			),
			expected: 42,
		},
		{
			desc: "FilterOrElse",
			result: result.FilterOrElse(result.Ok(42), func(x int) bool {
				return x > 40
			}, func(x int) error {
				return fmt.Errorf("filterOrElse: %d is not ok", x)
			}),
			expected: 42,
		},
		{
			desc: "FilterOrElse (false)",
			result: result.FilterOrElse(result.Ok(2), func(x int) bool {
				return x > 40
			}, func(x int) error {
				return fmt.Errorf("filterOrElse: %d is not ok", x)
			}),
			expectedErr: errors.New("filterOrElse: 2 is not ok"),
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			assertEq(t, tC.result, tC.expected, tC.expectedErr)
		})
	}
}

func assertEq(t *testing.T, res fpgo.Result[int], expected int, expectedErr error) {
	x, err := res(context.TODO())
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
