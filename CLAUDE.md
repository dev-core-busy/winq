# CLAUDE.md — winq (Windows Command Agent)

## Projekt

**winq** – Port von [bashq](../linux_cmd_agent) für Windows.  
Übersetzt natürliche Sprache direkt im Terminal in präzise PowerShell-Befehlsketten und führt sie lokal aus.  
Kommuniziert über eine OpenAI-kompatible API mit einem lokalen LLM.

Config: `%AppData%\winq\config.json` · Log: `%AppData%\winq\activities.log`

## Build

```powershell
go build -o winq.exe .          # Entwicklungsbuild
```

Cross-Compile von Linux:
```bash
GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -o winq.exe .
```

## Status: Ausgangsbasis kopiert — noch nicht Windows-fähig

Der Code kompiliert für Windows, ist aber inhaltlich noch Linux/bashq.  
Folgende Aufgaben müssen **in dieser Reihenfolge** erledigt werden:

---

## TODO-Liste (vollständig, sequenziell abarbeiten)

### 1. `main.go` — `syscall.Exec` ersetzen

`syscall.Exec` existiert unter Windows nicht. Auto-Restart nach Update muss anders gelöst werden.

**Aktuell:**
```go
import "syscall"
syscall.Exec(restartExecPath, os.Args, os.Environ())
```

**Ersetzen durch:**
```go
import "os/exec"
cmd := exec.Command(restartExecPath, os.Args[1:]...)
cmd.Stdin = os.Stdin
cmd.Stdout = os.Stdout
cmd.Stderr = os.Stderr
cmd.Start()
os.Exit(0)
```

---

### 2. `update.go` — `cmdRunCommand`: `bash -c` → PowerShell

```go
// Aktuell (Linux):
cmd := exec.Command("bash", "-c", command)

// Windows:
cmd := exec.Command("powershell", "-NoProfile", "-NonInteractive", "-Command", command)
```

---

### 3. `updater.go` — Asset-Name und Restart

- `currentVersion` auf `"v1.0.0"` belassen (eigene Release-Kette)
- `githubReleasesAPI` auf neues GitHub-Repo zeigen lassen (nach Repo-Erstellung)
- Asset-Name: `"winq-windows-amd64.exe"` statt `"bashq-linux-" + runtime.GOARCH`
- `syscall.Exec` im `updateDoneMsg`-Handler durch `exec.Command` + `os.Exit(0)` ersetzen (wie in Punkt 1)

---

### 4. `selfinstall.go` — Windows-Installation ohne Symlinks

Symlinks in `PATH` sind unter Windows unüblich. Stattdessen:

**Strategie:** Binary in `%LOCALAPPDATA%\Programs\winq\winq.exe` ablegen und  
diesen Pfad per `reg add` in `HKCU\Environment\PATH` eintragen (kein Admin nötig).

**Deinstallation:** Ordner löschen, PATH-Eintrag entfernen.

```go
func selfInstallToggle() (msg string, isErr bool) {
    exe, _ := os.Executable()
    installDir := filepath.Join(os.Getenv("LOCALAPPDATA"), "Programs", "winq")
    target := filepath.Join(installDir, "winq.exe")

    if _, err := os.Stat(target); err == nil {
        // Deinstallieren
        os.RemoveAll(installDir)
        removeFromUserPath(installDir)
        return "✓ winq deinstalliert", false
    }

    os.MkdirAll(installDir, 0755)
    // Binary kopieren (kein Symlink auf Windows)
    copyFile(exe, target)
    addToUserPath(installDir)  // reg add HKCU\Environment /v Path ...
    return fmt.Sprintf("✓ %s\n  Neues Terminal öffnen damit PATH aktiv wird", target), false
}
```

Hilfsfunktionen `addToUserPath` / `removeFromUserPath` via `golang.org/x/sys/windows/registry`.

---

### 5. `agent.go` — System-Prompt auf PowerShell/Windows umstellen

Den `systemPrompt` (in `i18n.go` unter `SystemPrompt`) komplett neu schreiben:

- Fokus: **PowerShell 5.1 / 7+**, `cmd.exe` nur als Fallback
- Windows-Pfadkonventionen (`C:\`, `%USERPROFILE%`, `$env:`)
- Paketmanagement: `winget`, `choco` (Chocolatey), `scoop`
- Dienste: `Get-Service`, `Start-Service`, `Stop-Service`
- Prozesse: `Get-Process`, `Stop-Process`
- Netzwerk: `Test-NetConnection`, `Get-NetAdapter`, `ipconfig`
- Dateisystem: `Get-ChildItem`, `Copy-Item`, `Remove-Item`
- Kein `bash`, kein `apt`, kein `systemctl`

---

### 6. `commands.go` + `i18n.go` — Slash-Commands auf Windows anpassen

Alle `Message`-Felder der Slash-Commands von Linux auf Windows-Äquivalente umstellen:

| Command     | Linux (aktuell)              | Windows (soll)                              |
|-------------|------------------------------|---------------------------------------------|
| `/update`   | `apt update && apt upgrade`  | `winget upgrade --all`                      |
| `/status`   | `systemctl status`           | `Get-Service \| Where Status -eq Running`   |
| `/disk`     | `df -h`                      | `Get-PSDrive -PSProvider FileSystem`        |
| `/memory`   | `free -h`                    | `Get-CimInstance Win32_OperatingSystem`     |
| `/network`  | `ip addr show`               | `Get-NetAdapter; ipconfig`                  |
| `/services` | `systemctl list-units`       | `Get-Service`                               |
| `/logs`     | `journalctl -n 50`           | `Get-EventLog -LogName System -Newest 50`   |
| `/optimize` | `apt autoremove`             | `winget upgrade --all; Optimize-Volume`     |

Außerdem `/colors`-Befehl aus `commands.go` und `i18n.go` **entfernen** — nicht benötigt unter Windows.

---

### 7. `config_persist.go` — Pfad von `bashq` auf `winq` umbenennen

```go
// Aktuell:
filepath.Join(dir, "bashq", "config.json")

// Soll:
filepath.Join(dir, "winq", "config.json")
```

`os.UserConfigDir()` gibt unter Windows `%AppData%\Roaming` zurück — das ist korrekt.

---

### 8. `activities.go` — Log-Pfad anpassen

Analog zu config_persist.go: `"linux_cmd_agent"` → `"winq"` in allen Pfaden.

---

### 9. `i18n.go` — Strings bereinigen

- Alle Linux-spezifischen Hinweise in Fehlermeldungen ersetzen (z.B. `sudo`, `~/.local/bin`)
- `SystemPrompt` komplett neu schreiben (Windows/PowerShell — siehe Punkt 5)
- `CmdColorsDesc` entfernen (Spalte in UIStrings-Struct + alle drei Sprachblöcke)

---

### 10. `colors.go` — bereits als Stub vorhanden, kein Handlungsbedarf

---

### 11. Abschluss: Bauen + Testen

```powershell
# Auf Windows:
go build -o winq.exe .
.\winq.exe

# Cross-Compile von Linux (für Release):
GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -o winq.exe .
```

Windows Terminal (WT) oder PowerShell 7 empfohlen — klassisches `cmd.exe` unterstützt keine vollen ANSI-Farben.

---

## Architektur (identisch zu bashq)

| Datei               | Verantwortung                                          |
|---------------------|--------------------------------------------------------|
| `model.go`          | Bubbletea-Modell, Typen, Hilfsmethoden                 |
| `update.go`         | Bubbletea `Update()`, Zustandsmaschine                 |
| `view.go`           | Bubbletea `View()`, alle Render-Funktionen             |
| `styles.go`         | Lipgloss-Styles                                        |
| `agent.go`          | HTTP-Client, Tool-Calling, Gesprächsverlauf            |
| `commands.go`       | Slash-Befehlsdefinitionen                              |
| `activities.go`     | Aktivitätsprotokoll                                    |
| `config_persist.go` | JSON-Persistenz                                        |
| `selfinstall.go`    | Installation via PATH-Eintrag (Windows)                |
| `updater.go`        | GitHub-Releases-Checker, Binary-Download               |
| `colors.go`         | No-op Stub (Windows braucht kein .bashrc-Setup)        |
| `i18n.go`           | Mehrsprachigkeit (de/en/zh)                            |
| `discovery.go`      | LLM-Profil-Erkennung                                   |
| `session_persist.go`| Sitzungsspeicherung                                    |

## Abhängigkeiten hinzufügen

Für `selfinstall.go` (Registry-Zugriff):
```bash
go get golang.org/x/sys
```

## Referenz

Originalprojekt: `../linux_cmd_agent` (bashq)  
Alle TODO-Punkte sind unabhängig voneinander außer Punkt 1+3 (beide betreffen syscall.Exec).
