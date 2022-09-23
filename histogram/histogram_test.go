package histogram

import (
	"fmt"
	"testing"
)

func TestHistogram(t *testing.T) {
	tt := []struct {
		min      int
		max      int
		buckets  int
		expected [][]int
	}{
		{
			min:      0,
			max:      5,
			buckets:  3,
			expected: [][]int{{0}, {1}, {2, 3}, {4, 5}},
		},
		{
			min:      0,
			max:      10,
			buckets:  5,
			expected: [][]int{{0}, {1, 2}, {3, 4}, {5, 6}, {7, 8}, {9, 10}},
		},
		{
			min:      0,
			max:      100,
			buckets:  4,
			expected: [][]int{{0}, seq(1, 25), seq(26, 50), seq(51, 75), seq(76, 100)},
		},
	}

	for _, tc := range tt {
		fit := fitter(tc.min, tc.max, tc.buckets)

		t.Run(fmt.Sprintf("fitter(%d,%d,%d)", tc.min, tc.max, tc.buckets), func(t *testing.T) {
			for expected, nums := range tc.expected {
				for _, num := range nums {
					actual := fit(num)
					if actual != expected {
						t.Errorf(
							"fit(%d) should eq %d, got %d",
							num,
							expected,
							actual,
						)
					}
				}
			}
		})
	}
}

// returns ints in the interval [start, stop]
func seq(start, stop int) []int {
	out := make([]int, stop-start+1)

	for i := range out {
		out[i] = i + start
	}

	return out
}
