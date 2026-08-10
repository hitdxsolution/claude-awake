# claude-awake

Keep your computer awake while the **Claude desktop app** is running — and let it sleep
normally as soon as you quit Claude.

Cross-platform (**macOS & Windows**). A single ~2 MB static binary with zero dependencies, so
the machine you install it on needs **no runtime** — no Node, no Python, nothing.

## Why

Long Claude sessions and background agents get interrupted when the machine goes to sleep.
`claude-awake` blocks sleep only while the Claude app is open, then gets out of the way.

## How it works

A tiny watchdog polls every 15 seconds for the Claude app. While it's running, it holds an
**OS built-in** sleep inhibitor; when Claude quits, it releases it.

| Step | macOS | Windows |
| --- | --- | --- |
| Detect the app | `pgrep` for the Claude process | `tasklist` for `Claude.exe` |
| Block sleep while open | `caffeinate -dimsu` | Win32 `SetThreadExecutionState` |
| Start at login + auto-restart | LaunchAgent | Scheduled Task (at logon) |

The app is detected by **process name**, not a fixed install path, so it works wherever Claude
is installed. On Windows the Win32 API is called directly — no PowerShell helper process.

## Install

The binaries are committed in [`dist/`](dist), so you can download this repo (Code → Download
ZIP, or `git clone`) and install straight away — nothing to build.

### macOS

```bash
bash install/macos-install.sh
```

Picks the right binary for your Mac (Apple Silicon or Intel), then registers a LaunchAgent that
starts at login and restarts if it dies. Log: `/tmp/claude-awake.log`.

> First launch: macOS Gatekeeper may block an unsigned binary. Right-click the binary →
> **Open** once to allow it (or sign/notarize for wider distribution).

Uninstall: `bash install/macos-uninstall.sh`

### Windows

```powershell
powershell -ExecutionPolicy Bypass -File .\install\windows-install.ps1
```

Registers a Scheduled Task that runs at logon and restarts on failure. The binary is built for
the GUI subsystem, so it never flashes a console window.
Log: `%LOCALAPPDATA%\claude-awake\claude-awake.log`.

> First launch: SmartScreen may warn about an unknown publisher → **More info → Run anyway**
> (or code-sign for wider distribution).

Uninstall: `powershell -ExecutionPolicy Bypass -File .\install\windows-uninstall.ps1`

## Build from source

Requires [Go](https://go.dev) **on the build machine only**.

```bash
make          # gofmt + go vet, then build all three binaries into dist/
make check    # formatting + static analysis only
go test ./... # verifies the inhibitor actually starts and leaves nothing behind
```

Go cross-compiles, so `make` on a Mac produces the Windows `.exe` too. Individual targets:
`make mac-arm64`, `make mac-x64`, `make windows-x64`.

Binaries are built with `-trimpath -s -w` (and `-H windowsgui` on Windows), which is why they
land at roughly 2 MB each.

## Run without installing

```bash
./dist/claude-awake-macos-arm64   # Ctrl+C to stop; it releases the block on exit
```

## Configuration

Edit `pollInterval` at the top of `main.go` and rebuild.

## Layout

| File | Purpose |
| --- | --- |
| `main.go` | The watchdog loop — platform-agnostic |
| `platform_darwin.go` | macOS detection + `caffeinate` inhibitor |
| `platform_windows.go` | Windows detection + `SetThreadExecutionState` |
| `platform_linux.go` | Best-effort Linux (`systemd-inhibit`) |
| `install/` | Autostart registration per OS |
| `dist/` | Prebuilt binaries (committed) |

---

### 한국어 요약

Claude 데스크톱 앱이 **켜져 있는 동안만** 컴퓨터가 절전에 들지 않게 막아 주고, 앱을 끄면
절전이 자동 복귀합니다. **macOS·Windows 공용**이며, 설치하는 PC엔 **런타임 설치가 전혀 필요
없습니다** — 2MB짜리 실행파일 하나가 전부입니다.

- 바이너리가 `dist/`에 **커밋돼 있어** repo만 받으면 바로 설치됩니다(빌드 불필요)
- 감지: 설치 경로가 아니라 **프로세스 이름**으로 → 어디에 깔려 있든 동작
- 절전 차단: mac `caffeinate`, Windows `SetThreadExecutionState`(Win32 API 직접 호출)
- 자동시작: mac LaunchAgent, Windows 작업 스케줄러(로그인 시)
- 빌드는 개발자 PC에서만 Go 필요(`make`) — Mac에서 Windows exe까지 크로스컴파일됩니다

## License

MIT
