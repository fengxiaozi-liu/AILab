package transform

import (
	"strings"
	"testing"
)

func TestMarkdownSubAgentToTOML(t *testing.T) {
	input := []byte(`---
name: specify
description: Create a spec.
---

## Instructions

Run the specify flow.
`)
	got, err := MarkdownSubAgentToTOML(input)
	if err != nil {
		t.Fatalf("MarkdownSubAgentToTOML returned error: %v", err)
	}
	text := string(got)
	for _, want := range []string{
		`name = "specify"`,
		`description = "Create a spec."`,
		`prompt = "## Instructions`,
	} {
		if !strings.Contains(text, want) {
			t.Fatalf("TOML output %q does not contain %q", text, want)
		}
	}
}
