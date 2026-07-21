$ErrorActionPreference = "Stop"

# Temporarily remove bridge from session PATH to prevent calling our own wrappers
$env:PATH = ($env:PATH -split ';' | Where-Object { $_ -notlike '*\.sudoku\bridge*' }) -join ';'

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

$serverSource = Join-Path $PWD "server.ps1"
$serverDest = Join-Path $env:USERPROFILE ".sudoku\server.ps1"
if (Test-Path $serverSource) {
    Copy-Item -Path $serverSource -Destination $serverDest -Force
}

$runKey = "HKCU:\Software\Microsoft\Windows\CurrentVersion\Run"
$runValue = "powershell -WindowStyle Hidden -NoProfile -ExecutionPolicy Bypass -File `"$serverDest`""
Set-ItemProperty -Path $runKey -Name "SudokuServer" -Value $runValue -ErrorAction SilentlyContinue

# Start the server now
Start-Process powershell -WindowStyle Hidden -ArgumentList "-NoProfile", "-ExecutionPolicy", "Bypass", "-File", "`"$serverDest`""

$bridgeGeneratorPath = Join-Path $env:USERPROFILE ".sudoku\generate_bridge.ps1"
$generatorContent = @"
`$ErrorActionPreference = "Stop"
`$BridgeDir = Join-Path `$env:USERPROFILE ".sudoku\bridge"
if (-not (Test-Path `$BridgeDir)) { New-Item -ItemType Directory -Force -Path `$BridgeDir | Out-Null }

`$cmdletObjs = Get-Command -CommandType Cmdlet, Function, Alias | Where-Object { `$_.Name -match "^[a-zA-Z\-]+`$" } | Group-Object Name | ForEach-Object { `$_.Group[0] }
foreach (`$cmdObj in `$cmdletObjs) {
    `$cmd = `$cmdObj.Name
    if (`$cmd -in "ls", "rm", "cp", "mv", "cat", "pwd", "mkdir", "clear", "sleep", "echo", "head", "tail", "wc") { continue }
    `$batPath = Join-Path `$BridgeDir "`$cmd.bat"
    
    `$batContent = "@`"`$env:USERPROFILE\.sudoku\bin\bridge.exe`" `"`$cmd`" %*"
    Set-Content -Path `$batPath -Value `$batContent
}

# 2.2 Download BusyBox and configure UNIX applets
`$busyboxUrl = "https://frippery.org/files/busybox/busybox.exe"
`$busyboxPath = Join-Path `$env:USERPROFILE ".sudoku\bin\busybox.exe"
`$appletsPath = Join-Path `$env:USERPROFILE ".sudoku\bin\applets.txt"

if (-not (Test-Path `$busyboxPath)) {
    Write-Host "[+] Downloading BusyBox..." -ForegroundColor Yellow
    Invoke-WebRequest -Uri `$busyboxUrl -OutFile `$busyboxPath
}

if (Test-Path `$busyboxPath) {
    Write-Host "[+] Configuring UNIX commands via BusyBox..." -ForegroundColor Yellow
    # Get list of all supported applets, using cmd /c to ensure pure ASCII encoding (no UTF-16 issues)
    cmd /c "`"`$busyboxPath`" --list > `"`$appletsPath`""

    # Create wrappers for all applets
    `$applets = Get-Content `$appletsPath
    foreach (`$applet in `$applets) {
        if (-not [string]::IsNullOrWhiteSpace(`$applet)) {
            `$batPath = Join-Path `$BridgeDir "`$applet.bat"
            `$batContent = "@`"`$env:USERPROFILE\.sudoku\bin\bridge.exe`" `"`$applet`" %*"
            Set-Content -Path `$batPath -Value `$batContent
        }
    }
}
"@

Set-Content -Path $bridgeGeneratorPath -Value $generatorContent
Start-Process powershell -WindowStyle Hidden -ArgumentList "-ExecutionPolicy", "Bypass", "-File", "`"$bridgeGeneratorPath`""
Write-Host "[+] Bridge wrappers generation started asynchronously in the background." -ForegroundColor Green

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
`$env:PATH = (`$env:PATH -split ';' | Where-Object { `$_.ToLower() -notlike "*\.sudoku\bridge*" }) -join ';'
`$aliasesToRemove = @('dir', 'echo', 'copy', 'del', 'move', 'type', 'cat', 'ls', 'rm', 'cp', 'mv', 'pwd', 'sleep', 'clear', 'mkdir', 'kill')
`$existing = Get-Alias | Select-Object -ExpandProperty Name
foreach (`$a in `$aliasesToRemove) {
    if (`$existing -contains `$a) { Remove-Item "Alias:`$a" -Force -ErrorAction Ignore }
}
# --- END SUDOKU ALIAS FIX ---
"@

$profileContent = ""
if (Test-Path $profilePath) {
    $profileContent = Get-Content $profilePath -Raw
}
if ($null -eq $profileContent) { $profileContent = "" }

if ($profileContent -match "(?s)# --- BEGIN SUDOKU ALIAS FIX ---.*# --- END SUDOKU ALIAS FIX ---") {
    $profileContent = $profileContent -replace "(?s)# --- BEGIN SUDOKU ALIAS FIX ---.*# --- END SUDOKU ALIAS FIX ---", ($profileSnippet -replace '\$', '$$$$')
    Set-Content -Path $profilePath -Value $profileContent
    Write-Host "[+] Updated alias remover in `$PROFILE." -ForegroundColor Green
} else {
    Add-Content -Path $profilePath -Value "`n$profileSnippet`n"
    Write-Host "[+] Added alias remover to `$PROFILE." -ForegroundColor Green
}

# Update PATH (Persistent for MACHINE to override System32)
if ($isAdmin) {
    $MachinePath = [Environment]::GetEnvironmentVariable("PATH", "Machine")
    $MachinePathArray = $MachinePath -split ';' | Where-Object { $_ -ne "" -and $_ -ne $BinDir -and $_ -ne $BridgeDir }
    $NewMachinePath = "$BinDir;" + ($MachinePathArray -join ';')
    [Environment]::SetEnvironmentVariable("PATH", $NewMachinePath, "Machine")
    Write-Host "[+] Binaries prepended to global SYSTEM PATH." -ForegroundColor Green
}

# Clean User PATH
$UserPath = [Environment]::GetEnvironmentVariable("PATH", "User")
if ($null -ne $UserPath) {
    $UserPathArray = $UserPath -split ';' | Where-Object { $_ -ne "" -and $_ -ne $BinDir -and $_ -ne $BridgeDir }
    $NewUserPath = "$BinDir;" + ($UserPathArray -join ';')
    [Environment]::SetEnvironmentVariable("PATH", $NewUserPath, "User")
    Write-Host "[+] Binaries added to User PATH." -ForegroundColor Green
}

# Setup CMD AutoRun to inject BridgeDir for CMD only
$cmdAutoRunKey = "HKCU:\Software\Microsoft\Command Processor"
if (-not (Test-Path $cmdAutoRunKey)) {
    New-Item -Path $cmdAutoRunKey -Force | Out-Null
}
$autoRunScript = Join-Path $env:USERPROFILE ".sudoku\cmd_autorun.cmd"
Set-Content -Path $autoRunScript -Value "@set PATH=%PATH%;$BridgeDir"
Set-ItemProperty -Path $cmdAutoRunKey -Name "AutoRun" -Value "`"$autoRunScript`""
Write-Host "[+] CMD AutoRun configured to inject Sudoku Bridge commands." -ForegroundColor Green

# Update PATH (Current session)
$SessionPathArray = $env:PATH -split ';' | Where-Object { $_ -ne "" -and $_ -ne $BinDir -and $_ -ne $BridgeDir }
$env:PATH = "$BinDir;" + ($SessionPathArray -join ';')
Write-Host "[+] Current session PATH updated." -ForegroundColor Green

Write-Host "=== Installation completed successfully! ===" -ForegroundColor Cyan
Write-Host "You can now use your new commands across both PowerShell and CMD."
