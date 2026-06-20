<div align="center">

# winq

### Your terminal. Your LLM. Your command.

Type what you want. winq turns plain language into PowerShell commands,  
explains them, confirms, and runs them — all locally, all instantly.

[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Go 1.22+](https://img.shields.io/badge/Go-1.22+-00ADD8?logo=go&logoColor=white)](go.mod)
[![Platform](https://img.shields.io/badge/Platform-Windows-0078D4?logo=windows&logoColor=white)](.)
[![Releases](https://img.shields.io/github/v/release/dev-core-busy/winq?color=4fc3f7)](../../releases)

**Sister project:** [bashq](https://github.com/dev-core-busy/bashq) — the Linux version

</div>

---

## What is winq?

**winq** is a minimalist TUI agent that lives entirely in your terminal.  
Describe what you need in plain English (or German, or Chinese) — winq figures out the exact PowerShell command chain, shows it to you with a clear explanation, and executes it on your confirmation.

> **Named after Q from Star Trek** — the omnipotent being of the Q Continuum who can answer any question and reshape reality with a thought. winq brings that same effortlessness to your Windows command line.

No cloud. No account. No data leaving your machine. Just you, your terminal, and a local LLM.

---

## Why winq?

| | |
|---|---|
| 🧠 **Plain language in, PowerShell out** | Type *"find the 10 largest files on this system"* — not `Get-ChildItem C:\ -Recurse -ErrorAction SilentlyContinue \| Sort-Object Length -Descending \| Select-Object -First 10` |
| 👁 **You stay in control** | Every command is shown with a plain-language explanation before anything runs. One keypress to confirm, one to cancel. |
| 🔒 **100% local & private** | Runs on any OpenAI-compatible local LLM. Nothing ever leaves your machine. |
| 📡 **Auto-discovery** | Enter an IP — winq scans common ports, detects the model, and sets up the profile in seconds. |
| ⚡ **Single static binary** | Download `winq-windows-amd64.exe` and run — no installer, no dependencies, no PowerShell module required. |
| 💾 **Persistent sessions** | Restart and pick up exactly where you left off. Chat history and LLM context are preserved. |
| 🌍 **Multi-language UI** | German, English, Chinese — auto-detected from your system locale, switchable at runtime. |
| 🎛 **F1–F9 one-key macros** | Pre-configure your most-used queries as keyboard shortcuts. |

---

## Install

### Quick start

1. Download `winq-windows-amd64.exe` from the **[Releases](../../releases)** page
2. Open **Windows Terminal** or **PowerShell 7** and run it:

```powershell
.\winq-windows-amd64.exe
```

3. Type `/setup` inside the app to install system-wide — winq copies itself to  
   `%LOCALAPPDATA%\Programs\winq\winq.exe` and registers it in your user PATH.  
   Open a new terminal and run `winq` from anywhere.

> **Recommended:** Windows Terminal or PowerShell 7 for full ANSI color support.  
> Classic `cmd.exe` works but colors may be limited.

### Connect your LLM

On first start, winq connects to `http://localhost:11434/v1` (Ollama default).  
Open `/config` to change the endpoint — or let winq find your LLM automatically:

1. Type `/config` and press Enter
2. Navigate to **LLM PROFILES → [ + New LLM Profile ]**
3. Enter your LLM server's IP address
4. winq scans ports `11434 1234 8080 8000 9081 7860 5000 3000` automatically
5. Select a model from the list — done ✓

### Ask anything

```
> Show me which processes are using the most memory
> List all Windows services that have stopped unexpectedly
> Find files larger than 1 GB modified in the last 7 days
> What's eating my disk space in C:\Users?
> Schedule a daily task to clean up the Temp folder
```

---

## Compatible LLM Servers

Any **OpenAI-compatible** local server works out of the box:

| Server | Default port | Notes |
|--------|-------------|-------|
| [Ollama](https://ollama.com) | 11434 | Recommended — easiest setup |
| [LM Studio](https://lmstudio.ai) | 1234 | Great GUI for model management |
| [vLLM](https://github.com/vllm-project/vllm) | 8000 | High-throughput production server |
| [llama.cpp server](https://github.com/ggerganov/llama.cpp) | 8080 | Lightweight, runs anywhere |
| [Jan](https://jan.ai) | 1234 | Cross-platform desktop app |
| Any OpenAI-compatible API | any | Including cloud providers |

> **Recommended models:** Qwen3, Llama 3.1/3.2, Mistral, DeepSeek-Coder, Gemma 2

---

## Settings & Configuration

Type `/config` to open the settings editor.

### LLM Profiles

Save multiple LLM endpoints and switch between them instantly.  
Mark one as **preferred** (`P`) — winq health-checks it on startup and suggests the next available profile if it's unreachable.

### Execution Modes

| Mode | Behaviour |
|------|-----------|
| **ASK** (default) | Shows command + explanation, waits for confirmation |
| **AUTO** | Executes commands immediately — for repetitive workflows |

Toggle with `Shift+Tab` from anywhere. The title bar always shows the active mode.

### Session Persistence

winq saves your entire conversation and LLM context on exit — including tool call history.  
Toggle with `Alt+S` or in `/config → SESSIONS`. The title bar shows a 💾 SESSION badge when active.

---

## Keyboard Reference

| Key | Action |
|-----|--------|
| `Enter` | Send message / confirm command |
| `J` / `Enter` | Confirm pending command |
| `N` / `Esc` | Cancel pending command |
| `↑ / ↓` | Scroll history / navigate lists |
| `Shift+Tab` | Toggle ASK ↔ AUTO execution mode |
| `Alt+S` | Toggle session persistence on/off |
| `F1–F9` / `Alt+1–9` | Custom one-key shortcuts (configure in `/config`) |
| `Tab` | Switch between profile list and settings in `/config` |
| `/` | Open command autocomplete |
| `Ctrl+C` | Save session and quit |

## Slash Commands

Type `/` for autocomplete with descriptions.

| Command | What it does |
|---------|-------------|
| `/install` | Install software packages (winget) |
| `/update` | Update all installed programs (winget) |
| `/status` | Full system overview |
| `/disk` | Disk space analysis |
| `/memory` | Memory usage breakdown |
| `/network` | Network interfaces & connectivity |
| `/services` | Manage Windows services |
| `/logs` | Recent system event log entries |
| `/config` | Open settings editor |
| `/setup` | Install winq system-wide (or remove — acts as a toggle) |
| `/activities` | Show command history with timestamps |
| `/clear` | Clear chat history and start fresh |
| `/help` | Show keyboard shortcuts and tips |

---

## Building from Source

Requires **Go 1.22+**. Cross-compile from Linux or build directly on Windows.

```bash
# Cross-compile from Linux (produces winq-windows-amd64.exe):
bash build.sh

# Build on Windows:
go build -o winq.exe .
```

```bash
go vet ./...           # static analysis
go build ./...         # compile-check without producing a binary
```

---

## Files

| Path | Purpose |
|------|---------|
| `%AppData%\winq\config.json` | Settings (LLM profiles, shortcuts, preferences) |
| `%AppData%\winq\activities.log` | Full history of every query, response and command |
| `%AppData%\winq\session.json` | Saved session state (chat + LLM context) |

All files are human-readable JSON / plain text. Delete any of them to reset that component.

---

## The Name

**win** + **q** — the Q Continuum for Windows.

In Star Trek, Q is an omnipotent entity who knows everything, can do anything, and never needs to look anything up. That's the energy winq brings to your Windows command line: ask it anything about your system, and it handles the rest.

---

## Sister Project

**[bashq](https://github.com/dev-core-busy/bashq)** — the Linux sibling.  
Same concept, same architecture, same local-LLM philosophy — built for bash and Debian/Ubuntu instead of PowerShell and Windows.

---

## License

MIT — see [LICENSE](LICENSE).

---

<div align="center">

**[⭐ Star this repo](../../stargazers)** if winq saves you time — it helps others find it.

*Runs entirely on your machine. Your queries never leave your terminal.*

</div>
