package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	args := os.Args[1:]
	dir := "."
	namePattern := ""

	for i := 0; i < len(args); i++ {
		if args[i] == "-name" && i+1 < len(args) {
			namePattern = args[i+1]
			i++
			continue
		}
		if dir == "." && !filepath.IsAbs(args[i]) && !strings.HasPrefix(args[i], "-") {
			dir = args[i]
		}
	}

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Fprintf(os.Stderr, "find: %s: %v\n", path, err)
			return nil
		}

		if namePattern != "" {
			matched, err := filepath.Match(namePattern, info.Name())
			if err == nil && matched {
				fmt.Println(path)
			}
		} else {
			fmt.Println(path)
		}
		return nil
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, "find error: %v\n", err)
	}
}
