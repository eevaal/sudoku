package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"syscall"
	"time"
)

type Request struct {
	Command string   `json:"command"`
	Args    []string `json:"args"`
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: bridge <command> [args...]")
		os.Exit(1)
	}

	cmdName := os.Args[1]
	args := os.Args[2:]

	req := Request{
		Command: cmdName,
		Args:    args,
	}

	payload, err := json.Marshal(req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error encoding payload: %v\n", err)
		os.Exit(1)
	}
	payload = append(payload, '\n')

	pipeName := `\\.\pipe\SudokuBridgePipe_v2`
	
	var conn io.ReadWriteCloser
	var openErr error
	for i := 0; i < 3; i++ {
		// On Windows, named pipes can be opened like files
		file, err := os.OpenFile(pipeName, os.O_RDWR, 0)
		if err == nil {
			conn = file
			break
		}
		openErr = err
		
		// If pipe not found on first try, attempt to launch server
		if i == 0 {
			userProfile := os.Getenv("USERPROFILE")
			serverScript := userProfile + `\.sudoku\server.ps1`
			if _, err := os.Stat(serverScript); err == nil {
				cmd := exec.Command("powershell", "-WindowStyle", "Hidden", "-NoProfile", "-ExecutionPolicy", "Bypass", "-File", serverScript)
				cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true, CreationFlags: 0x08000000}
				cmd.Start()
			}
		}
		time.Sleep(500 * time.Millisecond)
	}

	if conn == nil {
		fmt.Fprintf(os.Stderr, "Error: Could not connect to background PowerShell server: %v\n", openErr)
		os.Exit(1)
	}
	defer conn.Close()

	_, err = conn.Write(payload)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing to pipe: %v\n", err)
		os.Exit(1)
	}

	reader := bufio.NewReader(conn)
	io.Copy(os.Stdout, reader)
}
