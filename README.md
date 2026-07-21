<p align="center">
  <img src="./logo.png" width="256" height="256" alt="Sudoku Logo">
</p>

# Sudoku

Sudoku is a lightweight system package for Windows designed to erase the borders between PowerShell, CMD, and Unix terminals. It provides native, fast, and seamless access to over 150 Unix commands directly in your Windows environment — without WSL, virtual machines, or heavy dependencies.

## Features

- **150+ native Unix commands**: Powered by [BusyBox](https://frippery.org/busybox/), all classic utilities (`ls`, `grep`, `awk`, `sed`, `find`, `vi`, `tar`, and many more) work out of the box with full flag support.
- **Sudoku PowerBridge**: True two-way interoperability between PowerShell and CMD. Automatically generates over 2,000 asynchronous wrappers so that you can run any PowerShell Cmdlet directly from CMD.
- **Smart routing**: The universal `bridge.exe` intelligently routes Unix commands directly to BusyBox (millisecond response, full TTY support for interactive programs like `vi`) and PowerShell cmdlets through a background named-pipe server.
- **Lightweight**: Only two binaries — `bridge.exe` (~2 MB) and `busybox.exe` (~1 MB). No bloat.

## How It Works

```
User types: ls -la
    ↓
bridge.exe checks applets.txt → "ls" is a BusyBox applet
    ↓
Runs: busybox.exe ls -la  (direct, with full TTY)
    ↓
Output displayed instantly
```

```
User types: Get-Process
    ↓
bridge.exe checks applets.txt → NOT a BusyBox applet
    ↓
Sends JSON via named pipe → server.ps1 executes → returns result
```

## Available Commands

All standard Unix utilities are supported via BusyBox, including:

**Files & Directories:** `ls`, `cp`, `mv`, `rm`, `mkdir`, `rmdir`, `touch`, `ln`, `find`, `stat`, `chmod`

**Text Processing:** `cat`, `grep`, `egrep`, `fgrep`, `awk`, `sed`, `sort`, `uniq`, `wc`, `cut`, `tr`, `head`, `tail`, `less`, `diff`

**Archives:** `tar`, `gzip`, `gunzip`, `bzip2`, `xz`, `unzip`

**Networking:** `wget`, `nc` (netcat), `whois`, `httpd`

**System:** `ps`, `kill`, `killall`, `df`, `du`, `free`, `uptime`, `uname`, `whoami`, `env`, `id`

**Editors & Tools:** `vi`, `ed`, `bc`, `cal`, `seq`, `yes`, `watch`, `xargs`, `tee`

**Hashing & Encoding:** `md5sum`, `sha256sum`, `sha512sum`, `base64`, `xxd`, `hexdump`

**Shell:** `sh`, `bash`, `ash` — full POSIX shell for running scripts

Plus all PowerShell cmdlets are accessible from CMD via the bridge.

## Installation

There are two ways to install or update Sudoku.

### Automatic Installation (Recommended)

Open PowerShell (as Administrator for system-wide PATH) and run:

```powershell
iex (irm "https://raw.githubusercontent.com/eevaal/sudoku/main/install.ps1")
```

This will automatically clone the repository, compile the bridge, download BusyBox, and configure everything.

### Manual Installation

1. Clone the repository:
   ```
   git clone https://github.com/eevaal/sudoku.git
   cd sudoku
   ```
2. Run the installation script:
   ```powershell
   .\install.ps1
   ```

Once installed, the binaries are placed in `~/.sudoku/bin` and the Sudoku PowerBridge is configured globally.

## Usage

Use the commands exactly as you would in a Unix terminal.

```bash
# Remove directories recursively
rm -rf ./build/*

# Search for text with regex
grep -vE "^#" config.txt

# Process text with awk
ls -la | awk '{print $5, $9}'

# Archive a directory
tar czf backup.tar.gz ./project

# Open the built-in vi editor
vi README.md

# Enter a POSIX shell
sh
```

## Sudoku PowerBridge

Sudoku PowerBridge erases the boundaries between Windows CMD and PowerShell.

When you run the installation script, Sudoku generates a lightweight `bridge` directory containing asynchronous `.bat` wrappers for **over 2,000 PowerShell cmdlets**.

- **From CMD**: You can natively call any PowerShell cmdlet (e.g. `Get-Process`) directly, without typing `powershell -c ...`. Sudoku automatically resolves the module, handles auto-loading, and passes your arguments seamlessly.

**Note on Complex PowerShell Expressions in CMD:**
Because of the way Windows CMD parses parentheses, complex PowerShell expressions starting with an opening parenthesis, such as `(New-Object -ComObject SAPI.SpVoice).Speak(...)`, cannot be evaluated directly as the first token.
To execute such complex commands natively from CMD, wrap them in the `Write-Output` cmdlet:

```cmd
C:\> Write-Output ((New-Object -ComObject SAPI.SpVoice).Speak('Sudoku Power Bridge is alive'))
```

## License

This project is open-source. Please see the LICENSE file for more information.
