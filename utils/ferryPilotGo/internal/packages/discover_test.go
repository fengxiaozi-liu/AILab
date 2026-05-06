package packages

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"ferrypilot/internal/assets"
	"ferrypilot/internal/config"
)

func TestDiscoverOnlyFirstLevelDirectories(t *testing.T) {
	root := t.TempDir()
	for _, dir := range []string{
		filepath.Join(root, "AISupport", "speckit", "skills"),
		filepath.Join(root, "AISupport", "kratos"),
	} {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.WriteFile(filepath.Join(root, "AISupport", "README.md"), []byte("ignore"), 0o644); err != nil {
		t.Fatal(err)
	}
	source, err := assets.NewSource(context.Background(), config.DataSource{Type: "local", Path: root}, "")
	if err != nil {
		t.Fatal(err)
	}
	got, err := Discover(source)
	if err != nil {
		t.Fatalf("Discover returned error: %v", err)
	}
	if len(got) != 2 || got[0].Name != "kratos" || got[1].Name != "speckit" {
		t.Fatalf("packages = %#v, want kratos and speckit", got)
	}
}
