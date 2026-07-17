$ErrorActionPreference = "Stop"

# Check for Administrator privileges
$isAdmin = ([Security.Principal.WindowsPrincipal][Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)

if (-not $isAdmin) {
    Write-Host "This script requires Administrator privileges to override the native Windows sudo." -ForegroundColor Red
    Write-Host "Please open PowerShell as Administrator and run .\install.ps1 again." -ForegroundColor Yellow
    exit
}

# Paths
$BinDir = Join-Path $env:USERPROFILE ".sudoku\bin"
$SudoCmd = "cmd/sudo"
$RmCmd = "cmd/rm"

Write-Host "=== Sudoku Package Installation ===" -ForegroundColor Cyan

# Create directory
if (-not (Test-Path $BinDir)) {
    New-Item -ItemType Directory -Force -Path $BinDir | Out-Null
    Write-Host "[+] Created directory: $BinDir" -ForegroundColor Green
} else {
    Write-Host "[v] Directory already exists: $BinDir" -ForegroundColor DarkGreen
}

# Compile
Write-Host "[+] Compiling sudo.exe..." -ForegroundColor Yellow
go build -o "$BinDir\sudo.exe" "./$SudoCmd"
if ($LASTEXITCODE -ne 0 -and $LASTEXITCODE -ne $null) {
    Write-Error "Error compiling sudo.exe"
}

Write-Host "[+] Compiling rm.exe..." -ForegroundColor Yellow
go build -o "$BinDir\rm.exe" "./$RmCmd"
if ($LASTEXITCODE -ne 0 -and $LASTEXITCODE -ne $null) {
    Write-Error "Error compiling rm.exe"
}

# Update PATH (Persistent for MACHINE to override System32)
$MachinePath = [Environment]::GetEnvironmentVariable("PATH", "Machine")
$MachinePathArray = $MachinePath -split ';' | Where-Object { $_ -ne "" -and $_ -ne $BinDir }
$NewMachinePath = "$BinDir;" + ($MachinePathArray -join ';')
[Environment]::SetEnvironmentVariable("PATH", $NewMachinePath, "Machine")
Write-Host "[+] Directory prepended to global SYSTEM PATH (overrides native Windows sudo)." -ForegroundColor Green

# Remove from User PATH if it was added previously to avoid duplicates
$UserPath = [Environment]::GetEnvironmentVariable("PATH", "User")
if ($UserPath -like "*$BinDir*") {
    $UserPathArray = $UserPath -split ';' | Where-Object { $_ -ne "" -and $_ -ne $BinDir }
    $NewUserPath = $UserPathArray -join ';'
    [Environment]::SetEnvironmentVariable("PATH", $NewUserPath, "User")
    Write-Host "[-] Removed duplicate from User PATH." -ForegroundColor Gray
}

# Update PATH (Current session)
$SessionPathArray = $env:PATH -split ';' | Where-Object { $_ -ne "" -and $_ -ne $BinDir }
$env:PATH = "$BinDir;" + ($SessionPathArray -join ';')
Write-Host "[+] Current session PATH updated." -ForegroundColor Green

Write-Host "=== Installation completed successfully! ===" -ForegroundColor Cyan
Write-Host "You can now use 'sudo' and 'rm -rf' in your terminal."
Write-Host "Press any key to close..."
$Host.UI.RawUI.ReadKey("NoEcho,IncludeKeyDown") | Out-Null
