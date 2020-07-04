package main

import (
	"flag"
	"fmt"
	"io"
	"os"
)

func main() {
	delimiter := flag.String("d", "\t", "Word delimiter")
	//fieldSpec := flag.String("f", "1-", "Field list")
	flag.Parse()
	if _, err := io.Copy(os.Stdout, os.Stdin); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Fprintln(os.Stderr, delimiter)
	//fmt.Fprintln(os.Stderr, fieldSpec)
}
