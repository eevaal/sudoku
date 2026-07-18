$ErrorActionPreference = "Continue"
$pipeName = "SudokuBridgePipe"

# Prevent multiple instances
try {
    $testPipe = New-Object System.IO.Pipes.NamedPipeClientStream(".", $pipeName, 'InOut', 'None')
    $testPipe.Connect(100)
    $testPipe.Dispose()
    exit # Already running
} catch {
    # Not running, safe to start
}

while ($true) {
    $pipe = $null
    try {
        $pipe = New-Object System.IO.Pipes.NamedPipeServerStream($pipeName, 'InOut', 1, 'Byte', 'None', 65536, 65536)
        $pipe.WaitForConnection()

        $reader = New-Object System.IO.StreamReader($pipe, [System.Text.Encoding]::UTF8)
        $writer = New-Object System.IO.StreamWriter($pipe, [System.Text.Encoding]::UTF8)
        $writer.AutoFlush = $true

        $jsonLine = $reader.ReadLine()
        if (-not [string]::IsNullOrWhiteSpace($jsonLine)) {
            $payload = $jsonLine | ConvertFrom-Json
            $cmdName = $payload.command
            
            $argsToPass = @()
            if ($payload.args -ne $null) {
                foreach ($arg in $payload.args) {
                    $argsToPass += [string]$arg
                }
            }

            $cmdObj = Get-Command $cmdName -CommandType Cmdlet,Function,Alias -ErrorAction SilentlyContinue
            if ($cmdObj) {
                $moduleName = $cmdObj.ModuleName
                if (-not $moduleName -and $cmdObj.CommandType -eq 'Alias') {
                    $moduleName = $cmdObj.ResolvedCommand.ModuleName
                }
                if ($moduleName) {
                    Import-Module $moduleName -ErrorAction SilentlyContinue
                }

                $result = & $cmdObj @argsToPass | Out-String
                $writer.Write($result)
            } else {
                $writer.Write("Command not found: $cmdName`n")
            }
        }
    } catch {
        # ignore and continue
    } finally {
        if ($pipe -ne $null) {
            if ($pipe.IsConnected) {
                $pipe.Disconnect()
            }
            $pipe.Dispose()
        }
    }
}
