package svggen

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

// TestHandleProgressBar tests the HandleProgressBar function.
func TestHandleProgressBar(t *testing.T) {
	router := gin.Default()
	router.GET("/progress/bar", HandleProgressBar)

	testCases := []struct {
		name           string
		queryString    string
		expectedStatus int
		expectedInBody []string
	}{
		{
			name:           "Normal case",
			queryString:    "/progress/bar?width=221&height=33&percentage=54",
			expectedStatus: http.StatusOK,
			expectedInBody: []string{"width=\"221px\"", "height=\"33px\"", fmt.Sprintf("fill=\"%s\"", Colors.Green), "54%"},
		},
		{
			name:           "Normal case - Full Template",
			queryString:    "/progress/bar?width=221&height=33&percentage=54",
			expectedStatus: http.StatusOK,
			expectedInBody: []string{
				`<svg width="221px" height="33px" xmlns="http://www.w3.org/2000/svg">`,
				fmt.Sprintf("<rect rx=\"3\" ry=\"3\" x=\"0\" y=\"0\" width=\"221px\" height=\"33px\" fill=\"%s\" />", Colors.Grey),
				fmt.Sprintf("<rect rx=\"3\" ry=\"3\" x=\"0\" y=\"0\" width=\"119px\" height=\"33px\" fill=\"%s\" />", Colors.Green),
				fmt.Sprintf("<text x=\"110px\" y=\"16px\" font-size=\"16px\" dominant-baseline=\"central\" text-anchor=\"middle\" fill=\"%s\" font-family=\"Arial, Helvetica, sans-serif\" font-weight=\"bold\">54%%</text>", Colors.White),
				`</svg>`,
			},
		},
		{
			name:           "Percentage below zero",
			queryString:    "/progress/bar?width=200&height=30&percentage=-10",
			expectedStatus: http.StatusOK,
			expectedInBody: []string{"width=\"200px\"", "height=\"30px\"", fmt.Sprintf("fill=\"%s\"", Colors.Grey), "0%"},
		},
		{
			name:           "Percentage above 100",
			queryString:    "/progress/bar?width=200&height=30&percentage=150",
			expectedStatus: http.StatusOK,
			expectedInBody: []string{"width=\"200px\"", "height=\"30px\"", fmt.Sprintf("fill=\"%s\"", Colors.Green), "100%"},
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
