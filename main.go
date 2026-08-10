// claude-awake — Claude 데스크톱 앱이 켜져 있는 동안 컴퓨터가 잠들지 않게 한다.
//
//	Windows 를 주 대상으로 하고 macOS 도 지원한다. 의존성이 없는 단일 실행파일이라
//	설치하는 PC 에는 아무것도 미리 깔 필요가 없다.
//
//	동작: 몇 초마다 Claude 앱이 떠 있는지 확인해서, 떠 있으면 운영체제의 절전 차단을 걸고
//	앱이 꺼지면 바로 푼다. 운영체제별 부분(감지·차단·설치)은 platform_*.go 에 있다.
//
//	실행 방식이 두 가지인데, **자기 경로로 구분**한다 — 사용자는 더블클릭만 하면 된다:
//	  · 설치 위치 밖에서 실행됨(사용자가 내려받아 더블클릭) → 설치 프로그램으로 동작
//	  · 설치 위치에서 실행됨(로그인 시 OS 가 자동 실행)     → 감시 프로그램으로 동작
package main

import (
	"fmt"
	"io"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"time"
)

// Claude 가 떠 있는지 확인하는 주기.
const pollInterval = 15 * time.Second

var logOut io.Writer = os.Stdout

func logf(format string, args ...any) {
	fmt.Fprintf(logOut, "[claude-awake] %s %s\n", time.Now().UTC().Format(time.RFC3339), fmt.Sprintf(format, args...))
}

// 지금 실행 중인 파일이 "설치된 그 파일"인지. 대소문자·경로 표기 차이를 흡수한다.
func runningFromInstallPath() bool {
	exe, err := os.Executable()
	if err != nil {
		return false
	}
	exe, err = filepath.EvalSymlinks(exe)
	if err != nil {
		return false
	}
	target, err := filepath.EvalSymlinks(installedPath())
	if err != nil {
		return false // 아직 설치 전이면 설치 위치 파일이 없다
	}
	if runtime.GOOS == "windows" {
		return strings.EqualFold(exe, target)
	}
	return exe == target
}

// 감시 루프 — 설치된 뒤 로그인할 때마다 OS 가 이걸 실행한다.
func runWatchdog() {
	// Windows 는 절전 차단이 "그 스레드"에 걸리므로 고루틴을 한 스레드에 고정한다.
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	if closeLog := openPlatformLog(); closeLog != nil {
		defer closeLog()
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	blocking := false
	// 우리가 종료될 때 절전 차단을 걸어 둔 채로 두지 않는다.
	defer func() {
		if blocking {
			releaseSleepBlock()
		}
	}()

	// 한 번의 확인 — Claude 실행 여부에 차단 상태를 맞춘다.
	check := func() {
		running := isClaudeRunning()
		switch {
		case running && !blocking:
			if err := blockSleep(); err != nil {
				logf("절전 차단에 실패했습니다: %v", err)
				return
			}
			blocking = true
			logf("Claude 실행 감지 → 절전을 막습니다")
		case !running && blocking:
			releaseSleepBlock()
			blocking = false
			logf("Claude 종료 → 절전을 다시 허용합니다")
		}
	}

	logf("시작됨 — 플랫폼=%s, 확인주기=%s", runtime.GOOS, pollInterval)
	check()

	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			check()
		case <-stop:
			logf("종료합니다")
			return
		}
	}
}

// 설치 프로그램 — 사용자가 내려받은 파일을 더블클릭했을 때의 동작.
//
//	이미 설치돼 있으면 제거할지 물어본다. 그래서 파일 하나로 설치와 제거가 모두 된다.
func runInstaller() {
	if isInstalled() {
		if !askYesNo("claude-awake", "이미 설치되어 있습니다.\n\n제거할까요?") {
			return
		}
		if err := uninstall(); err != nil {
			showMessage("claude-awake — 제거 실패", "제거하지 못했습니다.\n\n"+err.Error())
			return
		}
		showMessage("claude-awake", "제거했습니다.\n\n절전 설정이 원래대로 돌아갑니다.")
		return
	}

	if err := install(); err != nil {
		showMessage("claude-awake — 설치 실패", "설치하지 못했습니다.\n\n"+err.Error())
		return
	}
	showMessage("claude-awake", "설치했습니다.\n\n"+
		"이제 컴퓨터를 켤 때 자동으로 시작되고,\n"+
		"Claude 앱이 켜져 있는 동안에는 컴퓨터가 잠들지 않습니다.\n\n"+
		"끄고 싶으면 이 파일을 다시 실행하세요.")
}

func main() {
	if runningFromInstallPath() {
		runWatchdog()
		return
	}
	runInstaller()
}
