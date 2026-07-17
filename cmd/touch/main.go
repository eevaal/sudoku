package main

import (
	"fmt"
	"os"
	"time"
)

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "usage: touch <file>...")
		os.Exit(1)
	}

	currentTime := time.Now()
	for _, file := range args {
		_, err := os.Stat(file)
		if os.IsNotExist(err) {
			// Create empty file
			f, err := os.Create(file)
			if err != nil {
				fmt.Fprintf(os.Stderr, "touch: cannot create '%s': %v\n", file, err)
				continue
			}
			f.Close()
		} else if err == nil {
			// Update timestamps
			err = os.Chtimes(file, currentTime, currentTime)
			if err != nil {
				fmt.Fprintf(os.Stderr, "touch: cannot touch '%s': %v\n", file, err)
			}
		} else {
			fmt.Fprintf(os.Stderr, "touch: cannot stat '%s': %v\n", file, err)
		}
	}
}
