package svggen

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func seq(start, end int) []int {
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

const calendarChartTemplateStr = `
		<svg width="370px" height="310px" xmlns="http://www.w3.org/2000/svg" font-family="Arial">
			<!-- Background with rounded corners and padding -->
			<rect x="5" y="5" width="360px" height="300px" fill="white" rx="15" />

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

func HandleCalendar(c *gin.Context) {
	funcMap := template.FuncMap{
		"seq":     seq,
		"mod":     mod,
		"div":     div,
		"mult":    mult,
		"add":     add,
		"hasElem": hasElem,
	}
	calendarChartTemplate := template.Must(template.New("calendarChart").Funcs(funcMap).Parse(calendarChartTemplateStr))

	// Get year and month from query parameters with defaults to the current year and month
	yearParam := c.DefaultQuery("year", strconv.Itoa(time.Now().Year()))
	monthParam := c.DefaultQuery("month", strconv.Itoa(int(time.Now().Month())))

	// Get progressDays from query parameter
	progressDaysParam := c.DefaultQuery("progressDays", "")
	var progressDays []int
	if progressDaysParam != "" {
		for _, dayStr := range strings.Split(progressDaysParam, ",") {
			day, err := strconv.Atoi(dayStr)
			if err != nil {
				c.String(http.StatusBadRequest, fmt.Sprintf("Invalid progress day format: %s", dayStr))
				return
			}
			progressDays = append(progressDays, day)
		}
	} else {
		// Default to the current day if no progressDays are provided
		progressDays = append(progressDays, time.Now().Day())
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
	fmt.Printf("Year: %d, Month: %d (%s), Days in Month: %d, Start Day: %d, Progress Days: %v\n",
		data.Year, data.Month, data.MonthName, data.DaysInMonth, data.StartDay, data.ProgressDays)
	err = calendarChartTemplate.Execute(c.Writer, data)
	if err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("Error rendering calendar: %v", err))
		return
	}
}
