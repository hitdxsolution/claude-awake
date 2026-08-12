# claude-awake

**Claude 앱을 켜 두면 컴퓨터가 잠들지 않습니다.** Claude 를 끄면 절전이 원래대로 돌아옵니다.

Claude 로 오래 작업하거나 뭔가를 돌려 두는 중에 컴퓨터가 절전에 들어가 작업이 끊기는 것을 막아 줍니다.
절전 설정을 아예 꺼 둘 필요가 없습니다.

---

## 다운로드

아래 링크를 누르면 바로 받아집니다. **받은 파일을 더블클릭하면 설치 끝입니다.**

### ⬇ [Windows용 받기](https://github.com/hitdxsolution/claude-awake/releases/latest/download/claude-awake-windows-x64.exe)

### ⬇ macOS용 받기 — [애플 실리콘(M1 이후)](https://github.com/hitdxsolution/claude-awake/releases/latest/download/claude-awake-macos-arm64) · [인텔 맥](https://github.com/hitdxsolution/claude-awake/releases/latest/download/claude-awake-macos-x64)

> 내 맥이 어느 쪽인지 모르겠으면 화면 왼쪽 위 사과 메뉴 → **이 Mac에 관하여** 에서 확인할 수 있습니다.
> 칩 이름에 **Apple** 이 있으면 애플 실리콘, **Intel** 이 있으면 인텔 맥입니다.

---

## 설치 (Windows)

### 1. 파일 받기

위 **[Windows용 받기](https://github.com/hitdxsolution/claude-awake/releases/latest/download/claude-awake-windows-x64.exe)** 를 누릅니다.

### 2. 받은 파일 더블클릭

보통 `다운로드` 폴더에 `claude-awake-windows-x64.exe` 가 받아집니다. **더블클릭**하세요.

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
| **바이러스가 발견됨** / 파일이 사라짐 | 오탐입니다. 아래 [백신이 막을 때](#백신이-막을-때) 를 보세요. |
| 아무 창도 안 뜸 | 압축을 **풀지 않고** ZIP 안에서 바로 실행하면 그럴 수 있습니다. 압축을 먼저 푸세요. |

---

## 백신이 막을 때

다운로드 중에 **"바이러스가 발견됨"** 이라며 파일이 사라지거나, 받은 파일이 조용히 격리될 수 있습니다.
**오탐입니다.** 실제로 아무 데도 접속하지 않고, 개인정보를 읽지도 보내지도 않습니다.
소스가 전부 이 저장소에 공개돼 있으니 직접 확인하실 수 있습니다.

### 왜 그런가

**만든 사람 서명(코드 서명 인증서)이 없기 때문입니다.** 인증서는 연 수십만 원짜리 유료라 붙이지 않았습니다.
서명이 없으면 브라우저의 SmartScreen 이 "받은 사람이 거의 없는, 서명 없는 프로그램" 으로 보고
내려받는 단계에서 먼저 끊습니다. 여기에 백신의 기계학습 판정이 겹치면
`Trojan:Win32/Wacatac` 같은 **이름만 그럴듯한 오탐**이 붙습니다.

0.2 버전에서 오탐을 줄이려고 아래를 손봤습니다. 그래도 백신에 따라 여전히 걸릴 수 있습니다.

- 설치할 때 PowerShell 을 부르지 않습니다 — 윈도우 기본 도구인 `schtasks.exe` 만 씁니다
- 실행파일에 만든 곳·이름·설명·버전·아이콘을 넣었습니다(비어 있으면 그 자체로 의심 점수가 올라갑니다)
- 관리자 권한을 요구하지 않는다는 것을 매니페스트에 명시했습니다
- 설치 위치를 `%LOCALAPPDATA%\Programs\` 아래로 옮겼습니다(VS Code 같은 정상 앱이 쓰는 자리)
- GitHub Actions 에서 빌드해, 어느 커밋에서 만들어졌는지 증명이 함께 올라갑니다

### 다운로드부터 막힐 때

브라우저가 파일을 지워 버리면 이렇게 하세요.

1. **Edge/Chrome**: 오른쪽 위 **다운로드** 목록을 열고, 지워진 항목의 **⋯** → **계속 유지** 를 누릅니다.
2. 그래도 안 되면 **다른 브라우저**로 받아 보세요.
3. 회사 PC 라 백신 설정을 못 바꾸는 경우에는 전산 담당자에게 이 문서를 보여 주세요.

### 받았는데 실행이 막힐 때

Windows 보안(Defender) 기준입니다.

1. **시작** → `Windows 보안` 을 검색해 엽니다.
2. **바이러스 및 위협 방지** → **보호 기록** 에서 격리된 `claude-awake` 를 찾아 **디바이스에서 허용** 을 누릅니다.
3. 파일이 이미 지워졌다면 다시 받은 뒤, **바이러스 및 위협 방지** → **설정 관리** →
   **제외 항목 추가 또는 제거** → **제외 항목 추가** → **파일** 로 받은 파일을 등록합니다.

V3·알약 등 다른 백신도 **격리 보관소 / 검사 제외** 메뉴에 같은 기능이 있습니다.

### 정 못 미더우면

[VirusTotal](https://www.virustotal.com/) 에 파일을 올리면 70여 개 백신이 어떻게 보는지 한 번에 확인할 수 있습니다.
소수의 엔진만 반응한다면 오탐으로 보시면 됩니다.
그래도 찜찜하면 설치하지 마시고, [소스](platform_windows.go)를 읽어 보신 뒤 직접 빌드해 쓰셔도 됩니다.

---

## 설치 (macOS)

### 1. 파일 받기

위 **macOS용 받기** 에서 내 맥에 맞는 것을 누릅니다
([애플 실리콘](https://github.com/hitdxsolution/claude-awake/releases/latest/download/claude-awake-macos-arm64) ·
[인텔 맥](https://github.com/hitdxsolution/claude-awake/releases/latest/download/claude-awake-macos-x64)).

### 2. 받은 파일을 우클릭 → 열기

> 처음에는 반드시 **우클릭 → 열기**로 실행하세요. 그냥 더블클릭하면
> "확인되지 않은 개발자" 경고만 뜨고 실행되지 않습니다. 한 번 이렇게 열어 두면 다음부터는 그냥 열립니다.

### 3. 끝

"설치했습니다" 창이 뜨면 완료입니다. 끄고 싶으면 **같은 파일을 다시 실행**하면 제거할지 물어봅니다.

---

## 잘 되고 있는지 보고 싶다면

설치만 하면 신경 쓸 것이 없지만, 확인하고 싶다면 기록 파일을 열어 보면 됩니다.

- **Windows**: `C:\Users\내계정\AppData\Local\Programs\claude-awake\claude-awake.log`
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
아이콘을 고칠 때만 [ImageMagick](https://imagemagick.org) 이 추가로 필요합니다.

```bash
make          # 포맷·정적검사 후 세 가지 실행파일을 dist/ 에 빌드
go test ./... # 절전 차단이 실제로 걸리고, 해제 후 찌꺼기가 남지 않는지 검증
```

Go 는 크로스컴파일이 되므로 맥에서 윈도우용 exe 까지 한 번에 만들어집니다.
윈도우 리소스(버전 정보·아이콘·매니페스트)는 `make` 가 `goversioninfo` 를 받아 `.syso` 로 만들어 링크합니다.
확인 주기는 `main.go` 의 `pollInterval` 에 있습니다.

| 파일 | 역할 |
| --- | --- |
| `main.go` | 감시 루프 + 설치 프로그램 분기 |
| `platform_windows.go` | Windows 감지·절전차단·작업스케줄러 등록·메시지창 |
| `platform_darwin.go` | macOS 감지·절전차단·LaunchAgent 등록·대화상자 |
| `platform_linux.go` | 리눅스(참고용) |
| `versioninfo.json` | exe 에 박히는 버전 리소스 정의 |
| `assets/` | 아이콘 원본(SVG·ICO)과 매니페스트 |
| `.github/workflows/release.yml` | 태그를 밀면 빌드·증명·릴리스 |
| `dist/` | 미리 빌드해 둔 실행파일 |

### 오탐이 다시 붙으면

서명이 없는 한 완전히 없앨 수는 없고, 릴리스마다 다시 걸릴 수 있습니다. 그때는 신고합니다.

1. [VirusTotal](https://www.virustotal.com/) 에 새 exe 를 올려 **어떤 엔진이 무슨 이름으로** 잡는지 확인합니다.
2. [Microsoft 오탐 신고](https://www.microsoft.com/en-us/wdsi/filesubmission) 에 제출합니다 — 보통 1~3 일 안에 정의가 갱신됩니다.
   구분은 **Software developer**, 유형은 **Incorrectly detected as malware** 를 고릅니다.
3. 국내 백신(V3·알약)은 각 벤더의 오진 신고 창구에 따로 넣어야 합니다.

근본 해결은 코드 서명 인증서입니다(OV 연 20~40 만원대, SmartScreen 경고까지 없애려면 EV).
붙이게 되면 `signtool` 을 `windows-x64` 타깃 뒤에 한 줄 추가하면 됩니다.

---

### In English

Keeps your computer awake while the **Claude desktop app** is running, and lets it sleep again as
soon as you quit Claude. Windows-first, macOS also supported.

Download the ZIP, unpack it, and double-click the executable in `dist/` — it installs itself, starts
at login, and needs no admin rights. Run the same file again to uninstall. A single ~2 MB static
binary; nothing else is required on the target machine. Messages are in Korean.

## License

MIT
