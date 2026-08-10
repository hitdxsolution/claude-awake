//go:build windows

package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"unsafe"
)

// SetThreadExecutionState 플래그. ES_CONTINUOUS 가 이 스레드에 걸려 있는 동안 윈도우는
// 컴퓨터(와 화면)를 재우지 않는다. Win32 API 를 직접 부르므로 뒤에서 도는 보조 프로그램이 없다.
const (
	esContinuous      = 0x80000000
	esSystemRequired  = 0x00000001
	esDisplayRequired = 0x00000002
)

// 메시지 창 플래그 — 일반 사용자에게 결과를 알려 주는 유일한 창구다(콘솔이 없으므로).
const (
	mbOK              = 0x00000000
	mbYesNo           = 0x00000004
	mbIconInformation = 0x00000040
	mbIconQuestion    = 0x00000020
	idYes             = 6
)

var (
	kernel32                    = syscall.NewLazyDLL("kernel32.dll")
	procSetThreadExecutionState = kernel32.NewProc("SetThreadExecutionState")
	user32                      = syscall.NewLazyDLL("user32.dll")
	procMessageBoxW             = user32.NewProc("MessageBoxW")
)

// 설치 폴더 — 사용자 계정 영역이라 관리자 권한이 필요 없다.
func installDir() string {
	base := os.Getenv("LOCALAPPDATA")
	if base == "" {
		base = os.Getenv("USERPROFILE")
	}
	return filepath.Join(base, "claude-awake")
}

func installedPath() string { return filepath.Join(installDir(), "claude-awake.exe") }

// GUI 서브시스템으로 빌드해(-H windowsgui) 콘솔이 없으므로, 기록은 파일로 남긴다.
func openPlatformLog() func() {
	if err := os.MkdirAll(installDir(), 0o755); err != nil {
		return nil
	}
	file, err := os.OpenFile(filepath.Join(installDir(), "claude-awake.log"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return nil
	}
	logOut = io.Writer(file)
	return func() { _ = file.Close() }
}

// 창을 띄우지 않고 명령을 실행한다(폴링 때마다 검은 창이 깜빡이지 않도록).
func hiddenCommand(name string, args ...string) *exec.Cmd {
	cmd := exec.Command(name, args...)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	return cmd
}

func isClaudeRunning() bool {
	out, err := hiddenCommand("tasklist", "/FI", "IMAGENAME eq Claude.exe", "/NH").Output()
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
	// ES_CONTINUOUS 만 걸면 앞서 설정한 요구사항이 해제된다.
	_, _, _ = procSetThreadExecutionState.Call(esContinuous)
}

// ─────────────────────────── 설치 / 제거 ───────────────────────────

const taskName = "claude-awake"

// PowerShell 을 창 없이 실행한다. 작업 스케줄러 등록은 이 방법이 가장 안정적이다.
func powershell(script string) error {
	cmd := hiddenCommand("powershell", "-NoProfile", "-NonInteractive", "-ExecutionPolicy", "Bypass", "-Command", script)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%w: %s", err, strings.TrimSpace(string(out)))
	}
	return nil
}

func isInstalled() bool {
	if _, err := os.Stat(installedPath()); err != nil {
		return false
	}
	return powershell("Get-ScheduledTask -TaskName '"+taskName+"' -ErrorAction Stop | Out-Null") == nil
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

	// 로그인할 때 시작하고, 죽으면 다시 뜨고, 실행 시간 제한은 두지 않는다.
	//   내 계정 작업으로 등록하므로 관리자 권한(UAC)이 필요 없다.
	script := strings.Join([]string{
		"$a=New-ScheduledTaskAction -Execute '" + installedPath() + "';",
		"$t=New-ScheduledTaskTrigger -AtLogOn;",
		"$s=New-ScheduledTaskSettingsSet -AllowStartIfOnBatteries -DontStopIfGoingOnBatteries",
		"-RestartCount 999 -RestartInterval (New-TimeSpan -Minutes 1) -ExecutionTimeLimit ([TimeSpan]::Zero);",
		"Register-ScheduledTask -TaskName '" + taskName + "' -Action $a -Trigger $t -Settings $s -Force | Out-Null;",
		"Start-ScheduledTask -TaskName '" + taskName + "'",
	}, " ")
	return powershell(script)
}

func uninstall() error {
	_ = powershell("Stop-ScheduledTask -TaskName '" + taskName + "' -ErrorAction SilentlyContinue;" +
		"Unregister-ScheduledTask -TaskName '" + taskName + "' -Confirm:$false -ErrorAction SilentlyContinue")
	// 실행 중인 자기 자신은 지울 수 없으므로, 지워지지 않아도 실패로 보지 않는다.
	_ = os.Remove(installedPath())
	return nil
}

// ─────────────────────────── 사용자 알림 ───────────────────────────

func messageBox(title, text string, flags uintptr) int {
	textPtr, err := syscall.UTF16PtrFromString(text)
	if err != nil {
		return 0
	}
	titlePtr, err := syscall.UTF16PtrFromString(title)
	if err != nil {
		return 0
	}
	ret, _, _ := procMessageBoxW.Call(0, uintptr(unsafe.Pointer(textPtr)), uintptr(unsafe.Pointer(titlePtr)), flags)
	return int(ret)
}

func showMessage(title, text string) {
	messageBox(title, text, mbOK|mbIconInformation)
}

func askYesNo(title, text string) bool {
	return messageBox(title, text, mbYesNo|mbIconQuestion) == idYes
}
