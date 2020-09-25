package style

import (
	"regexp"
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
	// FontRegex is an optional directive to specify the font of the code.
	// If not set, #font=courier_new is assumed by default.
	FontRegex = regexp.MustCompile("^#font=([\\w_]+)$")

	// FontSizeRegex is an optional directive to specify the font size of the code.
	// If not set, #size=11 is assumed by default.
	FontSizeRegex = regexp.MustCompile("^#size=(\\d+(\\.\\d+)?)$")

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
