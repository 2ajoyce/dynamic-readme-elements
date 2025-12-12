package svggen

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

func TestHandleCalendar(t *testing.T) {
	router := gin.Default()
	router.GET("/calendar", HandleCalendar)

	testCases := []struct {
		name           string
		queryString    string
		expectedStatus int
		expectInBody   []string
	}{
		{
			name:           "Valid parameters - Spot Checks",
			queryString:    "/calendar?year=2023&month=1&progressDays=1,15",
			expectedStatus: http.StatusOK,
			expectInBody: []string{
				`<svg width="370px" height="310px"`,                                                  // Height for 5wk month
				`xmlns="http://www.w3.org/2000/svg"`,                                                 // SVG namespace
				`<rect id="bgRect" x="5" y="5" width="360px" height="300px" fill="white" rx="15" />`, // Background
				`<text id="monthHeader"`,                                                             // Header element
				`January 2023</text>`,                                                                // Header text
				`onmousedown="prevMonth();"`,                                                         // Previous button
				`onmousedown="nextMonth();"`,                                                         // Next button
				// Check for a few specific days, including one that is marked as progress and one that is not
				`<rect x="15" y="45" width="40" height="40" fill="#4c1" stroke="#ddd" />`,       // Day 1 with progress
				`<text x="35" y="70" font-size="14" text-anchor="middle" fill="white">1</text>`, // Text for Day 1
				`<rect x="65" y="45" width="40" height="40" fill="#f0f0f0" stroke="#ddd" />`,    // Day 2 without progress
				`<text x="85" y="70" font-size="14" text-anchor="middle" fill="black">2</text>`, // Text for Day 2
			},
		},
		{
			name:           "Check height when month has 6 weeks",
			queryString:    "/calendar?year=2023&month=4",
			expectedStatus: http.StatusOK,
			expectInBody: []string{
				`<svg width="370px" height="360px"`,
				`xmlns="http://www.w3.org/2000/svg"`,
			},
		},
		{
			name:           "Check height when month has ALMOST 6 weeks",
			queryString:    "/calendar?year=2024&month=8",
			expectedStatus: http.StatusOK,
			expectInBody: []string{
				`<svg width="370px" height="310px"`,
				`xmlns="http://www.w3.org/2000/svg"`,
			},
		},
		{
			name:           "Valid parameters - Check Full Response",
			queryString:    "/calendar?year=2023&month=1&progressDays=1,15",
			expectedStatus: http.StatusOK,
			expectInBody: []string{
				`<svg width="370px" height="310px"`,
				`xmlns="http://www.w3.org/2000/svg"`,
				`<rect id="bgRect" x="5" y="5" width="360px" height="300px" fill="white" rx="15" />`,
				`January 2023</text>`,
				`<rect x="15" y="45" width="40" height="40" fill="#4c1" stroke="#ddd" />`,
				`<text x="35" y="70" font-size="14" text-anchor="middle" fill="white">1</text>`,
				`<rect x="65" y="45" width="40" height="40" fill="#f0f0f0" stroke="#ddd" />`,
				`<text x="85" y="70" font-size="14" text-anchor="middle" fill="black">2</text>`,
				`<rect x="15" y="145" width="40" height="40" fill="#4c1" stroke="#ddd" />`,
				`<text x="35" y="170" font-size="14" text-anchor="middle" fill="white">15</text>`,
				`<text x="135" y="270" font-size="14" text-anchor="middle" fill="black">31</text>`,
				`</svg>`,
				`<script type="text/ecmascript">`,
				`function updateCalendar()`,
				`function prevMonth()`,
				`function nextMonth()`,
			},
		},
		{
			name:           "Invalid year",
			queryString:    "/calendar?year=abc&month=1",
			expectedStatus: http.StatusBadRequest,
			expectInBody:   []string{"Invalid year format"},
		},
		{
			name:           "Invalid month",
			queryString:    "/calendar?year=2023&month=13",
			expectedStatus: http.StatusBadRequest,
			expectInBody:   []string{"Month must be between 1 and 12"},
		},
		{
			name:           "Default to current year and month",
			queryString:    "/calendar",
			expectedStatus: http.StatusOK,
			expectInBody:   []string{time.Now().Format("January"), strconv.Itoa(time.Now().Year())},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", tc.queryString, nil)
			router.ServeHTTP(w, req)

			if w.Code != tc.expectedStatus {
				t.Errorf("Expected status %d, got %d", tc.expectedStatus, w.Code)
			}

			body := w.Body.String()

			for _, str := range tc.expectInBody {
				if !strings.Contains(body, str) {
					t.Errorf("Expected to find %s in response body", str)
				}
			}
		})
	}
}
