package svggen

import (
	"github.com/gin-gonic/gin"
	"html/template"
	"strconv"
)

const circleTemplateStr = `

	<svg height="{{.Size}}px" width="{{.Size}}px" viewBox="0 0 {{.Size}} {{.Size}}" xmlns="http://www.w3.org/2000/svg">
		<circle cx="{{.Center}}" cy="{{.Center}}" r="{{.Radius}}" stroke="lightgrey" stroke-width="{{.StrokeWidth}}" fill="white" />
		<circle cx="{{.Center}}" cy="{{.Center}}" r="{{.Radius}}" stroke="#4c1" stroke-width="{{.StrokeWidth}}" fill="none" stroke-dasharray="{{.StrokeDasharrayFilled}}, {{.StrokeDasharrayUnfilled}}" stroke-dashoffset="0" transform="rotate(-90, {{.Center}}, {{.Center}})" />
		<text x="{{.Center}}" y="{{.Center}}" font-size="{{.FontSize}}px" dominant-baseline="central" text-anchor="middle" fill="black" font-family="Arial, Helvetica, sans-serif" font-weight="bold">{{.Percentage}}%</text>
	</svg>
	`

func HandleProgressCircle(c *gin.Context) {
	circleTemplate := template.Must(template.New("circle").Parse(circleTemplateStr))

	size, _ := strconv.Atoi(c.DefaultQuery("size", "100"))
	percentage, _ := strconv.Atoi(c.DefaultQuery("percentage", "0"))

	if percentage < 0 {
		percentage = 0
	} else if percentage > 100 {
		percentage = 100
	}

	strokeWidth := 15
	radius := float64(size)/2 - float64(strokeWidth)
	circumference := 2 * 3.14 * radius
	strokeDasharrayFilled := circumference * float64(percentage) / 100
	strokeDasharrayUnfilled := circumference - strokeDasharrayFilled

	data := struct {
		Size, StrokeWidth, Percentage                                            int
		Radius, StrokeDasharrayFilled, StrokeDasharrayUnfilled, FontSize, Center float64
	}{
		Size:                    size,
		StrokeWidth:             strokeWidth,
		Percentage:              percentage,
		Radius:                  radius,
		StrokeDasharrayFilled:   strokeDasharrayFilled,
		StrokeDasharrayUnfilled: strokeDasharrayUnfilled,
		FontSize:                float64(size) / 5,
		Center:                  float64(size) / 2,
	}

	c.Writer.Header().Set("Content-Type", "image/svg+xml")
	err := circleTemplate.Execute(c.Writer, data)
	if err != nil {
		return
	}
}
