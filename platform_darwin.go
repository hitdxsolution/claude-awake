//go:build darwin

package main

import "os/exec"

// The macOS inhibitor is a child process: `caffeinate -dimsu` blocks display/idle/system/disk
// sleep and marks the user active. It lives until we kill it, which is exactly the window we want.
var caffeinate *exec.Cmd

// openPlatformLog: stdout is the log channel here — the LaunchAgent captures it to a file.
func openPlatformLog() func() { return nil }

func isClaudeRunning() bool {
	// The main app executable is .../Claude.app/Contents/MacOS/Claude. Matching the process
	// (not a fixed install path) means it works wherever the user installed Claude.
	return exec.Command("pgrep", "-f", "Claude.app/Contents/MacOS/Claude").Run() == nil
}

func blockSleep() error {
	cmd := exec.Command("caffeinate", "-dimsu")
	if err := cmd.Start(); err != nil {
		return err
	}
	caffeinate = cmd
	return nil
}

func releaseSleepBlock() {
	if caffeinate == nil {
		return
	}
	_ = caffeinate.Process.Kill()
	_ = caffeinate.Wait() // reap it so we do not leave a zombie
	caffeinate = nil
}
