# claude-awake

**Claude 데스크톱 앱이 켜져 있는 동안 컴퓨터가 잠들지 않게** 해 줍니다. 앱을 끄면 절전이 자동으로 다시 켜집니다.

macOS·Windows 공용이고, **설치하는 PC에는 아무것도 미리 깔 필요가 없습니다.** 2MB짜리 실행파일 하나가 전부입니다.

---

## 설치

이 저장소를 받으면(초록색 **Code** 버튼 → **Download ZIP**, 또는 `git clone`) **바로 설치할 수 있습니다.**
실행파일이 [`dist/`](dist) 에 들어 있어서 따로 빌드할 필요가 없습니다.

### macOS

터미널에서 받은 폴더로 이동한 뒤:

```bash
bash install/macos-install.sh
```

내 맥에 맞는 실행파일(애플 실리콘/인텔)을 알아서 골라 설치하고, 로그인할 때 자동으로 켜지도록 등록합니다.

> **처음 실행할 때 "확인되지 않은 개발자" 경고가 뜨면**
> `dist` 폴더에서 실행파일을 **우클릭 → 열기** 를 한 번 눌러 주면 이후로는 그냥 실행됩니다.
> (개발자 서명이 없는 프로그램에 macOS 가 띄우는 기본 경고입니다.)

제거: `bash install/macos-uninstall.sh`

### Windows

**관리자 권한은 필요 없습니다.** 내 계정에만 등록되므로 UAC(사용자 계정 컨트롤) 창도 뜨지 않습니다.

#### 1단계 — 파일 준비

받은 폴더를 **통째로** 원하는 위치에 두세요(예: `C:\Users\내계정\Downloads\claude-awake`).
설치 스크립트가 `..\dist\claude-awake-windows-x64.exe` 를 상대경로로 찾기 때문에, **exe 하나만 빼서 옮기면 경로를 따로 알려 줘야 합니다**(아래 "문제가 생기면" 참고).

#### 2단계 — PowerShell 열기

시작 메뉴에서 **PowerShell** 검색 → 그냥 클릭해서 실행합니다(관리자로 실행할 필요 없음).

#### 3단계 — 폴더로 이동

```powershell
cd "C:\Users\내계정\Downloads\claude-awake"
```

> 경로가 헷갈리면, 탐색기에서 그 폴더를 연 뒤 주소창을 클릭해 경로를 복사해서 `cd "붙여넣기"` 하시면 됩니다.

#### 4단계 — (ZIP 으로 받았다면) 차단 해제

인터넷에서 받은 ZIP 안의 파일은 윈도우가 "차단됨" 표시를 붙여 둡니다. 한 번만 풀어 주세요.

```powershell
Get-ChildItem -Recurse | Unblock-File
```

#### 5단계 — 설치

```powershell
powershell -ExecutionPolicy Bypass -File .\install\windows-install.ps1
```

이렇게 나오면 성공입니다.

```
설치했습니다.
  이제 로그인할 때 자동으로 시작되고, Claude 앱이 켜져 있는 동안 절전을 막습니다.
  실행파일: C:\Users\내계정\AppData\Local\claude-awake\claude-awake.exe
  로그:     C:\Users\내계정\AppData\Local\claude-awake\claude-awake.log
  제거하려면: powershell -ExecutionPolicy Bypass -File .\install\windows-uninstall.ps1
```

#### 6단계 — 잘 되는지 확인

```powershell
# 작업이 등록되고 실행 중인지
Get-ScheduledTask -TaskName claude-awake | Select-Object TaskName, State

# 로그 실시간으로 보기 (Ctrl+C 로 중단)
Get-Content "$env:LOCALAPPDATA\claude-awake\claude-awake.log" -Wait -Tail 5
```

Claude 앱을 켜면 로그에 `Claude 실행 감지 → 절전을 막습니다`, 끄면 `Claude 종료 → 절전을 다시 허용합니다`
가 찍힙니다(최대 15초 안에). 작업 관리자 → 세부 정보 탭에서 `claude-awake.exe` 로도 확인할 수 있습니다.

> **검은 콘솔 창은 뜨지 않습니다.** 백그라운드 전용으로 빌드해서(`-H windowsgui`) 화면에 아무것도 안 보이는 게
> 정상입니다. 그래서 실행 기록은 화면이 아니라 위 로그 파일에 쌓입니다.

#### 문제가 생기면

| 증상 | 해결 |
| --- | --- |
| `실행파일을 찾을 수 없습니다` | 폴더째 옮기지 않았을 때입니다. exe 경로를 직접 알려 주세요: `powershell -ExecutionPolicy Bypass -File .\windows-install.ps1 -BinSrc "C:\경로\claude-awake-windows-x64.exe"` |
| 스크립트가 실행되지 않음 | 4단계(`Unblock-File`)를 먼저 하고 다시 시도해 주세요 |
| 백신이 차단하거나 파일을 지움 | 서명이 없는 프로그램이라 생기는 오탐입니다. 백신에서 `%LOCALAPPDATA%\claude-awake` 폴더를 예외로 등록해 주세요 |
| 한글이 깨져서 보임 | PowerShell 창에서 `chcp 65001` 을 먼저 실행한 뒤 다시 설치해 주세요 |
| 로그가 안 생김 | 작업이 등록만 되고 시작이 안 된 경우입니다: `Start-ScheduledTask -TaskName claude-awake` |

> **"Windows의 PC 보호"(SmartScreen) 창은 보통 뜨지 않습니다.**
> 설치 스크립트가 exe 를 직접 실행하지 않고 작업 스케줄러에 등록해서 실행시키기 때문입니다.
> 다만 exe 를 **탐색기에서 직접 더블클릭**하면 뜰 수 있는데, 그때는 **추가 정보 → 실행** 을 누르시면 됩니다.

#### 제거

```powershell
powershell -ExecutionPolicy Bypass -File .\install\windows-uninstall.ps1
```

작업 스케줄러 등록과 설치된 파일을 지웁니다. 절전 설정은 윈도우의 원래 값으로 돌아갑니다.

### 잘 되고 있는지 확인 (공통)

로그 파일에 아래처럼 찍히면 정상입니다.

- macOS: `/tmp/claude-awake.log` — 실시간으로 보려면 `tail -f /tmp/claude-awake.log`
- Windows: `%LOCALAPPDATA%\claude-awake\claude-awake.log`

```
[claude-awake] 2026-08-10T07:37:57Z 시작됨 — 플랫폼=darwin, 확인주기=15s
[claude-awake] 2026-08-10T07:37:57Z Claude 실행 감지 → 절전을 막습니다
```

macOS 에서는 절전이 실제로 막혔는지 운영체제에 직접 물어볼 수도 있습니다.

```bash
pmset -g assertions | grep PreventUserIdleSystemSleep   # 값이 1 이면 차단 중
```

---

## 왜 필요한가

Claude 를 오래 켜 두거나 백그라운드 작업을 돌리는 중에 컴퓨터가 절전에 들어가면 작업이 끊깁니다.
그렇다고 절전을 아예 꺼 두면 평소에 배터리·전기를 낭비하게 됩니다.

`claude-awake` 는 **Claude 가 켜져 있는 동안에만** 절전을 막고, 앱을 끄면 곧바로 손을 뗍니다.

## 어떻게 동작하나

15초마다 Claude 앱이 켜져 있는지 확인해서, 켜져 있으면 **운영체제에 원래 있는 기능**으로 절전을 막고
앱이 꺼지면 바로 풀어 줍니다.

| 하는 일 | macOS | Windows |
| --- | --- | --- |
| 앱 감지 | `pgrep` 으로 Claude 프로세스 확인 | `tasklist` 로 `Claude.exe` 확인 |
| 절전 막기 | `caffeinate -dimsu` | Win32 `SetThreadExecutionState` |
| 로그인 시 자동 시작 | LaunchAgent | 작업 스케줄러 |

설치 **경로**가 아니라 **프로세스 이름**으로 찾기 때문에 Claude 를 어디에 설치했든 동작합니다.
Windows 에서는 Win32 API 를 직접 호출해서, 뒤에서 도는 보조 프로그램이 없습니다.

## 직접 빌드하기

빌드하는 컴퓨터에만 [Go](https://go.dev) 가 필요합니다(설치하는 PC에는 필요 없습니다).

```bash
make          # 포맷·정적검사 후 세 가지 실행파일을 dist/ 에 빌드
make check    # 포맷·정적검사만
go test ./... # 절전 차단이 실제로 걸리고, 해제 후 찌꺼기가 남지 않는지 검증
```

Go 는 크로스컴파일이 되므로 **맥에서 윈도우용 exe 까지** 한 번에 만들어집니다.
개별 빌드: `make mac-arm64`, `make mac-x64`, `make windows-x64`.

## 설정 바꾸기

`main.go` 맨 위의 `pollInterval`(확인 주기)을 고친 뒤 다시 빌드하면 됩니다.

## 파일 구성

| 파일 | 역할 |
| --- | --- |
| `main.go` | 감시 루프 — 운영체제와 무관한 공통 부분 |
| `platform_darwin.go` | macOS 감지 + `caffeinate` |
| `platform_windows.go` | Windows 감지 + `SetThreadExecutionState` |
| `platform_linux.go` | 리눅스(참고용, `systemd-inhibit`) |
| `install/` | 운영체제별 자동 시작 등록 |
| `dist/` | 미리 빌드해 둔 실행파일 |

---

### In English

Keeps your computer awake while the **Claude desktop app** is running, and lets it sleep again as
soon as you quit Claude. macOS & Windows, a single ~2 MB static binary with no runtime required on
the target machine.

Prebuilt binaries are committed in `dist/`, so you can clone (or download the ZIP) and install with
`bash install/macos-install.sh` on macOS, or
`powershell -ExecutionPolicy Bypass -File .\install\windows-install.ps1` on Windows.
Messages are in Korean. Build from source with `make` (requires Go on the build machine only).

## License

MIT
