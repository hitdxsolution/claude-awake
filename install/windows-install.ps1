# claude-awake 설치 — 로그인할 때 자동 실행되고, 죽으면 다시 뜨도록 작업 스케줄러에 등록한다.
#   사용법:  powershell -ExecutionPolicy Bypass -File .\windows-install.ps1 [-BinSrc 실행파일경로]
#   인자를 안 주면 저장소의 ..\dist\claude-awake-windows-x64.exe 를 쓴다.
param([string]$BinSrc = "")
$ErrorActionPreference = "Stop"
# 콘솔 한글이 깨지지 않도록 출력 인코딩을 UTF-8 로 맞춘다.
$OutputEncoding = [Console]::OutputEncoding = [System.Text.Encoding]::UTF8

if (-not $BinSrc) { $BinSrc = Join-Path $PSScriptRoot "..\dist\claude-awake-windows-x64.exe" }
if (-not (Test-Path $BinSrc)) {
  Write-Error "실행파일을 찾을 수 없습니다: $BinSrc`n  make windows-x64 로 빌드하거나, -BinSrc 로 경로를 넘겨 주세요."
  exit 1
}

$dest = Join-Path $env:LOCALAPPDATA "claude-awake\claude-awake.exe"
New-Item -ItemType Directory -Force -Path (Split-Path $dest) | Out-Null
Copy-Item $BinSrc $dest -Force

# 이 실행파일은 GUI 서브시스템으로 빌드돼(-H windowsgui) 검은 콘솔 창이 뜨지 않는다.
#   그래서 별도 셸 래퍼 없이 그대로 등록하면 된다.
$action  = New-ScheduledTaskAction -Execute $dest
$trigger = New-ScheduledTaskTrigger -AtLogOn
$settings = New-ScheduledTaskSettingsSet `
  -AllowStartIfOnBatteries -DontStopIfGoingOnBatteries `
  -RestartCount 999 -RestartInterval (New-TimeSpan -Minutes 1) `
  -ExecutionTimeLimit ([TimeSpan]::Zero)

Register-ScheduledTask -TaskName "claude-awake" -Action $action -Trigger $trigger -Settings $settings -Force | Out-Null
Start-ScheduledTask -TaskName "claude-awake"

Write-Host "설치했습니다."
Write-Host "  이제 로그인할 때 자동으로 시작되고, Claude 앱이 켜져 있는 동안 절전을 막습니다."
Write-Host "  실행파일: $dest"
Write-Host "  로그:     $(Join-Path $env:LOCALAPPDATA 'claude-awake\claude-awake.log')"
Write-Host "  제거하려면: powershell -ExecutionPolicy Bypass -File .\install\windows-uninstall.ps1"
