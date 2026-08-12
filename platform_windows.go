//go:build windows

package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"unicode/utf16"
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
//
//	%LOCALAPPDATA%\Programs 아래에 두는 것은 VS Code·Slack 같은 정상 앱이 쓰는 관례다.
//	%LOCALAPPDATA% 바로 밑은 멀웨어가 즐겨 쓰는 자리라 백신 휴리스틱 점수가 그만큼 높다.
func installDir() string {
	base := os.Getenv("LOCALAPPDATA")
	if base == "" {
		base = os.Getenv("USERPROFILE")
	}
	return filepath.Join(base, "Programs", "claude-awake")
}

// 0.1 버전이 쓰던 설치 폴더. 그때 깔았던 사람이 새 파일을 실행하면 여기를 치운다.
func legacyInstallDir() string {
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

// 작업 스케줄러 등록은 schtasks.exe 로 한다.
//
//	예전에는 PowerShell 의 Register-ScheduledTask 를 썼는데, `powershell -ExecutionPolicy Bypass` 를
//	창을 숨긴 채 실행하는 조합은 백신 휴리스틱이 가장 강하게 반응하는 패턴이다. 서명이 없는 이
//	프로그램이 그 때문에 오탐으로 격리되는 일이 있었다. schtasks.exe 는 윈도우 기본 도구이고,
//	XML 을 넘기면 배터리·실행시간·재시작까지 PowerShell 로 주던 설정을 그대로 줄 수 있다.
func schtasks(args ...string) error {
	out, err := hiddenCommand("schtasks", args...).CombinedOutput()
	if err != nil {
		return fmt.Errorf("%w: %s", err, strings.TrimSpace(string(out)))
	}
	return nil
}

// 등록할 작업 정의. 요소 순서는 작업 스케줄러가 내보내는 XML 과 같게 맞춰 두었다.
//
//	ExecutionTimeLimit PT0S = 시간 제한 없음, RestartOnFailure = 죽으면 1 분 뒤 다시 띄움.
const taskXMLTemplate = `<?xml version="1.0" encoding="UTF-16"?>
<Task version="1.2" xmlns="http://schemas.microsoft.com/windows/2004/02/mit/task">
  <RegistrationInfo>
    <Description>Claude 앱이 켜져 있는 동안 컴퓨터가 잠들지 않게 합니다.</Description>
  </RegistrationInfo>
  <Triggers>
    <LogonTrigger>
      <Enabled>true</Enabled>
      <UserId>%[1]s</UserId>
    </LogonTrigger>
  </Triggers>
  <Principals>
    <Principal id="Author">
      <UserId>%[1]s</UserId>
      <LogonType>InteractiveToken</LogonType>
      <RunLevel>LeastPrivilege</RunLevel>
    </Principal>
  </Principals>
  <Settings>
    <MultipleInstancesPolicy>IgnoreNew</MultipleInstancesPolicy>
    <DisallowStartIfOnBatteries>false</DisallowStartIfOnBatteries>
    <StopIfGoingOnBatteries>false</StopIfGoingOnBatteries>
    <AllowHardTerminate>true</AllowHardTerminate>
    <StartWhenAvailable>true</StartWhenAvailable>
    <RunOnlyIfNetworkAvailable>false</RunOnlyIfNetworkAvailable>
    <IdleSettings>
      <StopOnIdleEnd>false</StopOnIdleEnd>
      <RestartOnIdle>false</RestartOnIdle>
    </IdleSettings>
    <AllowStartOnDemand>true</AllowStartOnDemand>
    <Enabled>true</Enabled>
    <Hidden>false</Hidden>
    <RunOnlyIfIdle>false</RunOnlyIfIdle>
    <WakeToRun>false</WakeToRun>
    <ExecutionTimeLimit>PT0S</ExecutionTimeLimit>
    <Priority>7</Priority>
    <RestartOnFailure>
      <Interval>PT1M</Interval>
      <Count>999</Count>
    </RestartOnFailure>
  </Settings>
  <Actions Context="Author">
    <Exec>
      <Command>%[2]s</Command>
    </Exec>
  </Actions>
</Task>
`

// 작업을 등록할 계정. 내 계정 작업으로 등록해야 관리자 권한(UAC)이 필요 없다.
func currentUserID() string {
	user := os.Getenv("USERNAME")
	domain := os.Getenv("USERDOMAIN")
	if domain == "" {
		domain = os.Getenv("COMPUTERNAME")
	}
	if domain == "" || user == "" {
		return user
	}
	return domain + `\` + user
}

func escapeXML(text string) string {
	var buf strings.Builder
	_ = xml.EscapeText(&buf, []byte(text))
	return buf.String()
}

// schtasks /XML 은 UTF-16 파일만 확실히 받아 준다. BOM 을 붙여 리틀엔디언으로 쓴다.
func writeUTF16File(path, text string) error {
	units := utf16.Encode([]rune(text))
	buf := make([]byte, 0, len(units)*2+2)
	buf = append(buf, 0xFF, 0xFE) // BOM
	for _, unit := range units {
		buf = append(buf, byte(unit), byte(unit>>8))
	}
	return os.WriteFile(path, buf, 0o644)
}

func isInstalled() bool {
	if _, err := os.Stat(installedPath()); err != nil {
		return false
	}
	return schtasks("/Query", "/TN", taskName) == nil
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

	xmlPath := filepath.Join(installDir(), "task.xml")
	definition := fmt.Sprintf(taskXMLTemplate, escapeXML(currentUserID()), escapeXML(installedPath()))
	if err := writeUTF16File(xmlPath, definition); err != nil {
		return err
	}
	defer func() { _ = os.Remove(xmlPath) }()

	// /F 는 같은 이름의 기존 작업을 덮어쓴다 — 예전 버전 위에 다시 깔아도 깨끗하게 교체된다.
	if err := schtasks("/Create", "/TN", taskName, "/XML", xmlPath, "/F"); err != nil {
		return err
	}
	removeLegacyInstall()
	return schtasks("/Run", "/TN", taskName)
}

// 예전 설치 폴더가 남아 있으면 지운다. 없으면 아무 일도 하지 않는다.
func removeLegacyInstall() {
	if legacyInstallDir() == installDir() {
		return
	}
	_ = os.RemoveAll(legacyInstallDir())
}

func uninstall() error {
	_ = schtasks("/End", "/TN", taskName)
	_ = schtasks("/Delete", "/TN", taskName, "/F")
	// 실행 중인 자기 자신은 지울 수 없으므로, 지워지지 않아도 실패로 보지 않는다.
	_ = os.Remove(installedPath())
	removeLegacyInstall()
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
