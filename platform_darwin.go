//go:build darwin

package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// macOS 의 절전 차단은 자식 프로세스로 건다. `caffeinate -dimsu` 는 화면·시스템·디스크 절전을
// 모두 막고 사용자가 활동 중인 것으로 표시하며, 우리가 죽일 때까지 살아 있다.
var caffeinate *exec.Cmd

const launchAgentLabel = "com.claude-awake"

func installDir() string { return filepath.Join(os.Getenv("HOME"), ".local", "bin") }

func installedPath() string { return filepath.Join(installDir(), "claude-awake") }

func plistPath() string {
	return filepath.Join(os.Getenv("HOME"), "Library", "LaunchAgents", launchAgentLabel+".plist")
}

// 기록은 stdout 으로 — LaunchAgent 가 파일로 받아 준다.
func openPlatformLog() func() { return nil }

func isClaudeRunning() bool {
	// 앱 본체는 .../Claude.app/Contents/MacOS/Claude 다. 설치 경로가 아니라 프로세스로 찾기 때문에
	// 사용자가 Claude 를 어디에 두었든 동작한다.
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
	_ = caffeinate.Wait() // 좀비로 남지 않게 회수한다
	caffeinate = nil
}

// ─────────────────────────── 설치 / 제거 ───────────────────────────

func launchctlDomain() string { return "gui/" + strconv.Itoa(os.Getuid()) }

func isInstalled() bool {
	if _, err := os.Stat(installedPath()); err != nil {
		return false
	}
	return exec.Command("launchctl", "print", launchctlDomain()+"/"+launchAgentLabel).Run() == nil
}

func install() error {
	self, err := os.Executable()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(installDir(), 0o755); err != nil {
		return err
	}
	data, err := os.ReadFile(self)
	if err != nil {
		return err
	}
	if err := os.WriteFile(installedPath(), data, 0o755); err != nil {
		return err
	}

	plist := `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0"><dict>
  <key>Label</key><string>` + launchAgentLabel + `</string>
  <key>ProgramArguments</key><array><string>` + installedPath() + `</string></array>
  <key>RunAtLoad</key><true/>
  <key>KeepAlive</key><true/>
  <key>StandardOutPath</key><string>/tmp/claude-awake.log</string>
  <key>StandardErrorPath</key><string>/tmp/claude-awake.log</string>
</dict></plist>
`
	if err := os.MkdirAll(filepath.Dir(plistPath()), 0o755); err != nil {
		return err
	}
	if err := os.WriteFile(plistPath(), []byte(plist), 0o644); err != nil {
		return err
	}

	// 이미 등록돼 있어도 깨끗하게 교체한다. bootout 은 비동기라, 완전히 내려가기 전에 bootstrap 하면
	// "Input/output error(5)" 로 실패한다 — 목록에서 사라질 때까지 기다린다(재설치 때 실제로 겪는 문제).
	_ = exec.Command("launchctl", "bootout", launchctlDomain()+"/"+launchAgentLabel).Run()
	for i := 0; i < 50; i++ {
		if exec.Command("launchctl", "print", launchctlDomain()+"/"+launchAgentLabel).Run() != nil {
			break
		}
		time.Sleep(100 * time.Millisecond)
	}
	return exec.Command("launchctl", "bootstrap", launchctlDomain(), plistPath()).Run()
}

func uninstall() error {
	_ = exec.Command("launchctl", "bootout", launchctlDomain()+"/"+launchAgentLabel).Run()
	_ = os.Remove(plistPath())
	_ = os.Remove(installedPath())
	return nil
}

// ─────────────────────────── 사용자 알림 ───────────────────────────

// osascript 로 기본 대화상자를 띄운다 — 터미널을 몰라도 결과를 볼 수 있게.
func osascript(script string) (string, error) {
	out, err := exec.Command("osascript", "-e", script).Output()
	return strings.TrimSpace(string(out)), err
}

func escapeForAppleScript(text string) string {
	return strings.ReplaceAll(strings.ReplaceAll(text, `\`, `\\`), `"`, `\"`)
}

func showMessage(title, text string) {
	_, _ = osascript(`display dialog "` + escapeForAppleScript(text) + `" with title "` + escapeForAppleScript(title) +
		`" buttons {"확인"} default button 1`)
}

func askYesNo(title, text string) bool {
	out, err := osascript(`display dialog "` + escapeForAppleScript(text) + `" with title "` + escapeForAppleScript(title) +
		`" buttons {"아니오", "예"} default button 2`)
	return err == nil && strings.Contains(out, "예")
}
