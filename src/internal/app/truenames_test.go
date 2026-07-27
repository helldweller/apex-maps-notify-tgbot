package app

import "testing"

func TestTrueName(t *testing.T) {
	cases := map[string]string{
		"Storm Point": "Shit Point 💩",
		"E-District":  "Dristrict 💩💩💩",
		"Olympus":     "Olympus",
		"":            "",
	}
	for in, want := range cases {
		if got := trueName(in); got != want {
			t.Errorf("trueName(%q) = %q, want %q", in, got, want)
		}
	}
}
