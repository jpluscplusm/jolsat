package jolsat

func FieldRangeToSliceIndices(first int, last int, length int) (start int, end int, ok bool) {
	// Quick escapes we don't need to spend time processing
	if length < 1 { // empty/fubar slice
		return 0, 0, false
	} else if first == 0 { // 0 is meaningless as a field
		return 0, 0, false
	} else if last == 0 { // 0 is meaningless as a field
		return 0, 0, false
	}

	// Translate negative (countback) fields into positive fields
	start = first
	if start < 0 {
		start += length + 1
	}
	end = last
	if end < 0 {
		end += length + 1
	}

	// Situations we don't need to process
	if end < 1 { // we can't end before field 1
		return 0, 0, false
	} else if end < start { // this function doesn't deal with possible requests to reverse slice input
		return 0, 0, false
	} else if start > length { // we can't start after the end of the slice
		return 0, 0, false
	}

	// Fixups
	if start < 1 {
		start = 1
	}
	if end > length {
		end = length
	}

	// Turn field references into slice indices
	start -= 1

	return start, end, true
}
