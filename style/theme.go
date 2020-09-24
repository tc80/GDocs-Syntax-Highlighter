package style

import "strings"

const (
	darkTheme = "dark"

	// DefaultTheme is the default theme.
	DefaultTheme = darkTheme
)

var (
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
