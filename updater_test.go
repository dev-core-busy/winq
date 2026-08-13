package main

import "testing"

func TestSemverGT(t *testing.T) {
	cases := []struct {
		rel, cur string
		want     bool
	}{
		// Kein Downgrade: ein älteres Release darf nie als neuer gelten
		{"v1.0.7", "v1.0.10", false},
		{"v1.0.10", "v1.0.7", true},
		// Zweistellig gegen einstellig – hier scheitert ein String-Vergleich
		{"v1.0.10", "v1.0.9", true},
		{"v1.0.9", "v1.0.10", false},
		// Gleichstand löst kein Update aus
		{"v1.0.11", "v1.0.11", false},
		{"1.0.11", "v1.0.11", false},
		// Minor/Major schlagen Patch
		{"v1.1.0", "v1.0.99", true},
		{"v2.0.0", "v1.9.9", true},
		// Fehlende Segmente werden als 0 ergänzt
		{"v1.1", "v1.1.0", false},
		{"v1.1.1", "v1.1", true},
		// Vorabversions-Suffix: nur die führenden Ziffern zählen
		{"v1.0.12-rc1", "v1.0.11", true},
		{"v1.0.11-rc1", "v1.0.12", false},
	}
	for _, c := range cases {
		if got := semverGT(c.rel, c.cur); got != c.want {
			t.Errorf("semverGT(%q, %q) = %v, erwartet %v", c.rel, c.cur, got, c.want)
		}
	}
}
