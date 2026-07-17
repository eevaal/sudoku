package main

import (
	"fmt"
	"os"
	"strings"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

// getParentProcessName пытается определить имя родительского процесса (оболочки).
func getParentProcessName() string {
	pid := uint32(os.Getppid())
	snapshot, err := windows.CreateToolhelp32Snapshot(windows.TH32CS_SNAPPROCESS, 0)
	if err != nil {
		return "cmd.exe" // Значение по умолчанию в случае ошибки
	}
	defer windows.CloseHandle(snapshot)

	var procEntry windows.ProcessEntry32
	procEntry.Size = uint32(unsafe.Sizeof(procEntry))

	if err = windows.Process32First(snapshot, &procEntry); err != nil {
		return "cmd.exe"
	}

	for {
		if procEntry.ProcessID == pid {
			return windows.UTF16ToString(procEntry.ExeFile[:])
		}
		if err = windows.Process32Next(snapshot, &procEntry); err != nil {
			break
		}
	}
	return "cmd.exe"
}

// shellExecute вызывает функцию ShellExecute из Windows API с нужным verb (runas для прав админа).
func shellExecute(verb, exe, args, cwd string, showCmd int32) error {
	verbPtr, _ := syscall.UTF16PtrFromString(verb)
	exePtr, _ := syscall.UTF16PtrFromString(exe)
	cwdPtr, _ := syscall.UTF16PtrFromString(cwd)
	var argsPtr *uint16
	if args != "" {
		argsPtr, _ = syscall.UTF16PtrFromString(args)
	}

	ret, _, _ := syscall.NewLazyDLL("shell32.dll").NewProc("ShellExecuteW").Call(
		0,
		uintptr(unsafe.Pointer(verbPtr)),
		uintptr(unsafe.Pointer(exePtr)),
		uintptr(unsafe.Pointer(argsPtr)),
		uintptr(unsafe.Pointer(cwdPtr)),
		uintptr(showCmd),
	)

	// Если возвращаемое значение <= 32, это ошибка
	if ret <= 32 {
		return fmt.Errorf("ShellExecute failed with code %d", ret)
	}
	return nil
}

func main() {
	args := os.Args[1:]

	openShell := false
	if len(args) == 0 {
		openShell = true
	} else if len(args) == 1 && (args[0] == "-s" || args[0] == "--shell") {
		openShell = true
	}

	parentProc := strings.ToLower(getParentProcessName())
	isPowerShell := strings.Contains(parentProc, "powershell") || strings.Contains(parentProc, "pwsh")

	var exe string
	var cmdArgs string

	if openShell {
		// Просто открываем новую оболочку от имени администратора
		if isPowerShell {
			exe = "powershell.exe"
		} else {
			exe = "cmd.exe"
		}
	} else {
		// Запускаем переданную команду в новой оболочке с флагом, чтобы она не закрывалась
		command := strings.Join(args, " ")
		if isPowerShell {
			exe = "powershell.exe"
			// -NoExit предотвращает закрытие PowerShell после выполнения команды
			cmdArgs = fmt.Sprintf("-NoExit -Command \"%s\"", command)
		} else {
			exe = "cmd.exe"
			// /k выполняет команду и оставляет CMD открытым
			cmdArgs = fmt.Sprintf("/k \"%s\"", command)
		}
	}

	cwd, _ := os.Getwd()

	// SW_SHOWNORMAL = 1
	err := shellExecute("runas", exe, cmdArgs, cwd, 1)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка при запросе повышения прав: %v\n", err)
		os.Exit(1)
	}
}
