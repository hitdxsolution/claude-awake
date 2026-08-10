#!/usr/bin/env bun
// claude-awake — keep the computer awake while the Claude desktop app is running.
//
//   Cross-platform (macOS · Windows · best-effort Linux). Shipped as a single compiled
//   binary, so the target machine needs NO runtime (Node/Bun) installed.
//
//   How it works: poll every few seconds for the Claude desktop app. While it is running,
//   hold an OS sleep inhibitor; release it as soon as the app quits.
//     macOS   → `caffeinate -dimsu`         (built in)
//     Windows → PowerShell SetThreadExecutionState (built in)
//     Linux   → `systemd-inhibit`           (best effort)
//   No native addons — we only shell out to tools that already ship with each OS.

import { spawn, spawnSync, type ChildProcess } from 'node:child_process';

/** How often to check whether Claude is running. */
const POLL_MS = 15_000;
const APP_LABEL = 'Claude';

/** The three platforms we adapt to; anything else falls back to the Linux path. */
type Platform = NodeJS.Platform;
const platform: Platform = process.platform;

/** stdout line (LaunchAgent / Task Scheduler capture this to a log file). */
function log(message: string): void {
  process.stdout.write(`[claude-awake] ${new Date().toISOString()} ${message}\n`);
}

/** Is the Claude desktop app currently running? Detected by process, not install path. */
function isClaudeRunning(): boolean {
  try {
    if (platform === 'darwin') {
      // main app executable is .../Claude.app/Contents/MacOS/Claude
      const result = spawnSync('pgrep', ['-f', 'Claude.app/Contents/MacOS/Claude'], { encoding: 'utf8' });
      return result.status === 0 && result.stdout.trim() !== '';
    }
    if (platform === 'win32') {
      const result = spawnSync('tasklist', ['/FI', 'IMAGENAME eq Claude.exe', '/NH'], { encoding: 'utf8' });
      return typeof result.stdout === 'string' && result.stdout.toLowerCase().includes('claude.exe');
    }
    const result = spawnSync('pgrep', ['-x', 'claude'], { encoding: 'utf8' }); // linux best-effort
    return result.status === 0 && result.stdout.trim() !== '';
  } catch {
    return false;
  }
}

/** Start a long-lived child that blocks system sleep. Killing it releases the block. */
function startInhibitor(): ChildProcess {
  if (platform === 'darwin') {
    // -dimsu: block display/idle/system/disk sleep + declare the user active. Lives until killed.
    return spawn('caffeinate', ['-dimsu'], { stdio: 'ignore' });
  }
  if (platform === 'win32') {
    // Hold ES_CONTINUOUS | ES_SYSTEM_REQUIRED | ES_DISPLAY_REQUIRED (0x80000003) on this thread.
    //   Windows clears the requirement automatically when this PowerShell process exits (we kill it).
    const script = [
      "$s='[DllImport(\"kernel32.dll\")] public static extern uint SetThreadExecutionState(uint e);';",
      '$k=Add-Type -MemberDefinition $s -Name P -Namespace W -PassThru;',
      '[void]$k::SetThreadExecutionState(0x80000003);',
      'while($true){Start-Sleep -Seconds 3600}',
    ].join(' ');
    return spawn('powershell', ['-NoProfile', '-NonInteractive', '-WindowStyle', 'Hidden', '-Command', script], {
      stdio: 'ignore',
      windowsHide: true,
    });
  }
  // linux best-effort: the block lasts while the child (sleep infinity) lives.
  return spawn('systemd-inhibit', ['--what=idle:sleep', '--who=claude-awake', '--why=Claude running', 'sleep', 'infinity'], {
    stdio: 'ignore',
  });
}

let inhibitor: ChildProcess | null = null;

function stopInhibitor(): void {
  if (inhibitor !== null) {
    try {
      inhibitor.kill();
    } catch {
      // already gone
    }
    inhibitor = null;
  }
}

/** One check: match the inhibitor to whether Claude is running. */
function tick(): void {
  const running = isClaudeRunning();
  if (running && inhibitor === null) {
    const child = startInhibitor();
    inhibitor = child;
    // If the inhibitor dies on its own, forget it so we respawn on the next tick.
    child.on('exit', () => {
      if (inhibitor === child) {
        inhibitor = null;
      }
    });
    log(`${APP_LABEL} detected → blocking sleep`);
  } else if (!running && inhibitor !== null) {
    stopInhibitor();
    log(`${APP_LABEL} closed → sleep allowed`);
  }
}

// Never leave the machine caffeinated if we are stopped.
const exitSignals: NodeJS.Signals[] = ['SIGINT', 'SIGTERM', 'SIGHUP'];
for (const signal of exitSignals) {
  process.on(signal, () => {
    stopInhibitor();
    process.exit(0);
  });
}
process.on('exit', stopInhibitor);

log(`started — platform=${platform}, poll=${POLL_MS / 1000}s`);
tick();
setInterval(tick, POLL_MS);
