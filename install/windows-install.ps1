# Install claude-awake as a Scheduled Task: runs at logon, restarts if it dies, no time limit.
#   Usage:  powershell -ExecutionPolicy Bypass -File .\windows-install.ps1 [-BinSrc path\to\exe]
#   With no argument it uses ..\dist\claude-awake-windows-x64.exe from a local build.
param([string]$BinSrc = "")
$ErrorActionPreference = "Stop"

if (-not $BinSrc) { $BinSrc = Join-Path $PSScriptRoot "..\dist\claude-awake-windows-x64.exe" }
if (-not (Test-Path $BinSrc)) {
  Write-Error "Binary not found: $BinSrc`nBuild it first (bun run build) or pass -BinSrc <path>."
  exit 1
}

$dest = Join-Path $env:LOCALAPPDATA "claude-awake\claude-awake.exe"
New-Item -ItemType Directory -Force -Path (Split-Path $dest) | Out-Null
Copy-Item $BinSrc $dest -Force

# The binary is a console app (cross-compiled from macOS), so launch it through a hidden
#   VBScript shim (window style 0) to avoid a flashing console window.
$vbs = Join-Path (Split-Path $dest) "run-hidden.vbs"
Set-Content -Path $vbs -Encoding ASCII -Value ('CreateObject("WScript.Shell").Run """' + $dest + '""", 0, False')

$action  = New-ScheduledTaskAction -Execute "wscript.exe" -Argument ('"' + $vbs + '"')
$trigger = New-ScheduledTaskTrigger -AtLogOn
$settings = New-ScheduledTaskSettingsSet `
  -AllowStartIfOnBatteries -DontStopIfGoingOnBatteries `
  -RestartCount 999 -RestartInterval (New-TimeSpan -Minutes 1) `
  -ExecutionTimeLimit ([TimeSpan]::Zero)

Register-ScheduledTask -TaskName "claude-awake" -Action $action -Trigger $trigger -Settings $settings -Force | Out-Null
Start-ScheduledTask -TaskName "claude-awake"

Write-Host "Installed as scheduled task 'claude-awake'. It runs at logon and while Claude is open."
Write-Host "  binary: $dest"
