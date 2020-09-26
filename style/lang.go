package style

import (
	"regexp"
	"strings"
)

// Language represents a programming language.
type Language struct {
	Name      string
	Format    FormatFunc
	Shortcuts []*Shortcut
	Themes    map[string]*Theme
}

var (
	// LangRegex is the regex for the optional directive
	// to specify the language of the code.
	// If not set, #lang=go is assumed by default.
	LangRegex = regexp.MustCompile("^#lang=([\\w_]+)$")

	goLang = &Language{
		Name:      "Go",
		Format:    FormatGo,
		Shortcuts: []*Shortcut{doubleQuotes, singleQuotes, goMainShortcut},
		Themes: map[string]*Theme{
			darkTheme: {
				DocBackground:    DarkThemeBackground,
				CodeForeground:   DarkThemeForeground,
				CodeBackground:   DarkThemeBackground,
				CodeHighlight:    Transparent,
				ConfigForeground: White,
				ConfigBackground: Black,
				ConfigHighlight:  Transparent,
				ConfigFont:       courierNew,
				ConfigFontSize:   11,
				ConfigItalics:    true,
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
					{go3, DarkThemeGreenCyan},
					{go4, DarkThemeYellow},
					{go5, DarkThemePaleGreen},
				},
			},
			lightTheme: {
				DocBackground:    White,
				CodeForeground:   Black,
				CodeBackground:   White,
				CodeHighlight:    Transparent,
				ConfigForeground: Black,
				ConfigBackground: LightGray,
				ConfigHighlight:  Transparent,
				ConfigFont:       courierNew,
				ConfigFontSize:   11,
				ConfigItalics:    true,
				Ranges: []*Range{
					{"//", "\n", LightThemeDarkGreen},
					{"/*", "*/", LightThemeDarkGreen},
					{"\"", "\"", LightThemeDarkRed},
					{"'", "'", LightThemeDarkRed},
					{"`", "`", LightThemeDarkRed},
				},
				Keywords: []Keyword{
					{go1, LightThemePink},
					{go2, Blue},
					{go3, LightThemeGreenCyan},
					{go4, LightThemeStrawYellow},
					{go5, LightThemePaleGreen},
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
