package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func seq(start, end int) []int {
	log.Printf("seq function called with start: %d, end: %d", start, end) // Add this log
	sequence := make([]int, end-start+1)
	for i := range sequence {
		sequence[i] = start + i
	}
	return sequence
}

func mod(a, b int) int {
	return a % b
}

func div(a, b int) int {
	return a / b
}

func mult(a, b int) int {
	return a * b
}

func add(a, b int) int {
	return a + b
}

func hasElem(slice []int, elem int) bool {
	for _, v := range slice {
		if v == elem {
			return true
		}
	}
	return false
}

const (
	rectTemplateStr = `
		<svg width="{{.Width}}px" height="{{.Height}}px" xmlns="http://www.w3.org/2000/svg">
			<rect rx="3" ry="3" x="0" y="0" width="{{.Width}}px" height="{{.Height}}px" fill="#555" />
			<rect rx="3" ry="3" x="0" y="0" width="{{.FillWidth}}px" height="{{.Height}}px" fill="#4c1" />
			<text x="{{.TextX}}px" y="{{.TextY}}px" font-size="{{.FontSize}}px" dominant-baseline="central" text-anchor="middle" fill="#fff" font-family="Arial, Helvetica, sans-serif" font-weight="bold">{{.Percentage}}%</text>
		</svg>
		`
	circleTemplateStr = `
		<svg height="{{.Size}}px" width="{{.Size}}px" viewBox="0 0 {{.Size}} {{.Size}}" xmlns="http://www.w3.org/2000/svg">
			<circle cx="{{.Center}}" cy="{{.Center}}" r="{{.Radius}}" stroke="lightgrey" stroke-width="{{.StrokeWidth}}" fill="none" />
			<circle cx="{{.Center}}" cy="{{.Center}}" r="{{.Radius}}" stroke="#4c1" stroke-width="{{.StrokeWidth}}" fill="none" stroke-dasharray="{{.StrokeDasharray}}" stroke-dashoffset="0" transform="rotate(-90, {{.Center}}, {{.Center}})" />
			<text x="{{.Center}}" y="{{.Center}}" font-size="{{.FontSize}}px" dominant-baseline="central" text-anchor="middle" fill="black" font-family="Arial, Helvetica, sans-serif" font-weight="bold">{{.Percentage}}%</text>
		</svg>
		`
	calendarChartTemplateStr = `
		<svg width="370px" height="310px" xmlns="http://www.w3.org/2000/svg" font-family="Arial">
			<!-- Background with rounded corners and padding -->
			<rect x="5" y="5" width="360px" height="300px" fill="white" rx="15" /> <!-- Adjust as needed -->
		
			<!-- Header for month and year -->
			<text x="180" y="35" font-size="20" text-anchor="middle" fill="black">{{.MonthName}} {{.Year}}</text>
		
			{{- $startDay := .StartDay -}}
			{{- $daysInMonth := .DaysInMonth -}}
			{{- $progressDays := .ProgressDays -}}
		
			<!-- Generating the grid -->
			{{- range $i := seq 1 $daysInMonth -}}
				{{- $positionIndex := add (add $i $startDay) -1 -}}
				{{- $x := mod $positionIndex 7 -}}
				{{- $y := div $positionIndex 7 -}}
				{{- $isProgress := hasElem $progressDays $i -}}
				<rect x="{{add (mult $x 50) 15}}" y="{{add (mult $y 50) 45}}" width="40" height="40" fill="{{if $isProgress}}#4c1{{else}}#f0f0f0{{end}}" stroke="#ddd" />
				<text x="{{add (mult $x 50) 35}}" y="{{add (mult $y 50) 70}}" font-size="14" text-anchor="middle" fill="{{if $isProgress}}white{{else}}black{{end}}">{{$i}}</text>
			{{- end }}
		</svg>
		`
)

func main() {
	router := gin.Default()

	// Parse and store the templates
	rectTemplate := template.Must(template.New("rect").Parse(rectTemplateStr))
	circleTemplate := template.Must(template.New("circle").Parse(circleTemplateStr))
	funcMap := template.FuncMap{
		"seq":     seq,
		"mod":     mod,
		"div":     div,
		"mult":    mult,
		"add":     add,
		"hasElem": hasElem,
	}
	calendarChartTemplate := template.Must(template.New("calendarChart").Funcs(funcMap).Parse(calendarChartTemplateStr))

	// Route for a rectangular loading bar
	router.GET("/bar", func(c *gin.Context) {
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
			Width, Height, FillWidth, TextX, TextY, FontSize, Percentage int
		}{
			Width:      width,
			Height:     height,
			FillWidth:  fillWidth,
			TextX:      width / 2,
			TextY:      height / 2,
			FontSize:   height / 2,
			Percentage: percentage,
		}

		c.Writer.Header().Set("Content-Type", "image/svg+xml")
		err := rectTemplate.Execute(c.Writer, data)
		if err != nil {
			return
		}
	})

	// Route for a circular progress bar
	router.GET("/circle", func(c *gin.Context) {
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
		strokeDasharray := circumference * float64(percentage) / 100

		data := struct {
			Size, StrokeWidth, Percentage             int
			Radius, StrokeDasharray, FontSize, Center float64
		}{
			Size:            size,
			StrokeWidth:     strokeWidth,
			Percentage:      percentage,
			Radius:          radius,
			StrokeDasharray: strokeDasharray,
			FontSize:        float64(size) / 5,
			Center:          float64(size) / 2,
		}

		c.Writer.Header().Set("Content-Type", "image/svg+xml")
		err := circleTemplate.Execute(c.Writer, data)
		if err != nil {
			return
		}
	})

	router.GET("/calendar", func(c *gin.Context) {
		// Get the current date
		now := time.Now()

		// Get year and month from query parameters, defaulting to the current year and month
		yearParam := c.DefaultQuery("year", strconv.Itoa(now.Year()))
		monthParam := c.DefaultQuery("month", strconv.Itoa(int(now.Month())))

		// Get progressDays from query parameter, defaulting to the current day
		progressDaysParam := c.DefaultQuery("progressDays", strconv.Itoa(now.Day()))
		var progressDays []int
		for _, dayStr := range strings.Split(progressDaysParam, ",") {
			day, err := strconv.Atoi(dayStr)
			if err != nil {
				c.String(http.StatusBadRequest, fmt.Sprintf("Invalid progress day format: %s", dayStr))
				return
			}
			progressDays = append(progressDays, day)
		}

		// Convert year and month to appropriate types
		year, err := strconv.Atoi(yearParam)
		if err != nil {
			c.String(http.StatusBadRequest, "Invalid year format")
			return
		}

		monthInt, err := strconv.Atoi(monthParam)
		if err != nil {
			c.String(http.StatusBadRequest, "Invalid month format")
			return
		}

		month := time.Month(monthInt)
		if month < time.January || month > time.December {
			c.String(http.StatusBadRequest, "Month must be between 1 and 12")
			return
		}

		// Calculate the first day of the month and number of days in the month
		firstDayOfMonth := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
		startDay := int(firstDayOfMonth.Weekday())
		lastDayOfMonth := firstDayOfMonth.AddDate(0, 1, -1)
		daysInMonth := lastDayOfMonth.Day()

		// Prepare data for the template
		data := struct {
			Year, Month, StartDay, DaysInMonth int
			MonthName                          string
			ProgressDays                       []int
		}{
			Year:         year,
			Month:        int(month),
			MonthName:    month.String(),
			StartDay:     startDay,
			DaysInMonth:  daysInMonth,
			ProgressDays: progressDays,
		}

		// Execute the template and write the response
		c.Writer.Header().Set("Content-Type", "image/svg+xml")
		err = calendarChartTemplate.Execute(c.Writer, data)
		if err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("Error rendering calendar: %v", err))
			return
		}
	})

	err := router.Run(":8080")
	if err != nil {
		return
	}
}
