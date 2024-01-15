package svggen

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"html/template"
	"math"
	"strconv"
)

const gaugeChartTemplateStr = `
<svg height="{{.Center}}px" width="{{.Size}}px" viewBox="0 0 {{.Size}} {{.Center}}" xmlns="http://www.w3.org/2000/svg">
	{{ range .PieSections }}
		<path d="{{.Path}}" fill="{{.FillColor}}"/>
	{{ end }}
    <path d="M{{.Center}},{{.Center}} L{{mult .Center 0.5}},{{.Center}} A{{mult .Center 0.5}},{{mult .Center 0.5}} 0 1,1 {{mult .Center 1.5}},{{.Center}} Z" fill="{{.ColorWhite}}"/>
    <polygon points="{{.Needle.X1}},{{.Needle.Y1}} {{.Needle.X2}},{{.Needle.Y2}} {{.Needle.X3}},{{.Needle.Y3}}" fill="{{.ColorBlack}}" />
    <path d="M{{.Center}},{{.Center}} L{{mult .Center 0.8}},{{.Center}} A{{mult .Center 0.2}},{{mult .Center 0.2}} 0 1,1 {{mult .Center 1.2}},{{.Center}} Z" fill="{{.ColorBlack}}"/>
    <path d="M{{.Center}},{{.Center}} L{{mult .Center 0.9}},{{.Center}} A{{mult .Center 0.1}},{{mult .Center 0.1}} 0 1,1 {{mult .Center 1.1}},{{.Center}} Z" fill="{{.ColorGrey}}"/>
</svg>
`

type Needle struct {
	X1, Y1 float64 // First point of the triangle
	X2, Y2 float64 // Second point of the triangle
	X3, Y3 float64 // Third point of the triangle
}

type PieSection struct {
	FillColor string
	Path      string // Path data with M, L, A, and Z commands
}

func radians(degrees float64) float64 {
	return degrees * (math.Pi / 180)
}

const angle = 36.0 // Angle for dividing the circle into sections

func HandleProgressGauge(c *gin.Context) {
	funcMap := template.FuncMap{
		"mult": multFloat64,
	}
	gaugeChartTemplate := template.Must(template.New("gaugeChart").Funcs(funcMap).Parse(gaugeChartTemplateStr))

	// Retrieve parameters or default
	width, _ := strconv.Atoi(c.DefaultQuery("width", "100"))
	percentage, _ := strconv.Atoi(c.DefaultQuery("percentage", "0"))
	if width <= 0 {
		width = 100
	}
	percentage = clamp(percentage, 0, 100)

	// Calculate center and needle position based on percentage
	center := float64(width) / 2
	needle := calculateNeedlePosition(center, percentage)

	// Create instances of PieSection for each section
	pieSections := []PieSection{
		{Colors.Red, createPiePath(center, center, 180, 180+angle)},
		{Colors.Orange, createPiePath(center, center, 180+angle, 180+2*angle)},
		{Colors.Yellow, createPiePath(center, center, 180+2*angle, 180+3*angle)},
		{Colors.LightGreen, createPiePath(center, center, 180+3*angle, 180+4*angle)},
		{Colors.Green, createPiePath(center, center, 180+4*angle, 180+5*angle)},
	}

	data := struct {
		ColorWhite, ColorBlack, ColorGrey string
		Size, Percentage                  int
		Center, NeedleX, NeedleY          float64
		Needle                            Needle
		PieSections                       []PieSection
	}{
		ColorWhite:  Colors.White,
		ColorBlack:  Colors.Black,
		ColorGrey:   Colors.Grey,
		Size:        width,
		Percentage:  percentage,
		Center:      center,
		Needle:      needle,
		PieSections: pieSections,
	}

	c.Writer.Header().Set("Content-Type", "image/svg+xml")
	if err := gaugeChartTemplate.Execute(c.Writer, data); err != nil {
		return
	}
}

func createPiePath(center, radius, startAngle, endAngle float64) string {
	startX := center + radius*math.Cos(radians(startAngle))
	startY := center + radius*math.Sin(radians(startAngle))
	endX := center + radius*math.Cos(radians(endAngle))
	endY := center + radius*math.Sin(radians(endAngle))

	path := fmt.Sprintf("M %f,%f A %f,%f 0 0 1 %f,%f L %f,%f L %f,%f Z", startX, startY, radius, radius, endX, endY, center, center, startX, startY)

	return path
}

func calculateNeedlePosition(center float64, percentage int) Needle {
	// Convert percentage to angle in radians
	angle := float64(percentage)/100*math.Pi + math.Pi

	// Calculate needle end point coordinates
	needleLength := center * 0.45
	needleX := center + needleLength*math.Cos(angle)
	needleY := center + needleLength*math.Sin(angle)

	// Calculate triangle points
	triangleBaseLength := center * 0.75 // Length of the base of the triangle
	triangleAngle := 15.0               // Angle of the triangle point in degrees
	degrees := math.Pi / 180

	// Calculate the coordinates of the triangle points
	needleX1 := needleX - (triangleBaseLength * 0.5 * math.Cos(angle-(triangleAngle*degrees)))
	needleY1 := needleY - (triangleBaseLength * 0.5 * math.Sin(angle-(triangleAngle*degrees)))

	needleX2 := needleX - (triangleBaseLength * 0.5 * math.Cos(angle+(triangleAngle*degrees)))
	needleY2 := needleY - (triangleBaseLength * 0.5 * math.Sin(angle+(triangleAngle*degrees)))

	needleX3 := needleX + (needleLength * math.Cos(angle))
	needleY3 := needleY + (needleLength * math.Sin(angle))

	needle := Needle{
		X1: needleX1, Y1: needleY1,
		X2: needleX2, Y2: needleY2,
		X3: needleX3, Y3: needleY3,
	}

	return needle
}
