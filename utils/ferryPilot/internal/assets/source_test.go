package assets

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"ferrypilot/internal/config"
)

func TestNewSourceUsesLocalDataSource(t *testing.T) {
	root := t.TempDir()
	if err := os.Mkdir(filepath.Join(root, "AISupport"), 0o755); err != nil {
		t.Fatal(err)
	}
	source, err := NewSource(context.Background(), config.DataSource{Type: "local", Path: root}, "")
	if err != nil {
		t.Fatalf("NewSource returned error: %v", err)
	}
	if source.Name != "local" {
		t.Fatalf("Name = %q, want local", source.Name)
	}
	if got := source.Path("speckit"); got != "AISupport/speckit" {
		t.Fatalf("Path = %q, want AISupport/speckit", got)
	}
}
