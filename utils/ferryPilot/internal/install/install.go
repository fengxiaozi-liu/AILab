package install

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"ferrypilot/internal/assets"
	"ferrypilot/internal/transform"
)

type InstalledFile struct {
	Source      string
	Destination string
	Transformed bool
}

func Package(ctx context.Context, source assets.Source, packageName string, target TargetPlan) ([]InstalledFile, error) {
	var installed []InstalledFile
	for _, mapping := range target.Mappings {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		sourceDir := source.Path(packageName, mapping.SourceDir)
		items, err := installTree(source, sourceDir, mapping)
		if err != nil {
			return nil, err
		}
		installed = append(installed, items...)
	}
	return installed, nil
}

func installTree(source assets.Source, sourceDir string, mapping ResolvedMapping) ([]InstalledFile, error) {
	var installed []InstalledFile
	err := fs.WalkDir(source.FS, sourceDir, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			if errors.Is(walkErr, fs.ErrNotExist) {
				return nil
			}
			return walkErr
		}
		if entry.IsDir() {
			return nil
		}
		rel, err := filepath.Rel(filepath.FromSlash(sourceDir), filepath.FromSlash(path))
		if err != nil {
			return err
		}
		dest := filepath.Join(mapping.Destination, rel)
		transformed := false
		if mapping.Transform && strings.EqualFold(filepath.Ext(dest), ".md") {
			dest = strings.TrimSuffix(dest, filepath.Ext(dest)) + ".toml"
			transformed = true
		}
		if err := copyAsset(source, path, dest, transformed); err != nil {
			return err
		}
		installed = append(installed, InstalledFile{
			Source:      path,
			Destination: dest,
			Transformed: transformed,
		})
		return nil
	})
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return installed, nil
		}
		return nil, fmt.Errorf("install %s: %w", sourceDir, err)
	}
	return installed, nil
}

func copyAsset(source assets.Source, sourcePath string, destination string, transformed bool) error {
	content, err := fs.ReadFile(source.FS, sourcePath)
	if err != nil {
		return err
	}
	if transformed {
		content, err = transform.MarkdownSubAgentToTOML(content)
		if err != nil {
			return err
		}
	}
	if err := os.MkdirAll(filepath.Dir(destination), 0o755); err != nil {
		return err
	}
	return os.WriteFile(destination, content, 0o644)
}
