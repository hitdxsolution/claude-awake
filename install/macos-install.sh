#!/bin/bash
# claude-awake 설치 — 로그인할 때 자동 실행되고 죽으면 다시 뜨도록 LaunchAgent 로 등록한다.
#   사용법: ./macos-install.sh [실행파일 경로]
#   인자를 안 주면 이 맥의 칩에 맞는 ../dist/claude-awake-macos-<칩> 을 자동으로 고른다.
set -euo pipefail

here="$(cd "$(dirname "$0")" && pwd)"
bin_src="${1:-}"
if [ -z "$bin_src" ]; then
  if [ "$(uname -m)" = "arm64" ]; then
    bin_src="$here/../dist/claude-awake-macos-arm64"   # 애플 실리콘(M1 이상)
  else
    bin_src="$here/../dist/claude-awake-macos-x64"     # 인텔 맥
  fi
fi
if [ ! -f "$bin_src" ]; then
  echo "실행파일을 찾을 수 없습니다: $bin_src" >&2
  echo "  make 로 빌드하거나, 실행파일 경로를 인자로 넘겨 주세요." >&2
  exit 1
fi

dest="$HOME/.local/bin/claude-awake"
mkdir -p "$HOME/.local/bin"
cp "$bin_src" "$dest"
chmod +x "$dest"

plist="$HOME/Library/LaunchAgents/com.claude-awake.plist"
cat > "$plist" <<EOF
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0"><dict>
  <key>Label</key><string>com.claude-awake</string>
  <key>ProgramArguments</key><array><string>$dest</string></array>
  <key>RunAtLoad</key><true/>
  <key>KeepAlive</key><true/>
  <key>StandardOutPath</key><string>/tmp/claude-awake.log</string>
  <key>StandardErrorPath</key><string>/tmp/claude-awake.log</string>
</dict></plist>
EOF

# 이미 등록돼 있어도 깨끗하게 교체되도록 내렸다가(bootout) 다시 올린다(bootstrap).
#   bootout 은 비동기라, 완전히 내려가기 전에 bootstrap 하면 "Input/output error(5)" 로 실패한다.
#   그래서 목록에서 사라질 때까지 잠깐 기다린다(재설치·업그레이드 때 실제로 겪는 문제).
domain="gui/$(id -u)"
launchctl bootout "$domain/com.claude-awake" 2>/dev/null || true
for _ in $(seq 1 50); do
  launchctl print "$domain/com.claude-awake" >/dev/null 2>&1 || break
  sleep 0.1
done
launchctl bootstrap "$domain" "$plist"

echo "설치했습니다."
echo "  이제 로그인할 때 자동으로 시작되고, Claude 앱이 켜져 있는 동안 절전을 막습니다."
echo "  실행파일: $dest"
echo "  로그:     /tmp/claude-awake.log"
echo "  제거하려면: bash install/macos-uninstall.sh"
