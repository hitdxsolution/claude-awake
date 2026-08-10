//go:build linux

package main

import "os/exec"

// Best-effort Linux support: systemd-inhibit holds the block for as long as its child lives,
// so we give it `sleep infinity` and kill the whole thing to release.
var inhibit *exec.Cmd

// openPlatformLog: stdout is the log channel here — a systemd user unit captures it.
func openPlatformLog() func() { return nil }

func isClaudeRunning() bool {
	return exec.Command("pgrep", "-x", "claude").Run() == nil
}

func blockSleep() error {
	cmd := exec.Command("systemd-inhibit", "--what=idle:sleep", "--who=claude-awake", "--why=Claude running", "sleep", "infinity")
	if err := cmd.Start(); err != nil {
		return err
	}
	inhibit = cmd
	return nil
}

func releaseSleepBlock() {
	if inhibit == nil {
		return
	}
	_ = inhibit.Process.Kill()
	_ = inhibit.Wait()
	inhibit = nil
}
