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

const calendarChartTemplateStr = `
	<svg width="370px" height="{{.Height}}px" xmlns="http://www.w3.org/2000/svg" font-family="Arial" xmlns:xlink="http://www.w3.org/1999/xlink">
		<script type="text/ecmascript">
		<![CDATA[
		var currentYear = {{.Year}};
		var currentMonth = {{.Month}};
		var progressDaysMap = {};
		{{- range $day := .ProgressDays -}}
		progressDaysMap[{{$day}}] = true;
		{{- end }}

		function getDaysInMonth(year, month) {
			return new Date(year, month, 0).getDate();
		}

		function getStartDay(year, month) {
			return new Date(year, month - 1, 1).getDay();
		}

		function getMonthName(month) {
			var months = ["January", "February", "March", "April", "May", "June",
				"July", "August", "September", "October", "November", "December"];
			return months[month - 1];
		}

		function updateCalendar() {
			var daysInMonth = getDaysInMonth(currentYear, currentMonth);
			var startDay = getStartDay(currentYear, currentMonth);
			
			// Update header
			var header = document.getElementById("monthHeader");
			header.textContent = getMonthName(currentMonth) + " " + currentYear;
			
			// Clear existing calendar grid
			var calendarGrid = document.getElementById("calendarGrid");
			while (calendarGrid.firstChild) {
				calendarGrid.removeChild(calendarGrid.firstChild);
			}
			
			// Generate new calendar grid
			for (var i = 1; i <= daysInMonth; i++) {
				var positionIndex = i + startDay - 1;
				var x = (positionIndex % 7) * 50 + 15;
				var y = Math.floor(positionIndex / 7) * 50 + 45;
				var isProgress = progressDaysMap[i] === true;
				
				var rect = document.createElementNS("http://www.w3.org/2000/svg", "rect");
				rect.setAttribute("x", x);
				rect.setAttribute("y", y);
				rect.setAttribute("width", 40);
				rect.setAttribute("height", 40);
				rect.setAttribute("fill", isProgress ? "#4c1" : "#f0f0f0");
				rect.setAttribute("stroke", "#ddd");
				calendarGrid.appendChild(rect);
				
				var text = document.createElementNS("http://www.w3.org/2000/svg", "text");
				text.setAttribute("x", x + 20);
				text.setAttribute("y", y + 25);
				text.setAttribute("font-size", 14);
				text.setAttribute("text-anchor", "middle");
				text.setAttribute("fill", isProgress ? "white" : "black");
				text.textContent = i;
				calendarGrid.appendChild(text);
			}

			// Update SVG height if needed
			var height = 310;
			if (startDay + daysInMonth > 35) {
				height = 360;
			}
			document.documentElement.setAttribute("height", height + "px");
			document.getElementById("bgRect").setAttribute("height", (height - 10) + "px");
		}

		function prevMonth() {
			currentMonth--;
			if (currentMonth < 1) {
				currentMonth = 12;
				currentYear--;
			}
			updateCalendar();
		}

		function nextMonth() {
			currentMonth++;
			if (currentMonth > 12) {
				currentMonth = 1;
				currentYear++;
			}
			updateCalendar();
		}
		]]>
		</script>

		<!-- Background with rounded corners and padding -->
		<rect id="bgRect" x="5" y="5" width="360px" height="{{add .Height -10}}px" fill="white" rx="15" />

		<!-- Header for month and year -->
		<text id="monthHeader" x="180" y="35" font-size="20" text-anchor="middle" fill="black">{{.MonthName}} {{.Year}}</text>

		<!-- Navigation Buttons -->
		<g id="prevButton" onmousedown="prevMonth();" style="cursor:pointer;">
			<rect x="15" y="15" width="60" height="25" fill="#007bff" rx="5" />
			<text x="45" y="32" font-size="12" text-anchor="middle" fill="white">← Prev</text>
		</g>
		<g id="nextButton" onmousedown="nextMonth();" style="cursor:pointer;">
			<rect x="295" y="15" width="60" height="25" fill="#007bff" rx="5" />
			<text x="325" y="32" font-size="12" text-anchor="middle" fill="white">Next →</text>
		</g>

		{{- $startDay := .StartDay -}}
		{{- $daysInMonth := .DaysInMonth -}}
		{{- $progressDays := .ProgressDays -}}

		<!-- Generating the grid -->
		<g id="calendarGrid">
		{{- range $i := seq 1 $daysInMonth -}}
			{{- $positionIndex := add (add $i $startDay) -1 -}}
			{{- $x := mod $positionIndex 7 -}}
			{{- $y := div $positionIndex 7 -}}
			{{- $isProgress := hasElem $progressDays $i -}}
			<rect x="{{add (mult $x 50) 15}}" y="{{add (mult $y 50) 45}}" width="40" height="40" fill="{{if $isProgress}}#4c1{{else}}#f0f0f0{{end}}" stroke="#ddd" />
			<text x="{{add (mult $x 50) 35}}" y="{{add (mult $y 50) 70}}" font-size="14" text-anchor="middle" fill="{{if $isProgress}}white{{else}}black{{end}}">{{$i}}</text>
		{{- end }}
		</g>
	</svg>
	`

func HandleCalendar(c *gin.Context) {
	funcMap := template.FuncMap{
		"seq":     seq,
		"mod":     mod,
		"div":     div,
		"mult":    multInt,
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

	height := 310 // If the month has 5 weeks
	if startDay+daysInMonth > 35 {
		height = 360 // If the month has 6 weeks
	}

	// Prepare data for the template
	data := struct {
		Year, Month, StartDay, DaysInMonth int
		MonthName                          string
		ProgressDays                       []int
		Height                             int
	}{
		Year:         year,
		Month:        int(month),
		MonthName:    month.String(),
		StartDay:     startDay,
		DaysInMonth:  daysInMonth,
		ProgressDays: progressDays,
		Height:       height,
	}

	// Execute the template and write the response
	c.Writer.Header().Set("Content-Type", "image/svg+xml")
	err = calendarChartTemplate.Execute(c.Writer, data)
	if err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("Error rendering calendar: %v", err))
		return
	}
}
