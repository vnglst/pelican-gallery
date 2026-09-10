package pages

import (
	"sort"
	"testing"
)

func TestModelVersionLess(t *testing.T) {
	models := []string{
		"google/gemini-3-pro-preview",
		"google/gemini-2.5-pro",
		"google/gemini-2.0-flash-001",
		"google/gemini-3.1-pro-preview",
	}

	sort.Slice(models, func(i, j int) bool {
		return modelVersionLess(models[i], models[j])
	})

	want := []string{
		"google/gemini-2.0-flash-001",
		"google/gemini-2.5-pro",
		"google/gemini-3-pro-preview",
		"google/gemini-3.1-pro-preview",
	}
	for i := range want {
		if models[i] != want[i] {
			t.Fatalf("models[%d] = %q, want %q", i, models[i], want[i])
		}
	}
}
