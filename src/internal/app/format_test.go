package app

import (
	"testing"
	"time"

	"package/main/internal/apexapi"
)

func TestEscapeMarkdownV2(t *testing.T) {
	cases := map[string]string{
		"World's Edge":      "World's Edge",                       // apostrophe is not reserved
		"E-District":        `E\-District`,                        // hyphen is reserved
		"a.b":               `a\.b`,                               // dot
		"x!y":               `x\!y`,                               // bang
		`_*[]()~>#+-=|{}.!`: `\_\*\[\]\(\)\~\>\#\+\-\=\|\{\}\.\!`, // full reserved set
	}
	for in, want := range cases {
		if got := escapeMarkdownV2(in); got != want {
			t.Errorf("escapeMarkdownV2(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestHoursAndMinutes(t *testing.T) {
	cases := []struct {
		d     time.Duration
		wantH int
		wantM int
	}{
		{90 * time.Minute, 1, 30},
		{135 * time.Minute, 2, 15},
		{0, 0, 0},
		{59 * time.Minute, 0, 59},
	}
	for _, c := range cases {
		h, m := hoursAndMinutes(c.d)
		if h != c.wantH || m != c.wantM {
			t.Errorf("hoursAndMinutes(%v) = %dh %dm, want %dh %dm", c.d, h, m, c.wantH, c.wantM)
		}
	}
}

// fixtureMaps returns a Maps value with a clean 1h30m gap until the next map
// and a 2h15m next-map duration, relative to the returned "now".
func fixtureMaps(currentMap, nextMap, asset string) (apexapi.Maps, time.Time) {
	now := time.Unix(1000, 0)
	nextStart := int64(1000 + 90*60)     // now + 1h30m
	nextEnd := nextStart + int64(135*60) // + 2h15m
	m := apexapi.Maps{
		Current: apexapi.Map{Start: 0, End: nextStart, Map: currentMap, Asset: asset},
		Next:    apexapi.Map{Start: nextStart, End: nextEnd, Map: nextMap},
	}
	return m, now
}

func TestFormatModeBattleRoyale(t *testing.T) {
	m, now := fixtureMaps("E-District", "Olympus", "https://example.com/we.png")
	want := "Карта сейчас *E\\-District* и продлится *1ч 30м*\n" +
		"Следующая карта *Olympus* и продлится *2ч 15м*\n" +
		"[](https://example.com/we.png)"
	if got := formatMode("", m, now, true); got != want {
		t.Errorf("formatMode battle royale:\n got: %q\nwant: %q", got, want)
	}
}

func TestFormatModeRanked(t *testing.T) {
	m, now := fixtureMaps("E-District", "Olympus", "ignored")
	want := "Карта в рейтинге сейчас *E\\-District* и продлится *1ч 30м*\n" +
		"Следующая карта *Olympus* и продлится *2ч 15м*"
	if got := formatMode("в рейтинге ", m, now, false); got != want {
		t.Errorf("formatMode ranked:\n got: %q\nwant: %q", got, want)
	}
}

func TestFormatModeLtm(t *testing.T) {
	m, now := fixtureMaps("E-District", "Olympus", "ignored")
	want := "Карта в ltm сейчас *E\\-District* и продлится *1ч 30м*\n" +
		"Следующая карта *Olympus* и продлится *2ч 15м*"
	if got := formatMode("в ltm ", m, now, false); got != want {
		t.Errorf("formatMode ltm:\n got: %q\nwant: %q", got, want)
	}
}
