package main

import (
	"os"

	"golang.org/x/sys/windows"
)

// terminalSize gibt die aktuelle Konsolengröße zurück (Spalten, Zeilen).
// Wird in cmdInitSize() genutzt um den Startup-Hänger auf Windows zu vermeiden.
func terminalSize() (int, int) {
	var info windows.ConsoleScreenBufferInfo
	h := windows.Handle(os.Stdout.Fd())
	if err := windows.GetConsoleScreenBufferInfo(h, &info); err != nil {
		return 0, 0
	}
	w := int(info.Window.Right - info.Window.Left + 1)
	ht := int(info.Window.Bottom - info.Window.Top + 1)
	return w, ht
}
