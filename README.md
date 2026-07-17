<p align="center">
  <img src="logo.png" alt="Sudoku Logo" width="200">
</p>

# Sudoku

Sudoku is a lightweight system package for Windows designed to erase the borders between PowerShell, CMD, and Unix terminals. It provides native, fast, and seamless implementations of popular Unix commands directly in your Windows environment without relying on WSL or heavy virtualization.

## Features

- **Native Windows API integration**: Commands like `sudo` automatically trigger the standard UAC prompt for privilege elevation.
- **Cross-shell compatibility**: Works seamlessly in both PowerShell and CMD.
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

## License

This project is open-source. Please see the LICENSE file for more information.
