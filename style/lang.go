package style

import (
	"strings"

	"google.golang.org/api/docs/v1"
)

// Language represents a programming language.
type Language struct {
	Name      string
	Format    FormatFunc
	Shortcuts map[string]string
	Themes    map[string]*Theme
}

// Theme represents a language's keywords, comments
// and associated colors for a particular theme.
type Theme struct {
	Foreground *docs.Color
	Background *docs.Color
	Ranges     []*Range
	Keywords   map[string]*docs.Color
}

// Range represents an area of text that will receive the same color.
// For instance, a comment.
// For now, there is no notion of precedence.
type Range struct {
	StartSymbol string
	EndSymbol   string
	Color       *docs.Color
}

var (
	goLang = &Language{
		Name:      "Go",
		Format:    FormatGo,
		Shortcuts: map[string]string{},
		Themes: map[string]*Theme{
			darkTheme: &Theme{
				Foreground: DarkThemeForeground,
				Background: DarkThemeBackground,
				Ranges: []*Range{
					&Range{"//", "\n", DarkThemePaleGreen},
					&Range{"/*", "*/", DarkThemePaleGreen},
					&Range{"\"", "\"", DarkThemeLightRedOrange},
					&Range{"`", "`", DarkThemeLightRedOrange},
				},
				Keywords: map[string]*docs.Color{
					"break":       DarkThemePink,
					"case":        DarkThemePink,
					"continue":    DarkThemePink,
					"default":     DarkThemePink,
					"defer":       DarkThemePink,
					"else":        DarkThemePink,
					"fallthrough": DarkThemePink,
					"for":         DarkThemePink,
					"go":          DarkThemePink,
					"goto":        DarkThemePink,
					"if":          DarkThemePink,
					"range":       DarkThemePink,
					"return":      DarkThemePink,
					"select":      DarkThemePink,
					"switch":      DarkThemePink,
					"chan":        DarkThemeDarkBlue,
					"const":       DarkThemeDarkBlue,
					"func":        DarkThemeDarkBlue,
					"interface":   DarkThemeDarkBlue,
					"map":         DarkThemeDarkBlue,
					"struct":      DarkThemeDarkBlue,
					"true":        DarkThemeDarkBlue,
					"false":       DarkThemeDarkBlue,
					"nil":         DarkThemeDarkBlue,
					"iota":        DarkThemeDarkBlue,
					"package":     DarkThemeDarkBlue,
					"type":        DarkThemeDarkBlue,
					"import":      DarkThemeDarkBlue,
					"var":         DarkThemeDarkBlue,
					"bool":        DarkThemeGreenCyan,
					"byte":        DarkThemeGreenCyan,
					"error":       DarkThemeGreenCyan,
					"complex64":   DarkThemeGreenCyan,
					"complex128":  DarkThemeGreenCyan,
					"float32":     DarkThemeGreenCyan,
					"float64":     DarkThemeGreenCyan,
					"int":         DarkThemeGreenCyan,
					"int8":        DarkThemeGreenCyan,
					"int16":       DarkThemeGreenCyan,
					"int32":       DarkThemeGreenCyan,
					"int64":       DarkThemeGreenCyan,
					"uint8":       DarkThemeGreenCyan,
					"uint16":      DarkThemeGreenCyan,
					"uint32":      DarkThemeGreenCyan,
					"uint64":      DarkThemeGreenCyan,
					"rune":        DarkThemeGreenCyan,
					"string":      DarkThemeGreenCyan,
					"uintptr":     DarkThemeGreenCyan,
					"append":      DarkThemeYellow,
					"cap":         DarkThemeYellow,
					"close":       DarkThemeYellow,
					"complex":     DarkThemeYellow,
					"copy":        DarkThemeYellow,
					"delete":      DarkThemeYellow,
					"imag":        DarkThemeYellow,
					"len":         DarkThemeYellow,
					"make":        DarkThemeYellow,
					"new":         DarkThemeYellow,
					"panic":       DarkThemeYellow,
					"print":       DarkThemeYellow,
					"println":     DarkThemeYellow,
					"real":        DarkThemeYellow,
					"recover":     DarkThemeYellow,
					//"'":           DarkThemeLightRedOrange,
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
