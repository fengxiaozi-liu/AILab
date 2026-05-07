package app

import (
	"strings"
	"testing"

	"ferryaide/internal/config"
)

func TestChooseTargetUsesExplicitTarget(t *testing.T) {
	fileMap := testFileMap()

	target, err := chooseTarget(Options{Target: "gemini"}, fileMap)
	if err != nil {
		t.Fatalf("chooseTarget returned error: %v", err)
	}
	if target != "gemini" {
		t.Fatalf("target = %q, want gemini", target)
	}
}

func TestChooseTargetPromptsWhenTargetMissing(t *testing.T) {
	fileMap := testFileMap()
	var output strings.Builder

	target, err := chooseTarget(Options{
		Input:  strings.NewReader("2\n"),
		Output: &output,
	}, fileMap)
	if err != nil {
		t.Fatalf("chooseTarget returned error: %v", err)
	}
	if target != "claude" {
		t.Fatalf("target = %q, want claude", target)
	}
	if !strings.Contains(output.String(), "Select a target agent:") {
		t.Fatalf("prompt output %q does not contain target selection heading", output.String())
	}
}

func TestChooseTargetRequiresTargetWhenNonInteractive(t *testing.T) {
	_, err := chooseTarget(Options{}, testFileMap())
	if err == nil {
		t.Fatal("chooseTarget returned nil error")
	}
	if !strings.Contains(err.Error(), "target agent is required") {
		t.Fatalf("error = %q, want target agent required", err.Error())
	}
}

func TestChooseTargetSingleTargetDoesNotPrompt(t *testing.T) {
	fileMap := config.FileMap{
		DefaultTarget: "codex",
		Targets: map[string]config.Target{
			"codex": {},
		},
	}

	target, err := chooseTarget(Options{}, fileMap)
	if err != nil {
		t.Fatalf("chooseTarget returned error: %v", err)
	}
	if target != "codex" {
		t.Fatalf("target = %q, want codex", target)
	}
}

func testFileMap() config.FileMap {
	return config.FileMap{
		DefaultTarget: "codex",
		Targets: map[string]config.Target{
			"codex":  {},
			"claude": {},
			"gemini": {},
		},
	}
}
