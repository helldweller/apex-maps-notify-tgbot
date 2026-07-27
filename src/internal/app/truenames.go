package app

// trueNames maps API map names to their joke aliases. Add new entries here.
var trueNames = map[string]string{
	"Storm Point": "Shit Point 💩",
	"E-District":  "Dristrict 💩💩💩",
}

// trueName returns the joke alias for name, or name unchanged if none is defined.
func trueName(name string) string {
	if alt, ok := trueNames[name]; ok {
		return alt
	}
	return name
}
