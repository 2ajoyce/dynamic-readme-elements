package svggen

import (
	"github.com/gin-gonic/gin"
	"html/template"
	"strconv"
)

const rectTemplateStr = `
		<svg width="{{.Width}}px" height="{{.Height}}px" xmlns="http://www.w3.org/2000/svg">
			<rect rx="3" ry="3" x="0" y="0" width="{{.Width}}px" height="{{.Height}}px" fill="{{.ColorInactive}}" />
			<rect rx="3" ry="3" x="0" y="0" width="{{.FillWidth}}px" height="{{.Height}}px" fill="{{.ColorActive}}" />
			<text x="{{.TextX}}px" y="{{.TextY}}px" font-size="{{.FontSize}}px" dominant-baseline="central" text-anchor="middle" fill="{{.ColorWhite}}" font-family="Arial, Helvetica, sans-serif" font-weight="bold">{{.Percentage}}%</text>
		</svg>
		`

func HandleProgressBar(c *gin.Context) {
	rectTemplate := template.Must(template.New("rect").Parse(rectTemplateStr))
	width, _ := strconv.Atoi(c.DefaultQuery("width", "200"))
	height, _ := strconv.Atoi(c.DefaultQuery("height", "30"))
	percentage, _ := strconv.Atoi(c.DefaultQuery("percentage", "0"))

	// Ensure percentage is within 0-100 range
	if percentage < 0 {
		percentage = 0
	} else if percentage > 100 {
		percentage = 100
	}

	fillWidth := (width * percentage) / 100

	data := struct {
		ColorActive, ColorInactive, ColorWhite                       string
		Width, Height, FillWidth, TextX, TextY, FontSize, Percentage int
	}{
		ColorActive:   Colors.Green,
		ColorInactive: Colors.Grey,
		ColorWhite:    Colors.White,
		Width:         width,
		Height:        height,
		FillWidth:     fillWidth,
		TextX:         width / 2,
		TextY:         height / 2,
		FontSize:      height / 2,
		Percentage:    percentage,
	}

	c.Writer.Header().Set("Content-Type", "image/svg+xml")
	err := rectTemplate.Execute(c.Writer, data)
	if err != nil {
		return
	}
}
