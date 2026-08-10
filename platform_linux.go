//go:build linux

package main

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
)

// 리눅스는 참고용 지원이다. systemd-inhibit 은 자식이 사는 동안 차단을 유지하므로
// `sleep infinity` 를 물려 두고, 통째로 죽여서 해제한다.
var inhibit *exec.Cmd

func installDir() string { return filepath.Join(os.Getenv("HOME"), ".local", "bin") }

func installedPath() string { return filepath.Join(installDir(), "claude-awake") }

// 기록은 stdout 으로 — systemd user 유닛이 받아 준다.
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

// 자동 시작 등록은 배포판마다 방식이 달라 자동화하지 않는다 — 직접 실행하거나 systemd 유닛을 만들면 된다.
func isInstalled() bool { return false }

func install() error {
	return errors.New("리눅스는 자동 설치를 지원하지 않습니다. 실행파일을 직접 실행하거나 systemd 유닛으로 등록해 주세요")
}

func uninstall() error { return nil }

func showMessage(title, text string) {
	logf("%s: %s", title, text)
}

func askYesNo(_, _ string) bool { return false }
