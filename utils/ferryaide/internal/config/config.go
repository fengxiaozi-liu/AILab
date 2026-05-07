package config

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

//go:embed file_map.json
var defaultFileMap []byte

type FileMap struct {
	DataSource    DataSource        `json:"data_source"`
	DefaultTarget string            `json:"default_target"`
	Targets       map[string]Target `json:"targets"`
}

type DataSource struct {
	Type       string `json:"type"`
	Repository string `json:"repository"`
	Ref        string `json:"ref"`
	Path       string `json:"path"`
}

type Mapping struct {
	Source      string `json:"source"`
	Destination string `json:"destination"`
	Transform   bool   `json:"transform"`
}

type Target struct {
	Global  []Mapping `json:"global"`
	Project []Mapping `json:"project"`
}

func Load(path string, workDir string) (FileMap, error) {
	if path == "" {
		return parse(defaultFileMap, "embedded default file_map.json")
	}
	content, err := os.ReadFile(path)
	if err != nil {
		return FileMap{}, fmt.Errorf("read file map %s: %w", path, err)
	}
	return parse(content, path)
}

func parse(content []byte, source string) (FileMap, error) {
	content = bytes.TrimPrefix(content, []byte{0xEF, 0xBB, 0xBF})
	var fileMap FileMap
	if err := json.Unmarshal(content, &fileMap); err != nil {
		return FileMap{}, fmt.Errorf("parse file map %s: %w", source, err)
	}
	if fileMap.DefaultTarget == "" {
		fileMap.DefaultTarget = "codex"
	}
	if len(fileMap.Targets) == 0 {
		return FileMap{}, fmt.Errorf("file map has no targets")
	}
	if fileMap.DataSource.Type == "" {
		return FileMap{}, fmt.Errorf("file map data_source.type is required")
	}
	return fileMap, nil
}

func FindDefault(workDir string) (string, error) {
	candidates := []string{}
	if workDir != "" {
		candidates = appendDefaultCandidates(candidates, workDir)
	}
	if exe, err := os.Executable(); err == nil {
		candidates = appendDefaultCandidates(candidates, filepath.Dir(exe))
	}
	for _, candidate := range candidates {
		if info, err := os.Stat(candidate); err == nil && !info.IsDir() {
			return candidate, nil
		}
	}
	return "", fmt.Errorf("file_map.json not found")
}

func appendDefaultCandidates(candidates []string, start string) []string {
	current, err := filepath.Abs(start)
	if err != nil {
		current = start
	}
	for {
		candidates = append(candidates,
			filepath.Join(current, "config", "file_map.json"),
			filepath.Join(current, "file_map.json"),
		)
		parent := filepath.Dir(current)
		if parent == current {
			return candidates
		}
		current = parent
	}
}

func LookupTarget(fileMap FileMap, name string) (string, Target, error) {
	if name == "" {
		name = fileMap.DefaultTarget
	}
	target, ok := fileMap.Targets[name]
	if !ok {
		return "", Target{}, fmt.Errorf("unknown target %q", name)
	}
	return name, target, nil
}

func MappingsFor(target Target, mode string) ([]Mapping, error) {
	switch mode {
	case "global":
		return target.Global, nil
	case "project":
		return target.Project, nil
	default:
		return nil, fmt.Errorf("invalid install mode %q", mode)
	}
}
