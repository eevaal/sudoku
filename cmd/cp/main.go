package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}

	info, err := os.Stat(src)
	if err == nil {
		os.Chmod(dst, info.Mode())
	}
	return nil
}

func copyDir(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		relPath, _ := filepath.Rel(src, path)
		targetPath := filepath.Join(dst, relPath)

		if info.IsDir() {
			os.MkdirAll(targetPath, info.Mode())
		} else {
			copyFile(path, targetPath)
		}
		return nil
	})
}

func main() {
	args := os.Args[1:]
	recursive := false
	var paths []string

	for _, arg := range args {
		if strings.HasPrefix(arg, "-") {
			if strings.Contains(arg, "r") || strings.Contains(arg, "R") {
				recursive = true
			}
			continue
		}
		paths = append(paths, arg)
	}

	if len(paths) < 2 {
		fmt.Fprintln(os.Stderr, "usage: cp [-r] <src>... <dst>")
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
		fmt.Fprintf(os.Stderr, "cp: target '%s' is not a directory\n", dst)
		os.Exit(1)
	}

	for _, src := range actualSrcs {
		srcInfo, err := os.Stat(src)
		if err != nil {
			fmt.Fprintf(os.Stderr, "cp: cannot stat '%s': %v\n", src, err)
			continue
		}

		targetPath := dst
		if dstIsDir {
			targetPath = filepath.Join(dst, filepath.Base(src))
		}

		if srcInfo.IsDir() {
			if !recursive {
				fmt.Fprintf(os.Stderr, "cp: -r not specified; omitting directory '%s'\n", src)
				continue
			}
			err = copyDir(src, targetPath)
			if err != nil {
				fmt.Fprintf(os.Stderr, "cp: error copying directory '%s': %v\n", src, err)
			}
		} else {
			err = copyFile(src, targetPath)
			if err != nil {
				fmt.Fprintf(os.Stderr, "cp: error copying file '%s': %v\n", src, err)
			}
		}
	}
}
