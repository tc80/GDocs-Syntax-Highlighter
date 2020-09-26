package style

import (
	"strings"
)

const (
	courierNew = "Courier New"
	consolas   = "Consolas"

	// DefaultFont is the default font.
	DefaultFont = courierNew

	// DefaultFontSize is the default font size.
	DefaultFontSize float64 = 11
)

var (
	fonts = map[string]string{
		"courier_new": courierNew,
		"consolas":    consolas,
	}
)

// GetFont attempts to get the Google Docs name
// of a font from its alias.
func GetFont(font string) (string, bool) {
	f, ok := fonts[strings.ToLower(font)]
	return f, ok
}
