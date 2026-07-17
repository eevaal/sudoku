package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	args := os.Args[1:]
	showLines, showWords, showBytes := false, false, false
	var files []string

	for _, arg := range args {
		if strings.HasPrefix(arg, "-") {
			if strings.Contains(arg, "l") {
				showLines = true
			}
			if strings.Contains(arg, "w") {
				showWords = true
			}
			if strings.Contains(arg, "c") {
				showBytes = true
			}
			continue
		}
		files = append(files, arg)
	}

	if !showLines && !showWords && !showBytes {
		showLines, showWords, showBytes = true, true, true
	}

	if len(files) == 0 {
		lines, words, bytes := count(os.Stdin)
		printCounts(lines, words, bytes, "", showLines, showWords, showBytes)
		return
	}

	var totalLines, totalWords, totalBytes int
	for _, file := range files {
		f, err := os.Open(file)
		if err != nil {
			fmt.Fprintf(os.Stderr, "wc: %s: %v\n", file, err)
			continue
		}
		l, w, c := count(f)
		f.Close()
		printCounts(l, w, c, file, showLines, showWords, showBytes)
		totalLines += l
		totalWords += w
		totalBytes += c
	}

	if len(files) > 1 {
		printCounts(totalLines, totalWords, totalBytes, "total", showLines, showWords, showBytes)
	}
}

func count(r io.Reader) (lines, words, bytes int) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		lines++
		text := scanner.Text()
		bytes += len(text) + 2 // approx for CRLF in windows
		words += len(strings.Fields(text))
	}
	return
}

func printCounts(l, w, c int, name string, showL, showW, showC bool) {
	var parts []string
	if showL {
		parts = append(parts, fmt.Sprintf("%8d", l))
	}
	if showW {
		parts = append(parts, fmt.Sprintf("%8d", w))
	}
	if showC {
		parts = append(parts, fmt.Sprintf("%8d", c))
	}
	if name != "" {
		parts = append(parts, name)
	}
	fmt.Println(strings.Join(parts, " "))
}
