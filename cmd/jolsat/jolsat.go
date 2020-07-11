package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"

	"github.com/jpluscplusm/jolsat/internal/jolsat"
)

var (
	p = fmt.Println
)

func main() {
	delimiter := flag.String("d", "\t", "Word delimiter")
	fieldFlag := flag.String("f", "1-", "Field list")
	flag.Parse()

	runtime.GOMAXPROCS(runtime.NumCPU())

	toProcessors, toDisplay := runProcessors(strings.Split(*fieldFlag, ","))
	go fanOut(os.Stdin, toProcessors, *delimiter)
	fanIn(os.Stdout, toDisplay, *delimiter)

}

func runProcessors(fields []string) (toProcessors, toDisplay []chan []string) {
	for _, f := range fields {
		fp, ok := jolsat.NewFieldProcessor(f)
		if ok {
			toFP := make(chan []string, 100)
			fromFP := make(chan []string, 100)
			toProcessors = append(toProcessors, toFP)
			toDisplay = append(toDisplay, fromFP)
			go fp.Run(toFP, fromFP)
		} else {
			panic("Couldn't make an FP with spec '" + f + "'")
		}
	}
	return
}

func fanOut(r io.Reader, receivers []chan []string, delimiter string) {
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
}

func fanIn(w io.Writer, senders []chan []string, delimiter string) {
	for {
		for i, sender := range senders {
			output, open := <-sender
			if !open {
				return
			}
			if i > 0 && len(output) > 0 {
				fmt.Fprint(w, delimiter)
			}
			fmt.Fprint(w, strings.Join(output, delimiter))
		}
		fmt.Fprintln(w)
	}
}
