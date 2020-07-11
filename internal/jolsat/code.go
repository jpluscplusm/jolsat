package jolsat

import (
	"fmt"
	"regexp"
	"strconv"
)

const (
	NegativeFieldPrefix = "^" // don't set this to "-"!
)

type FieldProcessor struct {
	first, last int
	statSpec    string
}

func (fp FieldProcessor) Run(fromInput <-chan []string, toOutput chan<- []string) {
	defer close(toOutput)
	for s := range fromInput {
		toOutput <- fp.process(s)
	}
	return
}

func (fp FieldProcessor) process(input []string) (output []string) {
	start, end, ok := fp.getSliceIndices(len(input))
	if !ok {
		return
	}
	output = input[start:end]
	return output
}

func (fp FieldProcessor) String() string {
	return fmt.Sprintf("%v-%v:%s", fp.first, fp.last, fp.statSpec)
}

func NewFieldProcessor(fieldSpec string) (fp FieldProcessor, ok bool) {
	fp.first, fp.last, fp.statSpec, ok = parseFieldSpec(fieldSpec)
	return
}

var (
	prefixRegex    = regexp.QuoteMeta(NegativeFieldPrefix)
	singleField    = regexp.MustCompile("^(" + prefixRegex + "?)([[:digit:]]+)$")
	openLeftRange  = regexp.MustCompile("^-(" + prefixRegex + "?)([[:digit:]]+)$")
	openRightRange = regexp.MustCompile("^(" + prefixRegex + "?)([[:digit:]]+)-$")
	specificRange  = regexp.MustCompile("^(" + prefixRegex + "?)([[:digit:]]+)-(" + prefixRegex + "?)([[:digit:]]+)$")
	statsField     = regexp.MustCompile("^(" + prefixRegex + "?)([[:digit:]]+):(.*)$")
)

// explodeFieldRange returns the 2 ints representing the requested field range
// Permitted patterns for field:
//  Type A: exactly analogous to `cut (1)`:
//   integer		(e.g. "4")	== "field 4 only"
//   -integer		(e.g. "-3")	== "the first field of the line to field 3, inclusive"
//   integer-		(e.g. "5-")	== "field 5 to the last field of the line, inclusive"
//   integer-integer	(e.g. "3-5")	== "field 3 to field 5, inclusive"
//  Type B: jolsat-specific ranges:
//   ^integer		(e.g. "^2")	== "field 2, counting back from the end of the line, only"
//   -^integer		(e.g. "-^3")	== "the first field of the line to field 3, counting back from the end of the line, inclusive"
//   ^integer-		(e.g. "^5-")	== "field 5, counting back from the end of the line, to the last field of the line, inclusive"
//   ^integer-^integer	(e.g. "^3-^5")	== "field 3 to field 5, both fields counting back from the end of the line, inclusive"
//  Type C: jolsat-specific stats
//   integer:stat
//   integer:stat[option]
//   integer:stat[option1+option2]
//   integer:stat[param1:setting1+option1+param2:setting2]
func parseFieldSpec(fieldSpec string) (first, last int, statSpec string, ok bool) {
	switch {
	case singleField.MatchString(fieldSpec):
		first, ok = parseInt(singleField.FindStringSubmatch(fieldSpec)[1:3])
		last = first
	case openLeftRange.MatchString(fieldSpec):
		first = 1
		last, ok = parseInt(openLeftRange.FindStringSubmatch(fieldSpec)[1:3])
	case openRightRange.MatchString(fieldSpec):
		first, ok = parseInt(openRightRange.FindStringSubmatch(fieldSpec)[1:3])
		last = -1
	case specificRange.MatchString(fieldSpec):
		var ok1, ok2 bool
		match := specificRange.FindStringSubmatch(fieldSpec)
		first, ok1 = parseInt(match[1:3])
		last, ok2 = parseInt(match[3:5])
		ok = ok1 && ok2
	case statsField.MatchString(fieldSpec):
		match := statsField.FindStringSubmatch(fieldSpec)
		first, ok = parseInt(match[1:3])
		last = first
		statSpec = match[3]
	}
	if first == 0 || last == 0 { // zeros are meaningless
		ok = false
	}
	return
}

func parseInt(s []string) (int, bool) {
	if len(s) < 2 {
		return 0, false
	}
	num, err := strconv.Atoi(s[1])
	if err != nil {
		return 0, false
	}
	switch s[0] {
	case "": // positive number
	case NegativeFieldPrefix:
		num = -num
	default: // we *shouldn't* have reached this point (i.e. a regexp matched something ... impossible?)
		return 0, false
	}
	return num, true
}

func (fp FieldProcessor) getSliceIndices(length int) (start int, end int, ok bool) {
	// Quick escapes we don't need to spend time processing
	if length < 1 { // empty/fubar slice
		return
	}

	// Translate negative (countback) fields into positive fields
	start = fp.first
	if start < 0 {
		start += length + 1
	}
	end = fp.last
	if end < 0 {
		end += length + 1
	}

	switch {
	case end < 1, end < start, start > length: // Situations we can't process
	default:
		ok = true
		if start < 1 {
			start = 1
		}
		if end > length {
			end = length
		}
		// Turn field references into slice indices
		start--
	}

	return
}
