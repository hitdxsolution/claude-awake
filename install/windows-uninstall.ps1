# claude-awake 제거 — 작업 스케줄러 등록을 지우고 설치된 실행파일을 삭제한다.
$ErrorActionPreference = "SilentlyContinue"
# 콘솔 한글이 깨지지 않도록 출력 인코딩을 UTF-8 로 맞춘다.
$OutputEncoding = [Console]::OutputEncoding = [System.Text.Encoding]::UTF8

Stop-ScheduledTask -TaskName "claude-awake"
Unregister-ScheduledTask -TaskName "claude-awake" -Confirm:$false
Remove-Item (Join-Path $env:LOCALAPPDATA "claude-awake") -Recurse -Force

Write-Host "제거했습니다. 절전 설정은 윈도우의 원래 값으로 돌아갑니다."
