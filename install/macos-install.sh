#!/bin/bash
# Install claude-awake as a macOS LaunchAgent: starts at login, restarts if it dies.
#   Usage: ./macos-install.sh [path-to-binary]
#   With no argument it uses ../dist/claude-awake-macos-<arch> from a local build.
set -euo pipefail

here="$(cd "$(dirname "$0")" && pwd)"
bin_src="${1:-}"
if [ -z "$bin_src" ]; then
  if [ "$(uname -m)" = "arm64" ]; then
    bin_src="$here/../dist/claude-awake-macos-arm64"
  else
    bin_src="$here/../dist/claude-awake-macos-x64"
  fi
fi
if [ ! -f "$bin_src" ]; then
  echo "Binary not found: $bin_src" >&2
  echo "Build it first:  bun run build      (or pass the binary path as an argument)" >&2
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

# Reload (bootout then bootstrap) so an existing agent is replaced cleanly.
launchctl bootout "gui/$(id -u)/com.claude-awake" 2>/dev/null || true
launchctl bootstrap "gui/$(id -u)" "$plist"

echo "Installed. It runs at login and while Claude is open."
echo "  binary: $dest"
echo "  log:    /tmp/claude-awake.log"
