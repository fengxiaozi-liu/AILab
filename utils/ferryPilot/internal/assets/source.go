package assets

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"

	"ferrypilot/internal/config"
)

type Source struct {
	FS      fs.FS
	Root    string
	Name    string
	cleanup func() error
}

func NewSource(ctx context.Context, dataSource config.DataSource, explicitRoot string) (Source, error) {
	if explicitRoot != "" {
		return localSource(explicitRoot, "local")
	}
	switch dataSource.Type {
	case "git":
		return gitSource(ctx, dataSource)
	case "local":
		return localSource(dataSource.Path, "local")
	default:
		return Source{}, fmt.Errorf("unsupported data source type %q", dataSource.Type)
	}
}

func (s Source) Close() error {
	if s.cleanup == nil {
		return nil
	}
	return s.cleanup()
}

func (s Source) Path(parts ...string) string {
	items := make([]string, 0, len(parts)+1)
	if s.Root != "." && s.Root != "" {
		items = append(items, s.Root)
	}
	items = append(items, parts...)
	if len(items) == 0 {
		return "."
	}
	return filepath.ToSlash(filepath.Join(items...))
}

func localSource(root string, name string) (Source, error) {
	if root == "" {
		return Source{}, fmt.Errorf("local data source path is required")
	}
	if ok, err := hasAISupport(root); err != nil {
		return Source{}, err
	} else if !ok {
		return Source{}, fmt.Errorf("AISupport directory not found under %s", root)
	}
	return Source{FS: os.DirFS(root), Root: "AISupport", Name: name}, nil
}

func gitSource(ctx context.Context, dataSource config.DataSource) (Source, error) {
	if dataSource.Repository == "" {
		return Source{}, fmt.Errorf("git data source repository is required")
	}
	tmp, err := os.MkdirTemp("", "ferryPilot-*")
	if err != nil {
		return Source{}, fmt.Errorf("create temp directory: %w", err)
	}
	args := []string{"clone", "--depth", "1"}
	if dataSource.Ref != "" {
		args = append(args, "--branch", dataSource.Ref)
	}
	args = append(args, dataSource.Repository, tmp)
	cmd := exec.CommandContext(ctx, "git", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		os.RemoveAll(tmp)
		return Source{}, fmt.Errorf("clone data source: %w: %s", err, string(output))
	}
	source, err := localSource(tmp, "git")
	if err != nil {
		os.RemoveAll(tmp)
		return Source{}, err
	}
	source.cleanup = func() error {
		return os.RemoveAll(tmp)
	}
	return source, nil
}

func hasAISupport(root string) (bool, error) {
	info, err := os.Stat(filepath.Join(root, "AISupport"))
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return info.IsDir(), nil
}
