package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: sleep <duration>")
		os.Exit(1)
	}

	arg := os.Args[1]
	
	if !strings.HasSuffix(arg, "s") && !strings.HasSuffix(arg, "m") && !strings.HasSuffix(arg, "h") {
		arg += "s"
	}

	d, err := time.ParseDuration(arg)
	if err != nil {
		// fallback to float seconds
		arg = strings.TrimSuffix(arg, "s")
		f, err2 := strconv.ParseFloat(arg, 64)
		if err2 != nil {
			fmt.Fprintf(os.Stderr, "sleep: invalid duration '%s'\n", os.Args[1])
			os.Exit(1)
		}
		d = time.Duration(f * float64(time.Second))
	}

	time.Sleep(d)
}
