package main

import (
	"bytes"
	"html"
	"html/template"
	"strings"
	"testing"
)

func TestFrameKeepsGeneratedMarkupInert(t *testing.T) {
	tmpl, err := template.ParseFiles("templates/frame.html")
	if err != nil {
		t.Fatal(err)
	}
	// Even malformed/truncated model output must not hide the next artwork or
	// install its IDs and styles in the gallery document.
	for _, source := range []string{
		`<svg><style>.bg{fill:red}</style><rect id="bg"/></svg>`,
		`<svg><rect fill="red"/><!-- unfinished comment`,
		`<svg><text>Tom & Jerry</text></svg></template><style>body{display:none}</style>`,
	} {
		var output bytes.Buffer
		if err := tmpl.ExecuteTemplate(&output, "frame", template.HTML(source)); err != nil {
			t.Fatal(err)
		}
		body := output.String()
		start := strings.Index(body, "<template data-svg-source>") + len("<template data-svg-source>")
		end := strings.Index(body, "</template>")
		if start < len("<template data-svg-source>") || end < start {
			t.Fatalf("missing inert SVG source: %s", body)
		}
		if got := html.UnescapeString(body[start:end]); got != source {
			t.Fatalf("SVG source changed: got %q, want %q", got, source)
		}
		if strings.Contains(body, "<svg") || strings.Contains(body, "<style") {
			t.Fatalf("generated markup leaked into the document: %s", body)
		}
	}
}
