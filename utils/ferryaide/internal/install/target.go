package install

import (
	"fmt"
	"path/filepath"

	"ferryaide/internal/config"
)

type TargetPlan struct {
	Mode     string
	Agent    string
	Root     string
	Mappings []ResolvedMapping
}

type ResolvedMapping struct {
	SourceDir   string
	Destination string
	Transform   bool
}

type TargetOptions struct {
	FileMap config.FileMap
	Mode    string
	Target  string
	HomeDir string
	WorkDir string
}

func ResolveTarget(options TargetOptions) (TargetPlan, error) {
	agent, target, err := config.LookupTarget(options.FileMap, options.Target)
	if err != nil {
		return TargetPlan{}, err
	}
	root, err := installRoot(options)
	if err != nil {
		return TargetPlan{}, err
	}
	mappings, err := config.MappingsFor(target, options.Mode)
	if err != nil {
		return TargetPlan{}, err
	}

	resolved := make([]ResolvedMapping, 0, len(mappings))
	for _, mapping := range mappings {
		resolved = append(resolved, ResolvedMapping{
			SourceDir:   filepath.FromSlash(mapping.Source),
			Destination: filepath.Join(root, filepath.FromSlash(mapping.Destination)),
			Transform:   mapping.Transform,
		})
	}

	return TargetPlan{
		Mode:     options.Mode,
		Agent:    agent,
		Root:     root,
		Mappings: resolved,
	}, nil
}

func installRoot(options TargetOptions) (string, error) {
	switch options.Mode {
	case "global":
		if options.HomeDir == "" {
			return "", fmt.Errorf("home directory is required for global install")
		}
		return options.HomeDir, nil
	case "project":
		if options.WorkDir == "" {
			return "", fmt.Errorf("working directory is required for project install")
		}
		return options.WorkDir, nil
	default:
		return "", fmt.Errorf("invalid install mode %q", options.Mode)
	}
}
