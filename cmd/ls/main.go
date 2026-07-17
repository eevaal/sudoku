package main

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
	"time"
)

func main() {
	args := os.Args[1:]
	showAll := false
	longFormat := false
	var paths []string

	for _, arg := range args {
		if strings.HasPrefix(arg, "-") {
			if strings.Contains(arg, "a") {
				showAll = true
			}
			if strings.Contains(arg, "l") {
				longFormat = true
			}
			continue
		}
		paths = append(paths, arg)
	}

	if len(paths) == 0 {
		paths = append(paths, ".")
	}

	for i, path := range paths {
		if len(paths) > 1 {
			fmt.Printf("%s:\n", path)
		}

		entries, err := os.ReadDir(path)
		if err != nil {
			// check if it's a file
			info, err2 := os.Stat(path)
			if err2 == nil && !info.IsDir() {
				printEntry(info, longFormat)
				continue
			}
			fmt.Fprintf(os.Stderr, "ls: cannot access '%s': %v\n", path, err)
			continue
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		for _, entry := range entries {
			if !showAll && strings.HasPrefix(entry.Name(), ".") {
				continue
			}
			info, err := entry.Info()
			if err != nil {
				continue
			}
			if longFormat {
				modTime := info.ModTime().Format(time.Stamp)
				size := info.Size()
				mode := info.Mode().String()
				name := formatName(info)
				fmt.Fprintf(w, "%s\t%d\t%s\t%s\n", mode, size, modTime, name)
			} else {
				fmt.Fprintf(w, "%s\t", formatName(info))
			}
		}
		if !longFormat {
			fmt.Fprintln(w)
		}
		w.Flush()

		if i < len(paths)-1 {
			fmt.Println()
		}
	}
}

func formatName(info os.FileInfo) string {
	name := info.Name()
	if info.IsDir() {
		return name + "/"
	}
	if info.Mode()&0111 != 0 {
		return name + "*"
	}
	return name
}

func printEntry(info os.FileInfo, longFormat bool) {
	if longFormat {
		modTime := info.ModTime().Format(time.Stamp)
		fmt.Printf("%s\t%d\t%s\t%s\n", info.Mode().String(), info.Size(), modTime, formatName(info))
	} else {
		fmt.Println(formatName(info))
	}
}
