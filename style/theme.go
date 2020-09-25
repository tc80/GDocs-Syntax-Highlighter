package style

import (
	"regexp"
	"strings"

	"google.golang.org/api/docs/v1"
)

const (
	// DarkTheme is the dark theme.
	DarkTheme = "dark"

	// DefaultTheme is the default theme.
	DefaultTheme = DarkTheme
)

var (
	// ThemeRegex is an optional directive to specify the theme of the code.
	// If not set, #theme=dark is assumed by default.
	ThemeRegex = regexp.MustCompile("^#theme=([\\w_]+)$")

	themes = map[string]bool{
		DarkTheme: true,
	}
)

// Theme represents a language's keywords, comments
// and associated colors for a particular theme.
type Theme struct {
	DocBackground    *docs.Color
	CodeForeground   *docs.Color
	CodeBackground   *docs.Color
	ConfigForeground *docs.Color
	ConfigBackground *docs.Color
	ConfigFont       string
	ConfigFontSize   float64
	ConfigItalics    bool
	Ranges           []*Range
	Keywords         []Keyword
}

// Range represents an area of text that will receive the same color.
// For instance, a comment.
// For now, there is no notion of precedence.
type Range struct {
	StartSymbol string
	EndSymbol   string
	Color       *docs.Color
}

// GetTheme returns the theme and if it exists.
func GetTheme(theme string) (string, bool) {
	lower := strings.ToLower(theme)
	_, ok := themes[lower]
	return lower, ok
}
