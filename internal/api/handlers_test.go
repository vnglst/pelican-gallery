package api

import (
	"testing"
	"time"

	"pelican-gallery/internal/models"
)

func TestEffectiveMaxTokens(t *testing.T) {
	messages := []models.Message{{Role: "user", Content: "draw a pelican"}}
	tests := []struct {
		name  string
		model models.ModelInfo
		in    int
		want  int
	}{
		{name: "requested below limits", model: models.ModelInfo{ContextLength: 16385, MaxCompletionTokens: 16000}, in: 8000, want: 8000},
		{name: "completion limit", model: models.ModelInfo{ContextLength: 100000, MaxCompletionTokens: 12000}, in: 50000, want: 12000},
		{name: "context reserves input", model: models.ModelInfo{ContextLength: 1000}, in: 50000, want: 934},
		{name: "prompt fills context", model: models.ModelInfo{ContextLength: 10}, in: 50000, want: 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := effectiveMaxTokens(tt.in, tt.model, messages); got != tt.want {
				t.Fatalf("effectiveMaxTokens(%d) = %d, want %d", tt.in, got, tt.want)
			}
		})
	}
}

func TestGenerationTimeout(t *testing.T) {
	t.Setenv("OPENROUTER_TIMEOUT_SECONDS", "12")
	if got := generationTimeout(); got != 12*time.Second {
		t.Fatalf("generationTimeout() = %v, want 12s", got)
	}

	t.Setenv("OPENROUTER_TIMEOUT_SECONDS", "invalid")
	if got := generationTimeout(); got != defaultGenerationTimeout {
		t.Fatalf("generationTimeout() = %v, want default %v", got, defaultGenerationTimeout)
	}
}

func TestGenerationTimeoutMessage(t *testing.T) {
	t.Setenv("OPENROUTER_TIMEOUT_SECONDS", "600")
	want := "Generation was stopped after 10 minutes while waiting for the AI provider. The model may be busy or producing a long response; please retry or choose a faster model."
	if got := generationTimeoutMessage(); got != want {
		t.Fatalf("generationTimeoutMessage() = %q, want %q", got, want)
	}
}
