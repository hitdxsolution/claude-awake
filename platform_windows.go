//go:build windows

package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
)

// SetThreadExecutionState flags. While ES_CONTINUOUS is set on this thread the OS keeps the
// machine (and display) awake; clearing it restores normal power behaviour. Calling the Win32
// API directly means no PowerShell helper and no extra process to babysit.
const (
	esContinuous      = 0x80000000
	esSystemRequired  = 0x00000001
	esDisplayRequired = 0x00000002
)

var (
	kernel32                    = syscall.NewLazyDLL("kernel32.dll")
	procSetThreadExecutionState = kernel32.NewProc("SetThreadExecutionState")
)

// openPlatformLog: the binary is built for the GUI subsystem (no console window), so stdout
// goes nowhere — write the log next to the installed executable instead.
func openPlatformLog() func() {
	dir := os.Getenv("LOCALAPPDATA")
	if dir == "" {
		return nil
	}
	dir = filepath.Join(dir, "claude-awake")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil
	}
	file, err := os.OpenFile(filepath.Join(dir, "claude-awake.log"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return nil
	}
	logOut = file
	return func() { _ = file.Close() }
}

func isClaudeRunning() bool {
	cmd := exec.Command("tasklist", "/FI", "IMAGENAME eq Claude.exe", "/NH")
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true} // no console flash on each poll
	out, err := cmd.Output()
	if err != nil {
		return false
	}
	return strings.Contains(strings.ToLower(string(out)), "claude.exe")
}

func blockSleep() error {
	if ret, _, err := procSetThreadExecutionState.Call(esContinuous | esSystemRequired | esDisplayRequired); ret == 0 {
		return err
	}
	return nil
}

func releaseSleepBlock() {
	// ES_CONTINUOUS on its own clears the requirements previously set on this thread.
	_, _, _ = procSetThreadExecutionState.Call(esContinuous)
}
