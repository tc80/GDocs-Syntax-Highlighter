package style

// Color represents an RGB color
type Color struct {
	Red   float64 // the red value from 0.0 to 1.0
	Green float64 // the green value from 0.0 to 1.0
	Blue  float64 // the blue value from 0.0 to 1.0
}

var (
	// Black is the color black
	Black = Color{}
	// White is the color white
	White = Color{1, 1, 1}
	// Red is the color red
	Red = Color{1, 0, 0}
	// Green is the color green
	Green = Color{0, 1, 0}
	// Blue is the color blue
	Blue = Color{0, 0, 1}
)
