package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

const currentVersion = "v1.0.13"
const githubReleasesAPI = "https://api.github.com/repos/dev-core-busy/winq/releases/latest"

type updateInfo struct {
	version     string
	downloadURL string
}

type updateCheckMsg struct {
	info *updateInfo // nil = kein Update verfügbar oder Fehler
	err  error
}

type updateDoneMsg struct {
	version string
	err     error
}

// scheduleUpdateCheckMsg feuert nach 30 Minuten Wartezeit.
type scheduleUpdateCheckMsg struct{}

type githubRelease struct {
	TagName string `json:"tag_name"`
	Assets  []struct {
		Name               string `json:"name"`
		BrowserDownloadURL string `json:"browser_download_url"`
	} `json:"assets"`
}

// cmdCheckUpdate fragt die GitHub-API nach der neuesten Version.
func cmdCheckUpdate() tea.Cmd {
	return func() tea.Msg {
		client := &http.Client{Timeout: 10 * time.Second}
		req, err := http.NewRequest("GET", githubReleasesAPI, nil)
		if err != nil {
			return updateCheckMsg{err: err}
		}
		req.Header.Set("Accept", "application/vnd.github+json")
		req.Header.Set("User-Agent", "winq/"+currentVersion)

		resp, err := client.Do(req)
		if err != nil {
			return updateCheckMsg{err: err}
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			return updateCheckMsg{} // kein Update, kein Fehler
		}

		var release githubRelease
		if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
			return updateCheckMsg{err: err}
		}

		if release.TagName == "" || !semverGT(release.TagName, currentVersion) {
			return updateCheckMsg{} // bereits aktuell oder älter
		}

		assetName := "winq-windows-amd64.exe"
		downloadURL := ""
		for _, asset := range release.Assets {
			if asset.Name == assetName {
				downloadURL = asset.BrowserDownloadURL
				break
			}
		}
		if downloadURL == "" {
			return updateCheckMsg{}
		}

		return updateCheckMsg{info: &updateInfo{
			version:     release.TagName,
			downloadURL: downloadURL,
		}}
	}
}

// cmdScheduleUpdateCheck wartet 30 Minuten und löst dann einen neuen Check aus.
func cmdScheduleUpdateCheck() tea.Cmd {
	return func() tea.Msg {
		time.Sleep(30 * time.Minute)
		return scheduleUpdateCheckMsg{}
	}
}

// cmdDownloadUpdate lädt die neue Binary als .pending-Datei herunter.
// Die eigentliche Ersetzung übernimmt applyPendingUpdate() beim nächsten Start —
// dann ist kein Prozess mehr aktiv, der die Datei sperrt (Windows-spezifisches Problem).
func cmdDownloadUpdate(info updateInfo) tea.Cmd {
	return func() tea.Msg {
		exe, err := os.Executable()
		if err != nil {
			return updateDoneMsg{err: fmt.Errorf("Pfad nicht ermittelbar: %w", err)}
		}
		exe, _ = filepath.EvalSymlinks(exe)

		client := &http.Client{Timeout: 120 * time.Second}
		resp, err := client.Get(info.downloadURL)
		if err != nil {
			return updateDoneMsg{err: fmt.Errorf("Download fehlgeschlagen: %w", err)}
		}
		defer resp.Body.Close()

		pendingFile := exe + ".pending"
		out, err := os.OpenFile(pendingFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
		if err != nil {
			return updateDoneMsg{err: fmt.Errorf("Temp-Datei: %w", err)}
		}
		if _, err = io.Copy(out, resp.Body); err != nil {
			out.Close()
			os.Remove(pendingFile)
			return updateDoneMsg{err: fmt.Errorf("Schreiben: %w", err)}
		}
		out.Close()

		return updateDoneMsg{version: info.version}
	}
}

// applyPendingUpdate prüft beim Programmstart ob ein heruntergeladenes Update (.pending)
// vorliegt, wendet es an und startet winq mit dem neuen Binary neu.
// Muss ganz am Anfang von main() aufgerufen werden — vor jeder TUI-Initialisierung,
// damit die Konsolen-Handles sauber vererbt werden.
func applyPendingUpdate() {
	exe, err := os.Executable()
	if err != nil {
		return
	}
	exe, _ = filepath.EvalSymlinks(exe)
	pendingFile := exe + ".pending"
	if _, err := os.Stat(pendingFile); os.IsNotExist(err) {
		return
	}

	oldFile := exe + ".old"
	os.Remove(oldFile)
	if err := os.Rename(exe, oldFile); err != nil {
		return
	}
	if err := os.Rename(pendingFile, exe); err != nil {
		os.Rename(oldFile, exe) // Rollback
		return
	}
	os.Remove(oldFile)

	// Neustart mit neuem Binary (vor TUI: Handles noch unverändert → klappt auf Windows)
	cmd := exec.Command(exe, os.Args[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		return
	}
	os.Exit(0)
}

// semverGT gibt true zurück wenn a semantisch größer als b ist (z.B. "v1.0.10" > "v1.0.9").
func semverGT(a, b string) bool {
	parse := func(v string) [3]int {
		v = strings.TrimPrefix(strings.TrimSpace(v), "v")
		parts := strings.SplitN(v, ".", 3)
		for len(parts) < 3 {
			parts = append(parts, "0")
		}
		var n [3]int
		for i := 0; i < 3; i++ {
			// Nur die führenden Ziffern auswerten: bei "12-rc1" lieferte
			// Atoi sonst einen Fehler und damit 0 – ein Vorabversions-Tag
			// hätte dadurch als älter gegolten als sein Vorgänger.
			p := parts[i]
			end := 0
			for end < len(p) && p[end] >= '0' && p[end] <= '9' {
				end++
			}
			n[i], _ = strconv.Atoi(p[:end])
		}
		return n
	}
	av, bv := parse(a), parse(b)
	for i := 0; i < 3; i++ {
		if av[i] != bv[i] {
			return av[i] > bv[i]
		}
	}
	return false
}
