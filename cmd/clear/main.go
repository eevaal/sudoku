package main

import (
	"fmt"
)

func main() {
	// ANSI escape code to clear screen and scrollback buffer
	// ESC [ H moves cursor to top left
	// ESC [ 2 J clears the screen
	// ESC [ 3 J clears the scrollback buffer
	fmt.Print("\033[H\033[2J\033[3J")
}
