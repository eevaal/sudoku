package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// fallbackMove handles moves across different drives/volumes
func fallbackMove(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}

	out, err := os.Create(dst)
	if err != nil {
		in.Close()
		return err
	}

	_, err = io.Copy(out, in)
	in.Close()
	out.Close()
	if err != nil {
		return err
	}

	return os.RemoveAll(src)
}

func main() {
	args := os.Args[1:]
	var paths []string

	for _, arg := range args {
		if strings.HasPrefix(arg, "-") {
			// Ignore flags for mv for simplicity
			continue
		}
		paths = append(paths, arg)
	}

	if len(paths) < 2 {
		fmt.Fprintln(os.Stderr, "usage: mv <src>... <dst>")
		os.Exit(1)
	}

	dst := paths[len(paths)-1]
	srcs := paths[:len(paths)-1]

	var actualSrcs []string
	for _, src := range srcs {
		matches, err := filepath.Glob(src)
		if err == nil && len(matches) > 0 {
			actualSrcs = append(actualSrcs, matches...)
		} else {
			actualSrcs = append(actualSrcs, src)
		}
	}

	dstInfo, err := os.Stat(dst)
	dstIsDir := err == nil && dstInfo.IsDir()

	if len(actualSrcs) > 1 && !dstIsDir {
		fmt.Fprintf(os.Stderr, "mv: target '%s' is not a directory\n", dst)
		os.Exit(1)
	}

	for _, src := range actualSrcs {
		targetPath := dst
		if dstIsDir {
			targetPath = filepath.Join(dst, filepath.Base(src))
		}

		err := os.Rename(src, targetPath)
		if err != nil {
			err = fallbackMove(src, targetPath)
			if err != nil {
				fmt.Fprintf(os.Stderr, "mv: cannot move '%s' to '%s': %v\n", src, targetPath, err)
			}
		}
	}
}
