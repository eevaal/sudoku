package main

import (
	"fmt"
	"io"
	"os"
)

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		// Read from stdin
		_, _ = io.Copy(os.Stdout, os.Stdin)
		return
	}

	for _, file := range args {
		f, err := os.Open(file)
		if err != nil {
			fmt.Fprintf(os.Stderr, "cat: %s: %v\n", file, err)
			continue
		}
		_, err = io.Copy(os.Stdout, f)
		if err != nil {
			fmt.Fprintf(os.Stderr, "cat: read error: %v\n", err)
		}
		f.Close()
	}
}
