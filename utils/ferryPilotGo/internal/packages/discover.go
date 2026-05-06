package packages

import (
	"fmt"
	"io/fs"
	"sort"

	"ferrypilot/internal/assets"
)

type Package struct {
	Name string
}

func Discover(source assets.Source) ([]Package, error) {
	entries, err := fs.ReadDir(source.FS, source.Path())
	if err != nil {
		return nil, fmt.Errorf("read AISupport packages: %w", err)
	}
	result := make([]Package, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			result = append(result, Package{Name: entry.Name()})
		}
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Name < result[j].Name
	})
	return result, nil
}
