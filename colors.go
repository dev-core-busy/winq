package main

// setupColors ist unter Windows ein No-op — Windows Terminal unterstützt
// ANSI-Farben nativ, keine .bashrc-Anpassung nötig.
func setupColors() (msg string, isErr bool) {
	return "Windows Terminal unterstützt ANSI-Farben nativ — keine Anpassung nötig.", false
}
