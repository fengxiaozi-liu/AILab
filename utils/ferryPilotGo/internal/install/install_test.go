package install

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"ferrypilot/internal/assets"
	"ferrypilot/internal/config"
)

func TestPackageCopiesAndTransforms(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "AISupport", "speckit", "skills", "one", "SKILL.md"), "# skill")
	writeFile(t, filepath.Join(root, "AISupport", "speckit", "sub-agents", "specify.md"), `---
name: specify
description: Specify things.
---

Prompt body.
`)

	source, err := assets.NewSource(context.Background(), config.DataSource{Type: "local", Path: root}, "")
	if err != nil {
		t.Fatal(err)
	}
	project := filepath.Join(root, "project")
	plan, err := ResolveTarget(TargetOptions{
		FileMap: installTestFileMap(),
		Mode:    "project",
		Target:  "codex",
		WorkDir: project,
	})
	if err != nil {
		t.Fatal(err)
	}

	installed, err := Package(context.Background(), source, "speckit", plan)
	if err != nil {
		t.Fatalf("Package returned error: %v", err)
	}
	if len(installed) != 2 {
		t.Fatalf("installed len = %d, want 2", len(installed))
	}
	if _, err := os.Stat(filepath.Join(project, ".codex", "skills", "one", "SKILL.md")); err != nil {
		t.Fatalf("skill not copied: %v", err)
	}
	got, err := os.ReadFile(filepath.Join(project, ".codex", "agents", "specify.toml"))
	if err != nil {
		t.Fatalf("agent not transformed: %v", err)
	}
	if !strings.Contains(string(got), `name = "specify"`) {
		t.Fatalf("transformed content = %q", string(got))
	}
}

func installTestFileMap() config.FileMap {
	return config.FileMap{
		DefaultTarget: "codex",
		Targets: map[string]config.Target{
			"codex": {
				Project: []config.Mapping{
					{Source: "skills", Destination: ".codex/skills"},
					{Source: "sub-agents", Destination: ".codex/agents", Transform: true},
				},
			},
		},
	}
}

func writeFile(t *testing.T, path string, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}
