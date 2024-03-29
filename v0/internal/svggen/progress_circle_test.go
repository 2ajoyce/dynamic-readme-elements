package svggen

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

// TestHandleProgressCircle tests the HandleProgressCircle function.
func TestHandleProgressCircle(t *testing.T) {
	router := gin.Default()
	router.GET("/progress/circle", HandleProgressCircle)

	testCases := []struct {
		name           string
		queryString    string
		expectedStatus int
		expectedInBody []string
	}{
		{
			name:           "Normal case",
			queryString:    "/progress/circle?size=103&percentage=58",
			expectedStatus: http.StatusOK,
			expectedInBody: []string{"<svg height=\"103px\" width=\"103px\"", "stroke-dasharray=\"132.9476, 96.2724\"", "58%"},
		},
		{
			name:           "Normal case - Full Template",
			queryString:    "/progress/circle?size=103&percentage=58",
			expectedStatus: http.StatusOK,
			expectedInBody: []string{
				`<svg height="103px" width="103px" viewBox="0 0 103 103" xmlns="http://www.w3.org/2000/svg">`,
				fmt.Sprintf("<circle cx=\"51.5\" cy=\"51.5\" r=\"36.5\" stroke=\"%s\" stroke-width=\"15\" fill=\"%s\" />", Colors.Grey, Colors.White),
				fmt.Sprintf("<circle cx=\"51.5\" cy=\"51.5\" r=\"36.5\" stroke=\"%s\" stroke-width=\"15\" fill=\"none\" stroke-dasharray=\"132.9476, 96.2724\" stroke-dashoffset=\"0\" transform=\"rotate(-90, 51.5, 51.5)\" />", Colors.Green),
				fmt.Sprintf("text x=\"51.5\" y=\"51.5\" font-size=\"20.6px\" dominant-baseline=\"central\" text-anchor=\"middle\" fill=\"%s\" font-family=\"Arial, Helvetica, sans-serif\" font-weight=\"bold\">58%%</text>", Colors.Black),
				`</svg>`,
			},
		},
		{
			name:           "Percentage below zero",
			queryString:    "/progress/circle?size=103&percentage=-10",
			expectedStatus: http.StatusOK,
			expectedInBody: []string{"stroke-dasharray=\"0, 229.22\"", "0%"},
		},
		{
			name:           "Percentage above 100",
			queryString:    "/progress/circle?size=103&percentage=150",
			expectedStatus: http.StatusOK,
			expectedInBody: []string{"stroke-dasharray=\"229.22, 0\"", "100%"},
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
			for _, str := range tc.expectedInBody {
				if !strings.Contains(body, str) {
					t.Errorf("Expected to find %s in response body", str)
				}
			}
		})
	}
}
