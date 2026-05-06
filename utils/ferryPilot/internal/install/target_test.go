package install

import (
	"path/filepath"
	"testing"

	"ferrypilot/internal/config"
)

func TestResolveTargetGlobalCodex(t *testing.T) {
	plan, err := ResolveTarget(TargetOptions{
		FileMap: targetTestFileMap(),
		Mode:    "global",
		Target:  "codex",
		HomeDir: "home",
	})
	if err != nil {
		t.Fatalf("ResolveTarget returned error: %v", err)
	}
	if plan.Agent != "codex" {
		t.Fatalf("Agent = %q, want codex", plan.Agent)
	}
	want := filepath.Join("home", ".codex", "skills")
	if plan.Mappings[0].Destination != want {
		t.Fatalf("Destination = %q, want %q", plan.Mappings[0].Destination, want)
	}
}

func TestResolveTargetProjectCursor(t *testing.T) {
	plan, err := ResolveTarget(TargetOptions{
		FileMap: targetTestFileMap(),
		Mode:    "project",
		Target:  "cursor",
		WorkDir: "project",
	})
	if err != nil {
		t.Fatalf("ResolveTarget returned error: %v", err)
	}
	want := filepath.Join("project", ".cursor", "skills")
	if plan.Mappings[0].Destination != want {
		t.Fatalf("Destination = %q, want %q", plan.Mappings[0].Destination, want)
	}
}

func TestResolveTargetAdditionalAgents(t *testing.T) {
	tests := []struct {
		name  string
		mode  string
		agent string
		root  string
		want  string
	}{
		{name: "claude global", mode: "global", agent: "claude", root: "home", want: filepath.Join("home", ".claude", "skills")},
		{name: "copilot project", mode: "project", agent: "copilot", root: "project", want: filepath.Join("project", ".github", "skills")},
		{name: "gemini global", mode: "global", agent: "gemini", root: "home", want: filepath.Join("home", ".gemini", "antigravity", "skills")},
		{name: "gemini project", mode: "project", agent: "gemini", root: "project", want: filepath.Join("project", ".agent", "skills")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			options := TargetOptions{
				FileMap: targetTestFileMap(),
				Mode:    tt.mode,
				Target:  tt.agent,
			}
			if tt.mode == "global" {
				options.HomeDir = tt.root
			} else {
				options.WorkDir = tt.root
			}
			plan, err := ResolveTarget(options)
			if err != nil {
				t.Fatalf("ResolveTarget returned error: %v", err)
			}
			if plan.Mappings[0].Destination != tt.want {
				t.Fatalf("Destination = %q, want %q", plan.Mappings[0].Destination, tt.want)
			}
		})
	}
}

func TestResolveTargetRejectsUnknownTarget(t *testing.T) {
	_, err := ResolveTarget(TargetOptions{
		FileMap: targetTestFileMap(),
		Mode:    "project",
		Target:  "unknown",
		WorkDir: "project",
	})
	if err == nil {
		t.Fatal("ResolveTarget returned nil error")
	}
}

func targetTestFileMap() config.FileMap {
	return config.FileMap{
		DefaultTarget: "codex",
		Targets: map[string]config.Target{
			"codex": {
				Global: []config.Mapping{
					{Source: "skills", Destination: ".codex/skills"},
					{Source: "sub-agents", Destination: ".codex/agents", Transform: true},
				},
			},
			"cursor": {
				Project: []config.Mapping{
					{Source: "skills", Destination: ".cursor/skills"},
					{Source: "sub-agents", Destination: ".cursor/agents", Transform: true},
				},
			},
			"claude": {
				Global: []config.Mapping{
					{Source: "skills", Destination: ".claude/skills"},
					{Source: "sub-agents", Destination: ".claude/agents"},
				},
			},
			"copilot": {
				Project: []config.Mapping{
					{Source: "skills", Destination: ".github/skills"},
					{Source: "sub-agents", Destination: ".github/agents"},
				},
			},
			"gemini": {
				Global: []config.Mapping{
					{Source: "skills", Destination: ".gemini/antigravity/skills"},
					{Source: "sub-agents", Destination: ".gemini/antigravity/workflows"},
				},
				Project: []config.Mapping{
					{Source: "skills", Destination: ".agent/skills"},
					{Source: "sub-agents", Destination: ".agent/workflows"},
				},
			},
		},
	}
}
