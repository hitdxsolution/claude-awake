# claude-awake

**Claude 앱을 켜 두면 컴퓨터가 잠들지 않습니다.** Claude 를 끄면 절전이 원래대로 돌아옵니다.

Claude 로 오래 작업하거나 뭔가를 돌려 두는 중에 컴퓨터가 절전에 들어가 작업이 끊기는 것을 막아 줍니다.
절전 설정을 아예 꺼 둘 필요가 없습니다.

---

## 설치 (Windows)

### 1. 파일 받기

위쪽 초록색 **Code** 버튼 → **Download ZIP** 을 눌러 내려받고, 압축을 풉니다.

### 2. 실행파일 실행하기

압축을 푼 폴더 안 `dist` 폴더로 들어가서 **`claude-awake-windows-x64.exe` 를 더블클릭**합니다.

### 3. 끝

이런 창이 뜨면 설치가 끝난 것입니다. **확인**을 누르면 됩니다.

```
설치했습니다.

이제 컴퓨터를 켤 때 자동으로 시작되고,
Claude 앱이 켜져 있는 동안에는 컴퓨터가 잠들지 않습니다.

끄고 싶으면 이 파일을 다시 실행하세요.
```

이제 아무것도 안 하셔도 됩니다. 컴퓨터를 껐다 켜도 알아서 다시 시작됩니다.
평소에는 화면에 아무것도 보이지 않는 것이 정상입니다(뒤에서 조용히 동작합니다).

### 끄고 싶을 때

**같은 파일을 다시 더블클릭**하면 "제거할까요?" 라고 물어봅니다. **예**를 누르면 완전히 지워집니다.

### 창이 안 뜨거나 막힐 때

| 이런 창이 뜨면 | 이렇게 하세요 |
| --- | --- |
| **Windows의 PC 보호** (파란 창) | **추가 정보** → **실행** 을 누르세요. 만든 사람 서명이 없는 프로그램에 뜨는 안내입니다. |
| 백신이 막거나 파일이 사라짐 | 백신에서 이 파일을 예외로 등록해 주세요. 서명이 없어서 생기는 오탐입니다. |
| 아무 창도 안 뜸 | 압축을 **풀지 않고** ZIP 안에서 바로 실행하면 그럴 수 있습니다. 압축을 먼저 푸세요. |

---

## 설치 (macOS)

### 1. 파일 받기

**Code** → **Download ZIP** 으로 내려받고 압축을 풉니다.

### 2. 실행파일 실행하기

`dist` 폴더에서 내 맥에 맞는 파일을 **우클릭 → 열기** 로 실행합니다.

- **`claude-awake-macos-arm64`** — M1 이후 맥 (애플 실리콘)
- **`claude-awake-macos-x64`** — 그 이전 인텔 맥

> 처음에는 반드시 **우클릭 → 열기**로 실행하세요. 그냥 더블클릭하면
> "확인되지 않은 개발자" 경고만 뜨고 실행되지 않습니다. 한 번 이렇게 열어 두면 다음부터는 그냥 열립니다.
>
> 내 맥이 어떤 것인지 모르겠으면 화면 왼쪽 위 사과 메뉴 → **이 Mac에 관하여** 에서 확인할 수 있습니다.

### 3. 끝

"설치했습니다" 창이 뜨면 완료입니다. 끄고 싶으면 **같은 파일을 다시 실행**하면 제거할지 물어봅니다.

---

## 잘 되고 있는지 보고 싶다면

설치만 하면 신경 쓸 것이 없지만, 확인하고 싶다면 기록 파일을 열어 보면 됩니다.

- **Windows**: `C:\Users\내계정\AppData\Local\claude-awake\claude-awake.log`
- **macOS**: `/tmp/claude-awake.log`

Claude 를 켜면 `Claude 실행 감지 → 절전을 막습니다`, 끄면 `Claude 종료 → 절전을 다시 허용합니다`
가 기록됩니다(최대 15초 안에).

---

## 어떻게 동작하나

15초마다 Claude 앱이 켜져 있는지 확인해서, 켜져 있으면 **운영체제에 원래 있는 기능**으로 절전을 막고
앱이 꺼지면 바로 풀어 줍니다. 별도로 설치되는 것도, 계속 떠 있는 창도 없습니다.

| 하는 일 | Windows | macOS |
| --- | --- | --- |
| 앱 감지 | `tasklist` 로 `Claude.exe` 확인 | `pgrep` 으로 Claude 프로세스 확인 |
| 절전 막기 | Win32 `SetThreadExecutionState` | `caffeinate -dimsu` |
| 자동 시작 | 작업 스케줄러(내 계정) | LaunchAgent |

설치 **경로**가 아니라 **프로세스 이름**으로 찾기 때문에 Claude 를 어디에 설치했든 동작합니다.
**관리자 권한은 필요 없습니다** — 내 계정 영역에만 설치되므로 Windows 에서 UAC 창도 뜨지 않습니다.

실행파일 하나가 설치 프로그램이면서 본체입니다. 자기가 설치된 위치에서 실행되면 감시 프로그램으로,
그 밖에서 실행되면(=사용자가 더블클릭) 설치 프로그램으로 동작합니다.

## 개발자용

빌드하는 컴퓨터에만 [Go](https://go.dev) 가 필요합니다.

```bash
make          # 포맷·정적검사 후 세 가지 실행파일을 dist/ 에 빌드
go test ./... # 절전 차단이 실제로 걸리고, 해제 후 찌꺼기가 남지 않는지 검증
```

Go 는 크로스컴파일이 되므로 맥에서 윈도우용 exe 까지 한 번에 만들어집니다.
확인 주기는 `main.go` 의 `pollInterval` 에 있습니다.

| 파일 | 역할 |
| --- | --- |
| `main.go` | 감시 루프 + 설치 프로그램 분기 |
| `platform_windows.go` | Windows 감지·절전차단·작업스케줄러 등록·메시지창 |
| `platform_darwin.go` | macOS 감지·절전차단·LaunchAgent 등록·대화상자 |
| `platform_linux.go` | 리눅스(참고용) |
| `dist/` | 미리 빌드해 둔 실행파일 |

---

### In English

Keeps your computer awake while the **Claude desktop app** is running, and lets it sleep again as
soon as you quit Claude. Windows-first, macOS also supported.

Download the ZIP, unpack it, and double-click the executable in `dist/` — it installs itself, starts
at login, and needs no admin rights. Run the same file again to uninstall. A single ~2 MB static
binary; nothing else is required on the target machine. Messages are in Korean.

## License

MIT
