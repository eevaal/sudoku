$ErrorActionPreference = "Stop"

# Check for Administrator privileges
$isAdmin = ([Security.Principal.WindowsPrincipal][Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)

if (-not $isAdmin) {
    Write-Host "[!] Note: Running without Administrator privileges. Machine-wide PATH will not be updated." -ForegroundColor Yellow
    Write-Host "    Only User PATH will be updated." -ForegroundColor Yellow
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

# Dynamic Compilation
$CmdDir = Join-Path $PWD "cmd"
if (Test-Path $CmdDir) {
    $tools = Get-ChildItem -Path $CmdDir -Directory
    foreach ($tool in $tools) {
        $toolName = $tool.Name
        Write-Host "[+] Compiling $toolName.exe..." -ForegroundColor Yellow
        go build -o "$BinDir\$toolName.exe" "./cmd/$toolName"
        if ($LASTEXITCODE -ne 0 -and $LASTEXITCODE -ne $null) {
            Write-Error "Error compiling $toolName.exe"
        }
    }
} else {
    Write-Warning "Directory 'cmd' not found!"
}

# Setup Bridge Directory
$BridgeDir = Join-Path $env:USERPROFILE ".sudoku\bridge"
if (-not (Test-Path $BridgeDir)) {
    New-Item -ItemType Directory -Force -Path $BridgeDir | Out-Null
    Write-Host "[+] Created bridge directory: $BridgeDir" -ForegroundColor Green
}

# 1. PowerShell -> CMD (Generate wrappers for all PowerShell Cmdlets)
Write-Host "[+] Generating PowerShell wrappers for CMD (this may take a few seconds on first run)..." -ForegroundColor Yellow
$cmdlets = Get-Command -CommandType Cmdlet, Function, Alias | Where-Object { $_.Name -match "^[a-zA-Z\-]+$" } | Select-Object -ExpandProperty Name -Unique
$pwshWrappersCount = 0
foreach ($cmd in $cmdlets) {
    # Skip short unix-like aliases so they don't override our Go binaries in CMD
    if ($cmd -in "ls", "rm", "cp", "mv", "cat", "pwd", "mkdir", "clear", "sleep", "echo", "head", "tail", "wc") { continue }
    $batPath = Join-Path $BridgeDir "$cmd.bat"
    
    $cmdObj = Get-Command $cmd -CommandType Cmdlet, Function, Alias -ErrorAction SilentlyContinue
    $moduleName = $cmdObj.ModuleName
    if (-not $moduleName -and $cmdObj.CommandType -eq 'Alias') {
        $moduleName = $cmdObj.ResolvedCommand.ModuleName
    }
    $importStr = ""
    if ($moduleName) {
        $importStr = "Import-Module $moduleName -ErrorAction SilentlyContinue; "
    }

    # Only write if it doesn't exist to save time on subsequent runs
    if (-not (Test-Path $batPath) -or (Get-Content $batPath) -notmatch "Import-Module") {
        Set-Content -Path $batPath -Value "@powershell -NoProfile -Command `"$importStr& (Get-Command $cmd -CommandType Cmdlet,Function,Alias) %*`""
        $pwshWrappersCount++
    }
}
Write-Host "[+] Generated $pwshWrappersCount new wrappers for PowerShell cmdlets." -ForegroundColor Green

# 2. CMD -> PowerShell (Generate wrappers for CMD builtins)
Write-Host "[+] Generating CMD built-in wrappers for PowerShell..." -ForegroundColor Yellow
$cmdBuiltins = @("assoc", "ftype", "mklink", "vol", "ver", "title", "color", "start", "md", "rd", "ren", "rename", "call", "pushd", "popd", "doskey")
$cmdWrappersCount = 0
foreach ($cmd in $cmdBuiltins) {
    $batPath = Join-Path $BridgeDir "$cmd.bat"
    if (-not (Test-Path $batPath)) {
        Set-Content -Path $batPath -Value "@cmd /c $cmd %*"
        $cmdWrappersCount++
    }
}
Write-Host "[+] Generated $cmdWrappersCount new wrappers for CMD built-ins." -ForegroundColor Green

# 3. Fix PowerShell Aliases in $PROFILE
Write-Host "[+] Patching PowerShell `$PROFILE to remove conflicting aliases..." -ForegroundColor Yellow
$profilePath = $PROFILE
if ($null -eq $profilePath -or $profilePath -eq "") {
    $profilePath = Join-Path (Join-Path $env:USERPROFILE "Documents") "WindowsPowerShell\Microsoft.PowerShell_profile.ps1"
}
if (-not (Test-Path $profilePath)) {
    $profileDir = Split-Path $profilePath
    if (-not (Test-Path $profileDir)) { New-Item -ItemType Directory -Force -Path $profileDir | Out-Null }
    New-Item -ItemType File -Force -Path $profilePath | Out-Null
}

$profileSnippet = @"
# --- BEGIN SUDOKU ALIAS FIX ---
`$aliasesToRemove = @('dir', 'echo', 'copy', 'del', 'move', 'type', 'cat', 'ls', 'rm', 'cp', 'mv', 'pwd', 'sleep', 'clear', 'mkdir', 'kill')
foreach (`$a in `$aliasesToRemove) {
    if (Test-Path "Alias:`$a") { Remove-Item "Alias:`$a" -Force -ErrorAction SilentlyContinue }
}
# --- END SUDOKU ALIAS FIX ---
"@

$profileContent = ""
if (Test-Path $profilePath) {
    $profileContent = Get-Content $profilePath -Raw
}
if ($null -eq $profileContent -or -not ($profileContent -match "BEGIN SUDOKU ALIAS FIX")) {
    Add-Content -Path $profilePath -Value "`n$profileSnippet`n"
    Write-Host "[+] Added alias remover to `$PROFILE." -ForegroundColor Green
} else {
    Write-Host "[v] Profile already patched." -ForegroundColor DarkGreen
}

# Update PATH (Persistent for MACHINE to override System32)
if ($isAdmin) {
    $MachinePath = [Environment]::GetEnvironmentVariable("PATH", "Machine")
    $MachinePathArray = $MachinePath -split ';' | Where-Object { $_ -ne "" -and $_ -ne $BinDir -and $_ -ne $BridgeDir }
    $NewMachinePath = "$BinDir;$BridgeDir;" + ($MachinePathArray -join ';')
    [Environment]::SetEnvironmentVariable("PATH", $NewMachinePath, "Machine")
    Write-Host "[+] Binaries and Bridge prepended to global SYSTEM PATH." -ForegroundColor Green
}

# Remove from User PATH if it was added previously to avoid duplicates
$UserPath = [Environment]::GetEnvironmentVariable("PATH", "User")
if ($UserPath -like "*$BinDir*" -or $UserPath -like "*$BridgeDir*") {
    $UserPathArray = $UserPath -split ';' | Where-Object { $_ -ne "" -and $_ -ne $BinDir -and $_ -ne $BridgeDir }
    $NewUserPath = $UserPathArray -join ';'
    [Environment]::SetEnvironmentVariable("PATH", $NewUserPath, "User")
    Write-Host "[-] Removed duplicate from User PATH." -ForegroundColor Gray
}

# Update PATH (Current session)
$SessionPathArray = $env:PATH -split ';' | Where-Object { $_ -ne "" -and $_ -ne $BinDir -and $_ -ne $BridgeDir }
$env:PATH = "$BinDir;$BridgeDir;" + ($SessionPathArray -join ';')
Write-Host "[+] Current session PATH updated." -ForegroundColor Green

Write-Host "=== Installation completed successfully! ===" -ForegroundColor Cyan
Write-Host "You can now use your new commands across both PowerShell and CMD."
