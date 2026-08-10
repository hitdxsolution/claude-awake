# Install claude-awake as a Scheduled Task: runs at logon, restarts if it dies, no time limit.
#   Usage:  powershell -ExecutionPolicy Bypass -File .\windows-install.ps1 [-BinSrc path\to\exe]
#   With no argument it uses ..\dist\claude-awake-windows-x64.exe from the repo.
param([string]$BinSrc = "")
$ErrorActionPreference = "Stop"

if (-not $BinSrc) { $BinSrc = Join-Path $PSScriptRoot "..\dist\claude-awake-windows-x64.exe" }
if (-not (Test-Path $BinSrc)) {
  Write-Error "Binary not found: $BinSrc`nPass -BinSrc <path> or build it with 'make windows-x64'."
  exit 1
}

$dest = Join-Path $env:LOCALAPPDATA "claude-awake\claude-awake.exe"
New-Item -ItemType Directory -Force -Path (Split-Path $dest) | Out-Null
Copy-Item $BinSrc $dest -Force

# The binary is built for the GUI subsystem (-H windowsgui), so it never opens a console
# window and can be launched directly — no wscript shim needed.
$action  = New-ScheduledTaskAction -Execute $dest
$trigger = New-ScheduledTaskTrigger -AtLogOn
$settings = New-ScheduledTaskSettingsSet `
  -AllowStartIfOnBatteries -DontStopIfGoingOnBatteries `
  -RestartCount 999 -RestartInterval (New-TimeSpan -Minutes 1) `
  -ExecutionTimeLimit ([TimeSpan]::Zero)

Register-ScheduledTask -TaskName "claude-awake" -Action $action -Trigger $trigger -Settings $settings -Force | Out-Null
Start-ScheduledTask -TaskName "claude-awake"

Write-Host "Installed as scheduled task 'claude-awake'. It runs at logon and while Claude is open."
Write-Host "  binary: $dest"
Write-Host "  log:    $(Join-Path $env:LOCALAPPDATA 'claude-awake\claude-awake.log')"
