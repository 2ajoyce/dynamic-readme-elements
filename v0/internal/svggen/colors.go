package svggen

// ColorSet defines a set of color constants
type ColorSet struct {
	Green      string
	Grey       string
	White      string
	Black      string
	Red        string
	Orange     string
	Yellow     string
	LightGreen string
}

// Colors holds the application-wide color constants
var Colors = ColorSet{
	Green:      "#44CC11",
	Grey:       "#7A7A7A",
	White:      "white",
	Black:      "black",
	Red:        "red",
	Orange:     "orange",
	Yellow:     "yellow",
	LightGreen: "#99F255",
}
