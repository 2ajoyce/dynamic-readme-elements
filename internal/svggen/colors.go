package svggen

// ColorSet defines a set of color constants
type ColorSet struct {
	ProgressActive   string
	ProgressInactive string
	White            string
	Black            string
}

// Colors holds the application-wide color constants
var Colors = ColorSet{
	ProgressActive:   "#44CC11",
	ProgressInactive: "#7A7A7A",
	White:            "white",
	Black:            "black",
}
