package style

import (
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
	themes = map[string]bool{
		darkTheme:  true,
		lightTheme: true,
	}
)

// Theme represents a language's keywords, comments
// and associated colors for a particular theme.
// For now, by default all code is not bolded, not underlined,
// not in italics, not in small caps, and not strikethrough.
// Since underlines are used in directives, at the moment
// they can not be removed from directive headers/footers.
type Theme struct {
	DocBackground       *docs.Color
	CodeForeground      *docs.Color
	CodeBackground      *docs.Color
	CodeHighlight       *docs.Color
	ConfigForeground    *docs.Color
	ConfigBackground    *docs.Color
	ConfigHighlight     *docs.Color
	ConfigFont          string
	ConfigFontSize      float64
	ConfigItalics       bool
	ConfigBold          bool
	ConfigSmallCaps     bool
	ConfigStrikethrough bool
	Ranges              []*Range
	Keywords            []Keyword
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
