package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/sys/windows/registry"
)

func selfInstallToggle() (msg string, isErr bool) {
	exe, err := os.Executable()
	if err != nil {
		return "Fehler: Pfad der aktuellen Binary nicht ermittelbar: " + err.Error(), true
	}
	exe, _ = filepath.EvalSymlinks(exe)

	installDir := filepath.Join(os.Getenv("LOCALAPPDATA"), "Programs", "winq")
	target := filepath.Join(installDir, "winq.exe")

	if _, err := os.Stat(target); err == nil {
		if err := os.RemoveAll(installDir); err != nil {
			return fmt.Sprintf("✗ Konnte %s nicht entfernen: %v", installDir, err), true
		}
		if err := removeFromUserPath(installDir); err != nil {
			return fmt.Sprintf("✓ winq deinstalliert (PATH-Eintrag konnte nicht entfernt werden: %v)", err), false
		}
		return "✓ winq deinstalliert", false
	}

	if err := os.MkdirAll(installDir, 0755); err != nil {
		return fmt.Sprintf("✗ Konnte Verzeichnis nicht anlegen: %v", err), true
	}
	if err := copyFile(exe, target); err != nil {
		return fmt.Sprintf("✗ Kopieren fehlgeschlagen: %v", err), true
	}
	if err := addToUserPath(installDir); err != nil {
		return fmt.Sprintf("✓ %s\n  ⚠ PATH-Eintrag fehlgeschlagen: %v", target, err), false
	}
	return fmt.Sprintf("✓ %s\n  Neues Terminal öffnen damit PATH aktiv wird", target), false
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, in)
	return err
}

func addToUserPath(dir string) error {
	k, err := registry.OpenKey(registry.CURRENT_USER, `Environment`, registry.QUERY_VALUE|registry.SET_VALUE)
	if err != nil {
		return err
	}
	defer k.Close()
	current, _, err := k.GetStringValue("Path")
	if err != nil {
		current = ""
	}
	if strings.Contains(strings.ToLower(current), strings.ToLower(dir)) {
		return nil
	}
	newPath := current
	if newPath != "" && !strings.HasSuffix(newPath, ";") {
		newPath += ";"
	}
	newPath += dir
	return k.SetStringValue("Path", newPath)
}

func removeFromUserPath(dir string) error {
	k, err := registry.OpenKey(registry.CURRENT_USER, `Environment`, registry.QUERY_VALUE|registry.SET_VALUE)
	if err != nil {
		return err
	}
	defer k.Close()
	current, _, err := k.GetStringValue("Path")
	if err != nil {
		return nil
	}
	parts := strings.Split(current, ";")
	filtered := make([]string, 0, len(parts))
	for _, p := range parts {
		if !strings.EqualFold(strings.TrimSpace(p), dir) {
			filtered = append(filtered, p)
		}
	}
	return k.SetStringValue("Path", strings.Join(filtered, ";"))
}
