// claude-awake — keep the computer awake while the Claude desktop app is running.
//
//	Cross-platform (macOS · Windows · best-effort Linux) and dependency-free: the whole
//	thing is one static binary, so the target machine needs no runtime at all.
//
//	How it works: poll every few seconds for the Claude desktop app. While it is running,
//	hold an OS sleep inhibitor; release it the moment the app quits. The per-OS pieces
//	(detection + inhibitor) live in platform_*.go behind two small functions.
package main

import (
	"fmt"
	"io"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"
)

// How often to check whether Claude is running.
const pollInterval = 15 * time.Second

var logOut io.Writer = os.Stdout

func logf(format string, args ...any) {
	fmt.Fprintf(logOut, "[claude-awake] %s %s\n", time.Now().UTC().Format(time.RFC3339), fmt.Sprintf(format, args...))
}

func main() {
	// Windows holds the sleep block on one specific OS thread, so keep this goroutine pinned
	// to a single thread for the whole run. Harmless on the other platforms.
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	if closeLog := openPlatformLog(); closeLog != nil {
		defer closeLog()
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	blocking := false
	// Never leave the machine caffeinated once we exit.
	defer func() {
		if blocking {
			releaseSleepBlock()
		}
	}()

	// One check: match the inhibitor to whether Claude is running.
	check := func() {
		running := isClaudeRunning()
		switch {
		case running && !blocking:
			if err := blockSleep(); err != nil {
				logf("could not block sleep: %v", err)
				return
			}
			blocking = true
			logf("Claude detected → blocking sleep")
		case !running && blocking:
			releaseSleepBlock()
			blocking = false
			logf("Claude closed → sleep allowed")
		}
	}

	logf("started — platform=%s, poll=%s", runtime.GOOS, pollInterval)
	check()

	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			check()
		case <-stop:
			logf("stopping")
			return
		}
	}
}
