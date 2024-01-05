package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func main() {
	router := gin.Default()

	// Route for a rectangular loading bar
	router.GET("/bar", func(c *gin.Context) {
		// Default values
		width := 200
		height := 30
		percentage := 0

		// Parse width, height, and percentage from query parameters
		if w, err := strconv.Atoi(c.DefaultQuery("width", "200")); err == nil {
			width = w
		}
		if h, err := strconv.Atoi(c.DefaultQuery("height", "30")); err == nil {
			height = h
		}
		if p, err := strconv.Atoi(c.DefaultQuery("percentage", "0")); err == nil {
			percentage = p
		}

		// Ensure percentage is within 0-100 range
		if percentage < 0 {
			percentage = 0
		} else if percentage > 100 {
			percentage = 100
		}

		// Calculate fill width
		fillWidth := (width * percentage) / 100

		// Generate SVG
		svg := fmt.Sprintf(
			`<svg width="%dpx" height="%dpx" xmlns="http://www.w3.org/2000/svg">
				<rect rx="3" ry="3" x="0" y="0" width="%dpx" height="%dpx" fill="#555" />
				<rect rx="3" ry="3" x="0" y="0" width="%dpx" height="%dpx" fill="#4c1" />
				<text x="%dpx" y="%dpx" font-size="%dpx" dominant-baseline="central" text-anchor="middle" fill="#fff" font-family="Arial, Helvetica, sans-serif" font-weight="bold">%d%%</text>
			</svg>`, width, height, width, height, fillWidth, height, width/2, height/2, height/2, percentage)

		c.Writer.Header().Set("Content-Type", "image/svg+xml")
		c.String(http.StatusOK, svg)
	})

	// Route for a circular progress bar
	router.GET("/circle", func(c *gin.Context) {
		// Default values
		size := 100
		percentage := 0

		// Parse size and percentage from query parameters
		if s, err := strconv.Atoi(c.DefaultQuery("size", "100")); err == nil {
			size = s
		}
		if p, err := strconv.Atoi(c.DefaultQuery("percentage", "0")); err == nil {
			percentage = p
		}

		// Ensure percentage is within 0-100 range
		if percentage < 0 {
			percentage = 0
		} else if percentage > 100 {
			percentage = 100
		}

		// Calculate font size
		fontSize := float64(size) / 5

		// Calculate stroke-dasharray for circular progress
		strokeWidth := 15
		radius := float64(size)/(2) - float64(strokeWidth)
		circumference := 2 * 3.14 * radius
		strokeDasharray := fmt.Sprintf("%f %f", (circumference*float64(percentage))/100, circumference)

		svg := fmt.Sprintf(
			`<svg height="%dpx" width="%dpx" viewBox="0 0 %d %d" xmlns="http://www.w3.org/2000/svg">
            <circle cx="%d" cy="%d" r="%f" stroke="lightgrey" stroke-width="%d" fill="none" />
            <circle cx="%d" cy="%d" r="%f" stroke="#4c1" stroke-width="%d" fill="none" stroke-dasharray="%s" stroke-dashoffset="0" transform="rotate(-90, %d, %d)" />
            <text x="%d" y="%d" font-size="%fpx" dominant-baseline="central" text-anchor="middle" fill="black" font-family="Arial, Helvetica, sans-serif" font-weight="bold">%d%%</text>
        </svg>`, size, size, size, size, size/2, size/2, radius, strokeWidth, size/2, size/2, radius, strokeWidth, strokeDasharray, size/2, size/2, size/2, size/2, fontSize, percentage)

		c.Writer.Header().Set("Content-Type", "image/svg+xml")
		c.String(http.StatusOK, svg)
	})

	err := router.Run(":8080")
	if err != nil {
		return
	}
}
