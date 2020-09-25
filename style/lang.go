package style

import (
	"strings"
)

// Language represents a programming language.
type Language struct {
	Name      string
	Format    FormatFunc
	Shortcuts map[string]string
	Themes    map[string]*Theme
}

var (
	goLang = &Language{
		Name:      "Go",
		Format:    FormatGo,
		Shortcuts: map[string]string{},
		Themes: map[string]*Theme{
			DarkTheme: {
				Foreground: DarkThemeForeground,
				Background: DarkThemeBackground,
				Ranges: []*Range{
					{"//", "\n", DarkThemeDarkGreen},
					{"/*", "*/", DarkThemeDarkGreen},
					{"\"", "\"", DarkThemeLightRedOrange},
					{"'", "'", DarkThemeLightRedOrange},
					{"`", "`", DarkThemeLightRedOrange},
				},
				Keywords: []Keyword{
					{go1, DarkThemePink},
					{go2, DarkThemeDarkBlue},
					{go3, DarkThemeDarkBlue},
					{go4, DarkThemeGreenCyan},
					{go5, DarkThemeYellow},
				},
			},
		},
	}
	languages = map[string]*Language{
		"go": goLang,
	}
)

// GetLanguage attempts to get a Language
// from a case insensitive string.
func GetLanguage(lang string) (*Language, bool) {
	l, ok := languages[strings.ToLower(lang)]
	return l, ok
}

// GetDefaultLanguage gets the default Language
// if the directive is not set.
func GetDefaultLanguage() *Language {
	return goLang
}
