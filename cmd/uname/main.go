package main

import (
	"fmt"
	"os"
	"runtime"
)

func main() {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}
	// Mimic 'uname -a' format
	fmt.Printf("Windows %s %s %s\n", hostname, runtime.Version(), runtime.GOARCH)
}
