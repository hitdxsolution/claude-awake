#!/bin/bash
# claude-awake 제거 — LaunchAgent 등록을 내리고 실행파일을 지운다.
set -euo pipefail
launchctl bootout "gui/$(id -u)/com.claude-awake" 2>/dev/null || true
rm -f "$HOME/Library/LaunchAgents/com.claude-awake.plist"
rm -f "$HOME/.local/bin/claude-awake"
echo "제거했습니다. 절전 설정은 맥의 원래 값으로 돌아갑니다."
