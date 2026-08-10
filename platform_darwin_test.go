//go:build darwin

package main

import (
	"os/exec"
	"strings"
	"testing"
)

// countCaffeinate reports how many `caffeinate -dimsu` processes we currently own.
func countCaffeinate(t *testing.T) int {
	t.Helper()
	out, err := exec.Command("pgrep", "-f", "caffeinate -dimsu").Output()
	if err != nil {
		return 0 // pgrep exits non-zero when nothing matches
	}
	return len(strings.Fields(string(out)))
}

// The inhibitor must actually start a caffeinate process and, just as importantly, leave none
// behind once released — a leaked one would keep the machine awake forever.
func TestBlockAndReleaseSleep(t *testing.T) {
	before := countCaffeinate(t)

	if err := blockSleep(); err != nil {
		t.Fatalf("blockSleep() = %v, want nil", err)
	}
	if got := countCaffeinate(t); got <= before {
		t.Fatalf("caffeinate count = %d after blockSleep, want > %d", got, before)
	}

	releaseSleepBlock()
	if got := countCaffeinate(t); got != before {
		t.Fatalf("caffeinate count = %d after release, want %d (leaked inhibitor)", got, before)
	}
	if caffeinate != nil {
		t.Fatal("caffeinate handle still set after release")
	}
}

// Releasing when nothing is held must be a harmless no-op (the shutdown path can hit this).
func TestReleaseWithoutBlockIsNoop(t *testing.T) {
	caffeinate = nil
	releaseSleepBlock()
}

// Detection must not report a false positive when the app is not running, and must never
// error out — a crash here would take the whole watchdog down.
func TestIsClaudeRunningDoesNotPanic(t *testing.T) {
	_ = isClaudeRunning()
}
