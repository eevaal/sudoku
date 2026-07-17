package main

import (
	"fmt"
	"syscall"
	"time"
)

var (
	kernel32       = syscall.NewLazyDLL("kernel32.dll")
	getTickCount64 = kernel32.NewProc("GetTickCount64")
)

func main() {
	r, _, _ := getTickCount64.Call()
	d := time.Duration(r) * time.Millisecond
	
	now := time.Now().Format("15:04:05")
	
	days := int(d.Hours()) / 24
	hours := int(d.Hours()) % 24
	mins := int(d.Minutes()) % 60
	
	fmt.Printf(" %s up %d days, %2d:%02d\n", now, days, hours, mins)
}
