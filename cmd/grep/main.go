package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

func main() {
	args := os.Args[1:]
	ignoreCase := false
	var nonFlags []string

	for _, arg := range args {
		if arg == "-i" {
			ignoreCase = true
			continue
		}
		if strings.HasPrefix(arg, "-") {
			continue
		}
		nonFlags = append(nonFlags, arg)
	}

	if len(nonFlags) == 0 {
		fmt.Fprintln(os.Stderr, "usage: grep [-i] <pattern> [file...]")
		os.Exit(1)
	}

	patternStr := nonFlags[0]
	if ignoreCase {
		patternStr = "(?i)" + patternStr
	}

	pattern, err := regexp.Compile(patternStr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "grep: invalid regex: %v\n", err)
		os.Exit(1)
	}

	files := nonFlags[1:]

	if len(files) == 0 {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			line := scanner.Text()
			if pattern.MatchString(line) {
				fmt.Println(line)
			}
		}
		return
	}

	for _, file := range files {
		f, err := os.Open(file)
		if err != nil {
			fmt.Fprintf(os.Stderr, "grep: %s: %v\n", file, err)
			continue
		}
		
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			line := scanner.Text()
			if pattern.MatchString(line) {
				if len(files) > 1 {
					fmt.Printf("%s:%s\n", file, line)
				} else {
					fmt.Println(line)
				}
			}
		}
		f.Close()
	}
}
