package main

import (
	"fmt"
	"os"
	"os/user"
	"runtime"
	"syscall"
	"time"
	"unsafe"
)

var (
	kernel32           = syscall.NewLazyDLL("kernel32.dll")
	getTickCount64     = kernel32.NewProc("GetTickCount64")
	globalMemoryStatus = kernel32.NewProc("GlobalMemoryStatusEx")
)

type memoryStatusEx struct {
	cbSize                  uint32
	dwMemoryLoad            uint32
	ullTotalPhys            uint64
	ullAvailPhys            uint64
	ullTotalPageFile        uint64
	ullAvailPageFile        uint64
	ullTotalVirtual         uint64
	ullAvailVirtual         uint64
	ullAvailExtendedVirtual uint64
}

func main() {
	// ASCII Art - Colored Clippy
	// Colors used: 
	// \033[38;5;250m - Silver/Grey (Body)
	// \033[38;5;16m  - Black (Eyebrows/Pupils)
	// \033[38;5;231m - White (Eyes)
	cBody := "\033[38;5;250m"
	cBlack := "\033[38;5;16m"
	cWhite := "\033[38;5;231m"
	cReset := "\033[0m"

	ascii := []string{
		cBody + "⠀⠀⠀⠀⠀⣠⢖⣭⣿⣿⣷⣄⠀⠀⠀⠀" + cReset,
		cBody + "⠀⠀⠀⢀⣠⣡⣟⠁⠀⠀⠹⣿⡇⠀⠀⠀" + cReset,
		cBody + "⠀⢠⣾" + cBlack + "⣿⡟⣿⠿⠃" + cBody + "⠀⠀⢸⣿⣧⣄⠀⠀" + cReset,
		cBody + "⢠⡮⢁" + cWhite + "⣤" + cBlack + "⣤" + cBody + "⡉⠳⡄⠀⣠⠾" + cWhite + "⠹" + cBlack + "⠿" + cBody + "⣻⣷⡄" + cReset,
		cBody + "⠸⣇" + cWhite + "⠼⠿⡿⠟" + cBlack + "⣤" + cBody + "⡇⣮⠀" + cWhite + "⣶" + cBlack + "⣶" + cBody + "⣦⡈⢻⠋" + cReset,
		cBody + "⠀⠈⠙⢳⣶⡞⠋⠀⠹⢦⡙⠛⠛⣠⡾⠀" + cReset,
		cBody + "⠀⠀⠀⢸⣿⣷⣶⡀⠀⠀⣿⣿⢉⣿⡆⠀" + cReset,
		cBody + "⠀⠀⠀⢸⣿⣿⣿⠀⠀⠀⣿⡇⣾⣿⠃⠀" + cReset,
		cBody + "⠀⠀⠀⢰⣿⣿⣿⡀⠀⠀⣿⡇⣿⣿⠀⠀" + cReset,
		cBody + "⠀⠀⠀⠀⣿⢻⣿⢇⠀⠀⣿⣿⣿⣿⠀⠀" + cReset,
		cBody + "⠀⠀⠀⠀⢿⣾⡿⣾⣲⣚⣽⠇⣿⣿⠀⠀" + cReset,
		cBody + "⠀⠀⠀⠀⠸⣧⢧⠈⠉⠉⠁⠀⣿⣿⠀⠀" + cReset,
		cBody + "⠀⠀⠀⠀⠀⢻⣞⢆⠀⠀⠀⡠⢻⡿⠀⠀" + cReset,
		cBody + "⠀⠀⠀⠀⠀⠀⠙⢷⣭⣤⣭⡴⠟⠀⠀⠀" + cReset,
	}

	// Fetch info
	u, _ := user.Current()
	hostname, _ := os.Hostname()

	// Uptime
	r, _, _ := getTickCount64.Call()
	uptimeDuration := time.Duration(r) * time.Millisecond
	uptimeStr := fmt.Sprintf("%d mins", int(uptimeDuration.Minutes()))
	if uptimeDuration.Hours() > 1 {
		uptimeStr = fmt.Sprintf("%d hours, %d mins", int(uptimeDuration.Hours()), int(uptimeDuration.Minutes())%60)
	}

	// Memory
	var mem memoryStatusEx
	mem.cbSize = uint32(unsafe.Sizeof(mem))
	globalMemoryStatus.Call(uintptr(unsafe.Pointer(&mem)))
	totalMemMB := mem.ullTotalPhys / 1024 / 1024
	usedMemMB := (mem.ullTotalPhys - mem.ullAvailPhys) / 1024 / 1024

	info := []string{
		fmt.Sprintf("\033[36m%s\033[0m@\033[36m%s\033[0m", u.Username, hostname),
		"-------------------------",
		fmt.Sprintf("\033[36mOS\033[0m: Windows %s", runtime.GOARCH),
		fmt.Sprintf("\033[36mUptime\033[0m: %s", uptimeStr),
		fmt.Sprintf("\033[36mMemory\033[0m: %d MB / %d MB", usedMemMB, totalMemMB),
		fmt.Sprintf("\033[36mShell\033[0m: sudoku"),
	}

	// Pad ascii or info so they match
	lines := len(ascii)
	if len(info) > lines {
		lines = len(info)
	}

	fmt.Println()
	for i := 0; i < lines; i++ {
		a := ""
		if i < len(ascii) {
			a = ascii[i]
		}
		// Calculate display length of ascii (strip ansi codes for padding calc)
		padLength := 40
		
		in := ""
		if i < len(info) {
			in = info[i]
		}

		fmt.Printf("  %-*s %s\n", padLength, a, in)
	}
	fmt.Println()
}
