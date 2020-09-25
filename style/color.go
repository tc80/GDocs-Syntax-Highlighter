package style

import (
	"encoding/hex"
	"log"

	"google.golang.org/api/docs/v1"
)

var (
	// Note that some of the following hex codes are taken/modified from the VSCode themes found here:
	// https://github.com/microsoft/vscode/tree/master/extensions/theme-defaults/themes

	// Transparent color.
	Transparent *docs.Color

	// White color.
	White = getColorFromHex("FFFFFF")

	// Black color.
	Black = getColorFromHex("000000")

	// Blue color.
	Blue = getColorFromHex("0000FF")

	// LightGray color.
	LightGray = getColorFromHex("F3F3F3")

	// LightThemePink is VSCode's light theme pink color.
	LightThemePink = getColorFromHex("AF00DB")

	// LightThemeGreenCyan is VSCode's light theme green-cyan color.
	LightThemeGreenCyan = getColorFromHex("267F99")

	// LightThemeStrawYellow is VSCode's light theme straw-yellow color.
	LightThemeStrawYellow = getColorFromHex("795E26")

	// LightThemePaleGreen is VSCode's light theme pale-green color.
	LightThemePaleGreen = getColorFromHex("098658")

	// LightThemeDarkGreen is VSCode's light theme dark-green color.
	LightThemeDarkGreen = getColorFromHex("008000")

	// LightThemeDarkRed is VSCode's light theme dark red color.
	LightThemeDarkRed = getColorFromHex("A31515")

	// DarkThemeBackground is VSCode's dark theme background color (dark gray).
	DarkThemeBackground = getColorFromHex("1E1E1E")

	// DarkThemeForeground is VSCode's dark theme default foreground color (light gray/white).
	DarkThemeForeground = getColorFromHex("D4D4D4")

	// DarkThemeYellow is VSCode's dark theme yellow color.
	DarkThemeYellow = getColorFromHex("DCDCAA")

	// DarkThemeGreenCyan is VSCode's dark theme green-cyan color.
	DarkThemeGreenCyan = getColorFromHex("4EC9B0")

	// DarkThemePaleGreen is VSCode's dark theme pale-green color.
	DarkThemePaleGreen = getColorFromHex("B5CEA8")

	// DarkThemeDarkGreen is VSCode's dark theme dark-green color.
	DarkThemeDarkGreen = getColorFromHex("6A9955")

	// DarkThemePink is VSCode's dark theme pink color.
	DarkThemePink = getColorFromHex("C586C0")

	// DarkThemeLightBlue is VSCode's dark theme light blue color.
	DarkThemeLightBlue = getColorFromHex("9CDCFE")

	// DarkThemeBlue is VSCode's dark theme blue color.
	DarkThemeBlue = getColorFromHex("4FC1FF")

	// DarkThemeDarkBlue is VSCode's dark theme dark blue color.
	DarkThemeDarkBlue = getColorFromHex("569CD6")

	// DarkThemeLightRedOrange is VSCode's dark theme dark light red-orange color.
	DarkThemeLightRedOrange = getColorFromHex("CE9178")

	// DarkThemeLightRed is VSCode's dark theme light red color.
	DarkThemeLightRed = getColorFromHex("D16969")

	// DarkThemeStrawYellow is VSCode's dark theme straw-yellow color.
	DarkThemeStrawYellow = getColorFromHex("D7BA7D")
)

// Gets an RGB color from red, green, blue values in [0.0, 1.0].
func getColorFromRGB(r, g, b float64) *docs.Color {
	return &docs.Color{
		RgbColor: &docs.RgbColor{
			Red:   r,
			Green: g,
			Blue:  b,
		},
	}
}

// Gets an RGB color from a hex code.
func getColorFromHex(h string) *docs.Color {
	b, err := hex.DecodeString(h)
	if err != nil {
		log.Fatalf("Failed to decode hex `%s`: %s\n", h, err)
	}
	return &docs.Color{
		RgbColor: &docs.RgbColor{
			Red:   float64(b[0]) / 255,
			Green: float64(b[1]) / 255,
			Blue:  float64(b[2]) / 255,
		},
	}
}
