package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	// Heruntergeladenes Update anwenden bevor die TUI startet.
	// Zu diesem Zeitpunkt ist die EXE noch nicht gesperrt → Rename klappt auf Windows.
	applyPendingUpdate()

	p := tea.NewProgram(
		newModel(),
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
		// Bubbletea liest unter Windows über die Console-Input-API und kann
		// Bracketed-Paste-Sequenzen (ESC[200~ … ESC[201~) nicht auswerten —
		// sie würden als Text im Eingabefeld landen. Strg+V wird stattdessen
		// selbst behandelt (siehe cmdReadClipboard in update.go).
		tea.WithoutBracketedPaste(),
	)
	if _, err := p.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "Fehler:", err)
		os.Exit(1)
	}
}
