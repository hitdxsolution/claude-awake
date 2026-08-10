# Remove the claude-awake scheduled task and binary.
$ErrorActionPreference = "SilentlyContinue"
Stop-ScheduledTask -TaskName "claude-awake"
Unregister-ScheduledTask -TaskName "claude-awake" -Confirm:$false
Remove-Item (Join-Path $env:LOCALAPPDATA "claude-awake") -Recurse -Force
Write-Host "Uninstalled."
