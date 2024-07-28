package io_test

import (
	"testing"

	"github.com/onur1/gofp"
	"github.com/onur1/gofp/io"
	"github.com/stretchr/testify/assert"
)

func TestIO(t *testing.T) {
	testCases := []struct {
		desc     string
		io       gofp.IO[int]
		expected int
	}{
		{
			desc:     "Nil",
			io:       io.Of(2),
			expected: 2,
		},
		{
			desc:     "Map",
			io:       io.Map(io.Of(2), double),
			expected: 4,
		},
		{
			desc:     "Ap",
			io:       io.Ap(io.Of(double), io.Of(2)),
			expected: 4,
		},
		{
			desc:     "ApFirst",
			io:       io.ApFirst(io.Of(1), io.Of(2)),
			expected: 1,
		},
		{
			desc:     "ApSecond",
			io:       io.ApSecond(io.Of(1), io.Of(2)),
			expected: 2,
		},
		{
			desc: "Chain",
			io: io.Chain(io.Of(2), func(a int) gofp.IO[int] {
				return io.Of(double(a))
			}),
			expected: 4,
		},
		{
			desc: "ChainFirst",
			io: io.ChainFirst(io.Of(2), func(a int) gofp.IO[int] {
				return io.Of(double(a))
			}),
			expected: 2,
		},
		{
			desc: "ChainRec",
			io: io.ChainRec(0, func(n int) gofp.IO[func() (int, int, bool)] {
				return io.Of(
					func() (int, int, bool) {
						if n < 15000 {
							return n + 1, 0, false
						} else {
							return 0, 15000, true
						}
					},
				)
			}),
			expected: 15000,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			assert.Equal(t, tC.expected, tC.io())
		})
	}
}

func double(n int) int {
	return n * 2
}
