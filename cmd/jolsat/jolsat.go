package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/jpluscplusm/jolsat/internal/jolsat"
)

func main() {
	delimiter := flag.String("d", "\t", "Word delimiter")
	fieldFlag := flag.String("f", "1-", "Field list")
	flag.Parse()

	var (
		fields = strings.Split(*fieldFlag, ",")
		fanOut []chan []string
		fanIn  []chan []string
		//p      = fmt.Println
	)

	for _, f := range fields {
		fp, ok := jolsat.NewFieldProcessor(f)
		switch ok {
		case true:
			toFP := make(chan []string, 100)
			fromFP := make(chan []string, 100)
			fanIn = append(fanIn, fromFP)
			fanOut = append(fanOut, toFP)
			go fp.Run(toFP, fromFP)
		case false:
			panic("Couldn't make an FP with spec '" + f + "'")
		}
	}

	go func(r io.Reader, receivers []chan []string, delimiter string) {
		scanner := bufio.NewScanner(r)

		for scanner.Scan() {
			tokens := strings.Split(scanner.Text(), delimiter)
			for _, r := range receivers {
				r <- tokens
			}
		}
		for _, r := range receivers {
			close(r)
		}
	}(os.Stdin, fanOut, *delimiter)

	var output []string
	for open := true; open; {
		for i, s := range fanIn {
			output, open = <-s
			if !open {
				break
			}
			if i > 0 && len(output) > 0 {
				fmt.Print(*delimiter)
			}
			fmt.Print(strings.Join(output, *delimiter))
		}
		if open {
			fmt.Print("\n")
		}
	}
}
