package app

import (
	"fmt"
	"regexp"
	"time"

	"package/main/internal/apexapi"
)

// md2regex matches every character that must be escaped in Telegram MarkdownV2.
var md2regex = regexp.MustCompile(`(\_|\*|\[|\]|\(|\)|\~|\>|\#|\+|\-|\=|\||\{|\}|\.|\!)`)

// escapeMarkdownV2 escapes reserved MarkdownV2 characters in s.
func escapeMarkdownV2(s string) string {
	return md2regex.ReplaceAllString(s, `\$1`)
}

// hoursAndMinutes splits a duration into whole hours and the remaining minutes.
func hoursAndMinutes(d time.Duration) (hours, minutes int) {
	hours = int(d.Hours())
	minutes = int(d.Minutes()) - hours*60
	return hours, minutes
}

// formatMode renders the notification text for a single game mode.
//
// prefix is inserted between "Карта " and "сейчас" ("" for battle royale,
// "в рейтинге " for ranked, "в ltm " for ltm). When withAsset is true the
// current map asset link is appended on its own line.
func formatMode(prefix string, maps apexapi.Maps, now time.Time, withAsset bool) string {
	nextStartAt := time.Unix(maps.Next.Start, 0)
	nextEndAt := time.Unix(maps.Next.End, 0)

	currentHours, currentMinutes := hoursAndMinutes(nextStartAt.Sub(now))
	nextHours, nextMinutes := hoursAndMinutes(nextEndAt.Sub(nextStartAt))

	text := fmt.Sprintf(
		"Карта %sсейчас *%s* и продлится *%dч %dм*\nСледующая карта *%s* и продлится *%dч %dм*",
		prefix,
		escapeMarkdownV2(maps.Current.Map),
		currentHours, currentMinutes,
		escapeMarkdownV2(maps.Next.Map),
		nextHours, nextMinutes,
	)

	if withAsset {
		text += fmt.Sprintf("\n[](%s)", maps.Current.Asset)
	}

	return text
}
