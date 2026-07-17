package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	args := os.Args[1:]
	var paths []string

	for _, arg := range args {
		if strings.HasPrefix(arg, "-") {
			// ignore flags like -p since we'll always act like -p
			continue
		}
		paths = append(paths, arg)
	}

	if len(paths) == 0 {
		fmt.Fprintln(os.Stderr, "usage: mkdir <dir>...")
		os.Exit(1)
	}

	for _, p := range paths {
		err := os.MkdirAll(p, 0755)
		if err != nil {
			fmt.Fprintf(os.Stderr, "mkdir: cannot create directory '%s': %v\n", p, err)
		}
	}
}
