package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

type FieldProcessor struct {
	Run    func(in, out chan []string)
	Input  chan []string
	Output chan []string
}

var (
	integerOnly          = regexp.MustCompile("^([0-9]+)$")
	dashInteger          = regexp.MustCompile("^-([0-9]+)$")
	integerDash          = regexp.MustCompile("^([0-9]+)-$")
	integerDashInteger   = regexp.MustCompile("^([0-9]+)-([0-9]+)$")
	integerColonAnything = regexp.MustCompile("^([0-9]+):")
	p                    = fmt.Println
	I                    = func(s string) int { a, _ := strconv.Atoi(s); return a }
)

// explodeFieldRange returns the 2 ints representing the requested field range
// Permitted patterns for field:
//  Type A: exactly analogous to `cut (1)`:
//   integer  (e.g. "4")          meaning "field 4 only"
//   -integer (e.g. "-3")         meaning "the first field of the line to field 3, inclusive"
//   integer- (e.g. "5-")         meaning "field 5 to the last field of the line, inclusive"
//   integer-integer (e.g. "3-5") meaning "field 3 to field 5, inclusive"
//  Type B: jolsat-specific:
//   integer:stat
//   integer:stat[option]
//   integer:stat[option1+option2]
//   integer:stat[param1:setting1+option1+param2:setting2]
func explodeFieldRange(fieldRange string) (first, last int) {

	switch {
	case integerOnly.MatchString(fieldRange):
		int1 := I(integerOnly.FindStringSubmatch(fieldRange)[1])
		first, last = int1, int1
	case dashInteger.MatchString(fieldRange):
		int1 := I(dashInteger.FindStringSubmatch(fieldRange)[1])
		first, last = 1, int1
	case integerDash.MatchString(fieldRange):
		int1 := I(integerDash.FindStringSubmatch(fieldRange)[1])
		first, last = int1, 0
	case integerDashInteger.MatchString(fieldRange):
		first = I(integerDashInteger.FindStringSubmatch(fieldRange)[1])
		last = I(integerDashInteger.FindStringSubmatch(fieldRange)[2])
	case integerColonAnything.MatchString(fieldRange):
		int1 := I(integerColonAnything.FindStringSubmatch(fieldRange)[1])
		first, last = int1, int1
	default:
		panic(fieldRange)
	}
	return
}

func main() {
	var (
		wg         = new(sync.WaitGroup)
		processors []FieldProcessor
	)

	delimiter := flag.String("d", "\t", "Word delimiter")
	fieldFlag := flag.String("f", "1-", "Field list")
	channelBufferSize := flag.Int("b", 99, "Buffer size")
	flag.Parse()

	fieldList := strings.Split(*fieldFlag, ",")
	fields := make([]string, len(fieldList))

	for i, s := range fieldList {
		fields[i] = s
	}

	for i, f := range fields {
		first, last := explodeFieldRange(f)
		processors = append(processors, FieldProcessor{
			Input:  make(chan []string, *channelBufferSize),
			Output: make(chan []string, *channelBufferSize),
		})
		processors[i].Run = func(in, out chan []string) {

			emptySlice := make([]string, 0)

			for x := range in {
				var start, end, lenx int
				lenx = len(x)
				start = first - 1
				// Quick escapes: empty input or not enough input for our range
				switch {
				case lenx == 0, first > lenx:
					out <- emptySlice
					continue
				}
				if last == 0 {
					end = lenx - 1
				} else {
					end = last - 1
				}
				if end > lenx-1 {
					end = lenx - 1
				}
				//p(start, end, lenx)
				out <- x[start : end+1]
			}
			close(out)
		}
		go processors[i].Run(processors[i].Input, processors[i].Output)
	}

	wg.Add(1)
	go func(fp []FieldProcessor, allDone *sync.WaitGroup) {
		for {
			for i, f := range fp {
				output, open := <-f.Output
				if !open {
					allDone.Done()
					return
				}
				if i > 0 && len(output) > 0 {
					fmt.Print(*delimiter)
				}
				fmt.Print(strings.Join(output, *delimiter))
			}
			fmt.Print("\n")
		}
	}(processors, wg)

	scanner := bufio.NewScanner(os.Stdin)

	var tokens []string
	for scanner.Scan() {
		tokens = strings.Split(scanner.Text(), *delimiter)
		for _, processor := range processors {
			processor.Input <- tokens
		}
	}

	for _, processor := range processors {
		close(processor.Input)
	}

	wg.Wait()
}
