package config

import "testing"

func TestIsOpenSourceModelRecognizesGemmaWithoutMetadata(t *testing.T) {
	for _, modelID := range []string{"google/gemma-3n-e4b-it", "google/gemma-3-27b-it:free"} {
		if !IsOpenSourceModel(modelID, "") {
			t.Errorf("IsOpenSourceModel(%q, empty metadata) = false, want true", modelID)
		}
	}
}
