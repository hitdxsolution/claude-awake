#!/bin/bash
# Remove the claude-awake LaunchAgent and binary.
set -euo pipefail
launchctl bootout "gui/$(id -u)/com.claude-awake" 2>/dev/null || true
rm -f "$HOME/Library/LaunchAgents/com.claude-awake.plist"
rm -f "$HOME/.local/bin/claude-awake"
echo "Uninstalled."
