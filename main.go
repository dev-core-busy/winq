package main

import (
	"fmt"
	"os"
	"os/exec"

	tea "github.com/charmbracelet/bubbletea"
)

// restartAfterUpdate + restartExecPath: gesetzt nach erfolgreichem auto-update.
// Der Pfad wird VOR dem Rename gesichert, damit /proc/self/exe nicht auf einen
// gelöschten Inode zeigt wenn syscall.Exec aufgerufen wird.
var restartAfterUpdate bool
var restartExecPath string

func main() {
	p := tea.NewProgram(
		newModel(),
		tea.WithAltScreen(),
	)
	if _, err := p.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "Fehler:", err)
		os.Exit(1)
	}

	if restartAfterUpdate && restartExecPath != "" {
		cmd := exec.Command(restartExecPath, os.Args[1:]...)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Start(); err != nil {
			fmt.Fprintln(os.Stderr, "winq Neustart fehlgeschlagen:", err)
			fmt.Fprintln(os.Stderr, "Bitte manuell neu starten:", restartExecPath)
			os.Exit(1)
		}
		os.Exit(0)
	}
}
