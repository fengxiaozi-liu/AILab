package app

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunInstallsSelectedPackage(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "AISupport", "speckit", "skills", "demo", "SKILL.md"), "# demo")
	project := filepath.Join(root, "project")

	result, err := Run(context.Background(), Options{
		Mode:       InstallModeProject,
		Target:     "codex",
		WorkDir:    project,
		HomeDir:    filepath.Join(root, "home"),
		AssetRoot:  root,
		ConfigPath: writeConfig(t, root),
		Package:    "speckit",
	})
	if err != nil {
		t.Fatalf("Run returned error: %v", err)
	}
	if result.Package != "speckit" {
		t.Fatalf("Package = %q, want speckit", result.Package)
	}
	if !strings.Contains(result.Message, ToolName) {
		t.Fatalf("Message = %q, want tool name", result.Message)
	}
	if _, err := os.Stat(filepath.Join(project, ".codex", "skills", "demo", "SKILL.md")); err != nil {
		t.Fatalf("installed file missing: %v", err)
	}
}

func writeConfig(t *testing.T, root string) string {
	t.Helper()
	path := filepath.Join(root, "file_map.json")
	content := `{
  "data_source": {"type": "local", "path": "` + filepath.ToSlash(root) + `"},
  "default_target": "codex",
  "targets": {
    "codex": {
      "project": [
        {"source": "skills", "destination": ".codex/skills", "transform": false},
        {"source": "sub-agents", "destination": ".codex/agents", "transform": true}
      ]
    }
  }
}`
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	return path
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
