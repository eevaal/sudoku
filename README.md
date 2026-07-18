<p align="center">
  <img src="./logo.png" width="256" height="256" alt="Sudoku Logo">
</p>

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

There are two ways to install or update Sudoku: the recommended automatic method or the manual method. Both methods automatically compile all commands and prepend them to your system's PATH variable to override native Windows commands where necessary.

### Automatic Installation (Recommended)

You can install Sudoku with a single command. Open PowerShell as Administrator and run the following:

```powershell
iex (irm "https://raw.githubusercontent.com/eevaal/sudoku/main/install.ps1")
```

### Manual Installation

If you prefer to download the files manually:

1. Download the repository archive and extract it.
2. Open PowerShell as Administrator.
3. Navigate to the extracted project directory.
4. Run the installation script:

```powershell
.\install.ps1
```

Once installed, the binaries are placed in `~/.sudoku/bin` and the Sudoku PowerBridge is configured globally.


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

**Note on Complex PowerShell Expressions in CMD:**
Because of the way Windows CMD parses parentheses, complex PowerShell expressions starting with an opening parenthesis, such as `(New-Object -ComObject SAPI.SpVoice).Speak(...)`, cannot be evaluated directly as the first token. 
To execute such complex commands natively from CMD, wrap them in the `Write-Output` cmdlet. This correctly invokes the wrapper and allows PowerShell to evaluate the entire pipeline:

```cmd
C:\> Write-Output ((New-Object -ComObject SAPI.SpVoice).Speak("Sudoku Power Bridge is alive"))
```

The generation process is highly optimized and runs completely asynchronously in the background, so it doesn't block your terminal or cause any black screens during installation.

## License

This project is open-source. Please see the LICENSE file for more information.
