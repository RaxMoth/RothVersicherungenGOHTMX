package handler

import (
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandler(t *testing.T) {
	tests := []struct {
		path string
		want int
	}{
		{"/", 200},
		{"/roth-versicherungen/firmenkunden/cyber-police", 200},
		{"/roth-finanz/sterbegeldversicherung", 200},
		{"/static/css/output.css", 200},
		{"/does-not-exist", 404},
	}
	for _, tt := range tests {
		rec := httptest.NewRecorder()
		Handler(rec, httptest.NewRequest("GET", tt.path, nil))
		if rec.Code != tt.want {
			t.Errorf("GET %s = %d, want %d", tt.path, rec.Code, tt.want)
		}
	}

	// The home page must render real content, not an error fallback.
	rec := httptest.NewRecorder()
	Handler(rec, httptest.NewRequest("GET", "/", nil))
	if !strings.Contains(rec.Body.String(), "Finanzmakler in Langen") {
		t.Error("home page is missing expected content")
	}
}
