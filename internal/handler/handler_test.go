package handler

import (
	"testing"
)

func TestReservedPaths(t *testing.T) {
	expected := []string{"privacy", "endpoints", "health", "static", "robots.txt", "sitemap.xml", "favicon.ico"}
	for _, path := range expected {
		if !reservedPaths[path] {
			t.Errorf("expected %q to be reserved", path)
		}
	}
}

func TestReservedPathsRejectsUnknown(t *testing.T) {
	unreserved := []string{"my-endpoint", "brave-lazy-fox", "test", ""}
	for _, path := range unreserved {
		if reservedPaths[path] {
			t.Errorf("expected %q to NOT be reserved", path)
		}
	}
}

func TestDefaultBodyForStatus(t *testing.T) {
	tests := []struct {
		status int
		delay  int
		want   string
	}{
		{200, 0, `{"status": "ok"}`},
		{201, 0, `{"status": "created", "id": "example-id"}`},
		{400, 0, `{"error": "bad_request", "message": "Invalid payload"}`},
		{500, 0, `Internal Server Error`},
		{200, 35000, ``},
		{999, 0, `{"status": "ok"}`}, // unknown falls back to default
	}

	for _, tt := range tests {
		got := defaultBodyForStatus(tt.status, tt.delay)
		if got != tt.want {
			t.Errorf("defaultBodyForStatus(%d, %d) = %q, want %q", tt.status, tt.delay, got, tt.want)
		}
	}
}

func TestResponsePresetsComplete(t *testing.T) {
	if len(responsePresets) == 0 {
		t.Fatal("responsePresets should not be empty")
	}

	for _, p := range responsePresets {
		if p.Label == "" {
			t.Errorf("preset with status %d has empty label", p.Status)
		}
		if p.DefaultBody == "" && p.Delay == 0 {
			t.Errorf("preset %q has empty default body", p.Label)
		}
	}
}
