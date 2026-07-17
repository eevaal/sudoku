package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func main() {
	args := os.Args[1:]
	lines := 10
	var files []string

	for i := 0; i < len(args); i++ {
		if args[i] == "-n" && i+1 < len(args) {
			n, err := strconv.Atoi(args[i+1])
			if err == nil {
				lines = n
			}
			i++
			continue
		}
		files = append(files, args[i])
	}

	if len(files) == 0 {
		processTail(os.Stdin, lines)
		return
	}

	for i, file := range files {
		f, err := os.Open(file)
		if err != nil {
			fmt.Fprintf(os.Stderr, "tail: %s: %v\n", file, err)
			continue
		}
		if len(files) > 1 {
			if i > 0 {
				fmt.Println()
			}
			fmt.Printf("==> %s <==\n", file)
		}
		processTail(f, lines)
		f.Close()
	}
}

func processTail(f *os.File, lines int) {
	scanner := bufio.NewScanner(f)
	var buffer []string

	for scanner.Scan() {
		buffer = append(buffer, scanner.Text())
		if len(buffer) > lines {
			buffer = buffer[1:]
		}
	}
	
	for _, line := range buffer {
		fmt.Println(line)
	}
}
