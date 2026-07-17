package main

import (
	"fmt"
	"os"
)

func main() {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "pwd error: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(dir)
}
