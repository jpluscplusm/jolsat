package jolsat

import "testing"

func TestFieldRangeToSliceIndices(t *testing.T) {
	// wrapping these in nicer types results in *not* nicer code: https://groups.google.com/g/golang-nuts/c/hvraQJtcY9U/m/Lh8YQch2UiQJ
	var (
		fieldtests = []struct {
			first       int
			last        int
			length      int
			start       int
			end         int
			ok          bool
			description string
		}{
			{1, 5, 5, 0, 5, true, "1 to 5 [1,2,3,4,5]"},
			{1, -1, 5, 0, 5, true, "1 to -1 [1,2,3,4,5]"},
			{2, 3, 5, 1, 3, true, "2 to 3 [1,2,3,4,5]"},
			{2, -2, 5, 1, 4, true, "2 to -2 [1,2,3,4,5]"},
			{2, -4, 5, 1, 2, true, "2 to -4 [1,2,3,4,5]"},
			{-3, -1, 5, 2, 5, true, "-3 to -1 [1,2,3,4,5]"},
			{1, -10, 5, 0, 0, false, "1 to -10 [1,2,3,4,5]"},
			{4, -4, 5, 0, 0, false, "4 to -4 [1,2,3,4,5]"},
			{10, 11, 5, 0, 0, false, "10 to 11 [1,2,3,4,5]"},
			{-10, 3, 5, 0, 3, true, "-10 to 3 [1,2,3,4,5]"},
			{1, 10, 5, 0, 5, true, "1 to 10 [1,2,3,4,5]"},
			{2, 3, 0, 0, 0, false, "2 to 3 []"},
			{0, 2, 5, 0, 0, false, "0 to 2 [1,2,3,4,5]"},
			{2, 0, 5, 0, 0, false, "2 to 0 [1,2,3,4,5]"},
		}

		start, end int
		ok         bool
	)

	for _, tt := range fieldtests {
		t.Run(tt.description, func(t *testing.T) {
			start, end, ok = FieldRangeToSliceIndices(tt.first, tt.last, tt.length)
			if start != tt.start || end != tt.end || ok != tt.ok {
				t.Errorf("wanted: %v, %v, %v; got: %v, %v, %v, '%v'", start, end, ok, tt.start, tt.end, tt.ok, tt.description)
			}
		})
	}
}
