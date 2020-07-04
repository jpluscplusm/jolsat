package main

import (
	"fmt"
	"io"
	"os"
)

func main() {
	if _, err := io.Copy(os.Stdout, os.Stdin); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
