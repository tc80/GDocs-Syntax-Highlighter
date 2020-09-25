package style

import (
	"regexp"
	"strings"
)

const (
	darkTheme = "dark"

	// DefaultTheme is the default theme.
	DefaultTheme = darkTheme
)

var (
	// ThemeRegex is an optional directive to specify the theme of the code.
	// If not set, #theme=dark is assumed by default.
	ThemeRegex = regexp.MustCompile("^#theme=([\\w_]+)$")

	themes = map[string]bool{
		darkTheme: true,
	}
)

// GetTheme returns the theme and if it exists.
func GetTheme(theme string) (string, bool) {
	lower := strings.ToLower(theme)
	_, ok := themes[lower]
	return lower, ok
}
