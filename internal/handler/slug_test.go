package handler

import (
	"strings"
	"testing"
)

func TestGenerateSlug(t *testing.T) {
	slug := generateSlug()
	if slug == "" {
		t.Fatal("generateSlug() returned empty string")
	}

	parts := strings.Split(slug, "-")
	if len(parts) != 3 {
		t.Errorf("expected 3 hyphen-separated words, got %d: %q", len(parts), slug)
	}
}

func TestGenerateSlugUniqueness(t *testing.T) {
	seen := make(map[string]bool)
	for i := 0; i < 100; i++ {
		slug := generateSlug()
		if seen[slug] {
			t.Errorf("duplicate slug after %d generations: %q", i, slug)
		}
		seen[slug] = true
	}
}
