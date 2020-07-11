package jolsat

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

type ib struct {
	Num int
	Ok  bool
}

type strs []string

func TestParseInt(t *testing.T) {
	n := NegativeFieldPrefix
	tests := map[string]struct {
		input strs
		want  ib
	}{
		"zero":                                 {strs{"", "0"}, ib{0, true}},
		"one":                                  {strs{"", "1"}, ib{1, true}},
		"nine":                                 {strs{"", "9"}, ib{9, true}},
		"ten":                                  {strs{"", "10"}, ib{10, true}},
		"negative zero":                        {strs{n, "0"}, ib{0, true}},
		"negative one":                         {strs{n, "1"}, ib{-1, true}},
		"negative nine":                        {strs{n, "9"}, ib{-9, true}},
		"negative ten":                         {strs{n, "10"}, ib{-10, true}},
		"empty slice":                          {strs{}, ib{0, false}},
		"too few elements":                     {strs{"5"}, ib{0, false}},
		"first element is meaningless":         {strs{"XX", "5"}, ib{0, false}},
		"positive; second element not an int":  {strs{"", "ten"}, ib{0, false}},
		"negative; second element not an int":  {strs{n, "ten"}, ib{0, false}},
		"positive; extra data in slice":        {strs{"", "5", "extra data"}, ib{5, true}},
		"negative; extra data in slice":        {strs{n, "5", "extra data"}, ib{-5, true}},
		"check that negative prefix isn't '-'": {strs{"-", "99"}, ib{0, false}},
	}
	var got ib
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			got.Num, got.Ok = parseInt(tt.input)
			diff := cmp.Diff(tt.want, got)
			if diff != "" {
				t.Errorf(diff)
			}
		})
	}
}

func TestProcess(t *testing.T) {
	a2e := []string{"a", "b", "c", "d", "e"}
	tests := map[string]struct {
		given []string
		want  []string
	}{
		// In bounds tests
		"2":     {a2e, strs{"b"}},
		"^3":    {a2e, strs{"c"}},
		"1-":    {a2e, a2e},
		"^3-":   {a2e, strs{"c", "d", "e"}},
		"-3":    {a2e, strs{"a", "b", "c"}},
		"-^2":   {a2e, strs{"a", "b", "c", "d"}},
		"2-3":   {a2e, strs{"b", "c"}},
		"2-^2":  {a2e, strs{"b", "c", "d"}},
		"^3-^2": {a2e, strs{"c", "d"}},
		// One or both ends out of bounds
		"1-6":   {a2e, a2e},
		"6":     {a2e, nil},
		"^6":    {a2e, nil},
		"^8-^6": {a2e, nil},
		"7-8":   {a2e, nil},
		"^8-10": {a2e, a2e},
		// Dodgy input
		"1": {nil, nil},
		"3": {strs{}, nil},
	}

	for spec, tt := range tests {
		t.Run(spec, func(t *testing.T) {
			fp, ok := NewFieldProcessor(spec)
			if !ok {
				t.Errorf(spec)
			}
			got := fp.process(tt.given)
			diff := cmp.Diff(tt.want, got)
			if diff != "" {
				t.Errorf(spec, diff)
			}
		})
	}
}

func TestNewFieldProcessor(t *testing.T) {

	tests := map[string]struct {
		fieldSpec string
		want      FieldProcessor
	}{
		"single positive field": {"2", FieldProcessor{first: 2, last: 2, statSpec: ""}},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			got, _ := NewFieldProcessor(tt.fieldSpec)
			diff := cmp.Diff(tt.want, got, cmp.AllowUnexported(FieldProcessor{}))
			if diff != "" {
				t.Errorf(diff)
			}
		})
	}

}

func TestFieldRangeToSliceIndices(t *testing.T) {
	var (
		fieldtests = []struct {
			spec               string
			length, start, end int
			ok                 bool
			description        string
		}{
			{"1-5", 5, 0, 5, true, "1 to 5 [1,2,3,4,5]"},
			{"1-^1", 5, 0, 5, true, "1 to -1 [1,2,3,4,5]"},
			{"2-3", 5, 1, 3, true, "2 to 3 [1,2,3,4,5]"},
			{"2-^2", 5, 1, 4, true, "2 to -2 [1,2,3,4,5]"},
			{"2-^4", 5, 1, 2, true, "2 to -4 [1,2,3,4,5]"},
			{"^3-^1", 5, 2, 5, true, "-3 to -1 [1,2,3,4,5]"},
			{"1-^10", 5, 0, 0, false, "1 to -10 [1,2,3,4,5]"},
			{"4-^4", 5, 0, 0, false, "4 to -4 [1,2,3,4,5]"},
			{"10-11", 5, 0, 0, false, "10 to 11 [1,2,3,4,5]"},
			{"^10-3", 5, 0, 3, true, "-10 to 3 [1,2,3,4,5]"},
			{"1-10", 5, 0, 5, true, "1 to 10 [1,2,3,4,5]"},
			{"2-3", 0, 0, 0, false, "2 to 3 []"},
		}
	)

	for _, tt := range fieldtests {
		t.Run(tt.description, func(t *testing.T) {
			fp, _ := NewFieldProcessor(tt.spec)
			start, end, ok := fp.getSliceIndices(tt.length)
			switch tt.ok {
			case true:
				if start != tt.start || end != tt.end || ok != tt.ok {
					t.Errorf("wanted: %v, %v, %v; got: %v, %v, %v, '%v'", tt.start, tt.end, tt.ok, start, end, ok, tt.description)
				}
			case false:
				if tt.ok != ok {
					t.Errorf("wanted: %v, %v, %v; got: %v, %v, %v, '%v'", tt.start, tt.end, tt.ok, start, end, ok, tt.description)
				}
			}
		})
	}
}
