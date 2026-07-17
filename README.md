![Sudoku Logo](./logo.png)

# Sudoku

Sudoku is a lightweight system package for Windows designed to erase the borders between PowerShell, CMD, and Unix terminals. It provides native, fast, and seamless implementations of popular Unix commands directly in your Windows environment without relying on WSL or heavy virtualization.

## Features

- **Native Windows API integration**: Commands like `sudo` automatically trigger the standard UAC prompt for privilege elevation.
- **Sudoku PowerBridge**: True two-way interoperability between PowerShell and CMD. Automatically generates over 2,000 asynchronous wrappers so that you can run any PowerShell Cmdlet directly from CMD, and all CMD built-ins from PowerShell.
- **Zero dependencies**: Written in Go and compiled to standalone lightweight binaries.

## Available Commands

The following Unix-like commands are currently implemented:

- `sudo` - Execute commands with Administrator privileges.
- `rm` - Safely remove files and directories (supports globbing and -rf).
- `cp` - Copy files and directories.
- `mv` - Move or rename files and directories.
- `ls` - List directory contents.
- `pwd` - Print working directory.
- `touch` - Create empty files or update timestamps.
- `cat` - Concatenate and print files.
- `grep` - Search text using regular expressions.
- `find` - Search for files in a directory hierarchy.
- `head` / `tail` - Output the first/last parts of files.
- `wc` - Print newline, word, and byte counts.
- `whoami` - Print effective user ID.
- `sleep` - Delay for a specified amount of time.
- `clear` - Clear the terminal screen.
- `winch` - Built-in WINdows fastfetCH clone with ASCII logo.
- `uname` - Print system information.
- `uptime` - Print how long the system has been running.
- `nvim` - Wraps and automatically installs the official Neovim editor.

## Installation

To install or update the package, run the included PowerShell script. The script automatically compiles all commands and prepends them to your system's PATH variable to override native Windows commands where necessary.

1. Open PowerShell as Administrator.
2. Navigate to the project directory.
3. Run the installation script:

```powershell
.\install.ps1
```

Once installed, the binaries are placed in `~/.sudoku/bin` and are immediately available globally.

## Usage

Use the commands exactly as you would in a Unix terminal.

```powershell
# Elevate privileges for a single command
sudo notepad.exe

# Open an elevated shell
sudo -s

# Remove directories recursively
rm -rf ./build/*

# Search for text in a file
grep -i "error" app.log
```

## Sudoku PowerBridge

Sudoku PowerBridge is a core feature that erases the boundaries between Windows CMD and PowerShell. 

When you run the installation script, Sudoku generates a lightweight `bridge` directory containing asynchronous `.bat` wrappers for **over 2,000 PowerShell cmdlets** and all **CMD built-in commands** (like `mklink`, `assoc`, `title`, etc.).

- **From CMD**: You can natively call any PowerShell cmdlet (e.g. `Get-Process`) directly, without typing `powershell -c ...`. Sudoku automatically resolves the module, handles auto-loading, and passes your arguments seamlessly.
- **From PowerShell**: You can call CMD built-in commands directly as if they were native executables.

The generation process is highly optimized and runs completely asynchronously in the background, so it doesn't block your terminal or cause any black screens during installation.

## License

This project is open-source. Please see the LICENSE file for more information.
