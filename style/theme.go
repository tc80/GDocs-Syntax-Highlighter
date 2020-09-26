package style

import (
	"regexp"
	"strings"

	"google.golang.org/api/docs/v1"
)

const (
	// The dark theme.
	darkTheme = "dark"

	// The light theme.
	lightTheme = "light"

	// DefaultTheme is the default theme.
	DefaultTheme = lightTheme
)

var (
	// ThemeRegex is an optional directive to specify the theme of the code.
	// If not set, #theme=dark is assumed by default.
	ThemeRegex = regexp.MustCompile("^#theme=([\\w_]+)$")

	themes = map[string]bool{
		darkTheme:  true,
		lightTheme: true,
	}
)

// Theme represents a language's keywords, comments
// and associated colors for a particular theme.
type Theme struct {
	DocBackground    *docs.Color
	CodeForeground   *docs.Color
	CodeBackground   *docs.Color
	CodeHighlight    *docs.Color
	ConfigForeground *docs.Color
	ConfigBackground *docs.Color
	ConfigHighlight  *docs.Color
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
