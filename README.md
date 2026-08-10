# claude-awake

Keep your computer awake while the **Claude desktop app** is running — and let it sleep
normally as soon as you quit Claude.

Cross-platform (**macOS & Windows**). Ships as a single compiled binary, so the machine
you install it on needs **no runtime** (no Node, no Bun) — just the one file.

## Why

Long Claude sessions and background agents get interrupted when the machine goes to sleep.
`claude-awake` blocks sleep only while the Claude app is open, then gets out of the way.

## How it works

A tiny watchdog polls every 15 seconds for the Claude app. While it's running, it holds an
**OS built-in** sleep inhibitor; when Claude quits, it releases it.

| Step                          | macOS                          | Windows                              |
| ----------------------------- | ------------------------------ | ------------------------------------ |
| Detect the app                | `pgrep` for the Claude process | `tasklist` for `Claude.exe`          |
| Block sleep while open        | `caffeinate -dimsu`            | PowerShell `SetThreadExecutionState` |
| Start at login + auto-restart | LaunchAgent                    | Scheduled Task (at logon)            |

No native addons — it only calls tools that already ship with each OS. The app is detected by
**process name**, not a fixed install path, so it works wherever Claude is installed.

## Install

Grab the binary for your OS from the [Releases](../../releases) page (or build it — see below),
then run the installer for your platform.

### macOS

```bash
# from the repo root, after placing the binary in dist/ (or run `bun run build`)
bash install/macos-install.sh
```

Registers a LaunchAgent that starts at login and restarts if it dies. Log: `/tmp/claude-awake.log`.

> First launch: macOS Gatekeeper may block an unsigned binary. Right-click the binary →
> **Open** once to allow it (or sign/notarize for wider distribution).

Uninstall: `bash install/macos-uninstall.sh`

### Windows

```powershell
powershell -ExecutionPolicy Bypass -File .\install\windows-install.ps1
```

Registers a Scheduled Task that runs at logon and restarts on failure.

> First launch: SmartScreen may warn about an unknown publisher → **More info → Run anyway**
> (or code-sign for wider distribution).

Uninstall: `powershell -ExecutionPolicy Bypass -File .\install\windows-uninstall.ps1`

## Build from source

Requires [Bun](https://bun.sh) **on the build machine only** (the target machine needs nothing).

```bash
bun install            # dev-only toolchain (types, eslint, prettier)
bun run build          # typecheck + lint + build all three: macOS arm64, macOS x64, Windows x64
# or individually:
bun run typecheck      # tsc --noEmit
bun run lint           # eslint --fix
bun run build:mac-arm64
bun run build:mac-x64
bun run build:win-x64
```

Output goes to `dist/`. Bun cross-compiles, so you can build the Windows `.exe` from a Mac.

The source is **TypeScript** (`index.ts`) — Bun runs and compiles `.ts` natively, so there is no
separate transpile step. The build is gated on `typecheck` **and** `lint`, so it fails rather than
shipping a binary that does not pass both.

### Code quality

Strict by default — `tsconfig.json` enables `strict`, `noImplicitReturns`, `noUnusedLocals`,
`noUnusedParameters`, `noFallthroughCasesInSwitch`, `exactOptionalPropertyTypes` and
`erasableSyntaxOnly`. ESLint runs `typescript-eslint` **strictTypeChecked** with `no-explicit-any`,
explicit return types, no floating promises, and Prettier enforced as a lint rule
(single quotes, trailing commas, 2-space indent, 150 columns).

## Run without installing (dev)

```bash
bun run start
```

## Configuration

Edit the constants at the top of `index.ts` (poll interval, etc.) and rebuild.

---

### 한국어 요약

Claude 데스크톱 앱이 **켜져 있는 동안만** 컴퓨터가 절전에 들지 않게 막아 주고, 앱을 끄면
절전이 자동 복귀합니다. **macOS·Windows 공용**이고, 설치하는 PC엔 **런타임(노드/번) 설치가
전혀 필요 없습니다** — 컴파일된 실행파일 하나만 받으면 됩니다.

- 감지: 설치 경로가 아니라 **프로세스 이름**으로 → 어디에 깔려 있든 동작
- 절전 차단: mac `caffeinate`, Windows `SetThreadExecutionState`(둘 다 OS 내장)
- 자동시작: mac LaunchAgent, Windows 작업 스케줄러(로그인 시)
- 빌드는 개발자 PC에서만 Bun 필요(`bun run build`) — 배포는 `dist/`의 바이너리를 Releases로.

## License

MIT
