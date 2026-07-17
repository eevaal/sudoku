package main

import (
	"fmt"
	"os"
	"os/user"
)

func main() {
	u, err := user.Current()
	if err != nil {
		fmt.Fprintf(os.Stderr, "whoami error: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(u.Username)
}
