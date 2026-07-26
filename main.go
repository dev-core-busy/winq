package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

// Nach erfolgreichem Auto-Update: winq beendet sich und druckt diese Meldung
// ins normale Terminal. Kein Auto-Restart — Windows-Console-Handles werden beim
// Prozesswechsel unzuverlässig, was die neue Instanz taub für Tastatureingaben macht.
var restartAfterUpdate bool
var updatedToVersion string

func main() {
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

	if restartAfterUpdate {
		fmt.Printf("\n  ✓ winq %s installiert — bitte winq neu starten.\n\n", updatedToVersion)
	}
}
