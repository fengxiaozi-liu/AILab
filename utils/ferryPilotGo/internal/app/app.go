package app

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"

	"ferrypilot/internal/assets"
	"ferrypilot/internal/config"
	"ferrypilot/internal/install"
	"ferrypilot/internal/packages"
)

const ToolName = "ferryPilot"

var ErrNotImplemented = errors.New("ferryPilot Go implementation is not complete yet")

type InstallMode string

const (
	InstallModeUnset   InstallMode = ""
	InstallModeGlobal  InstallMode = "global"
	InstallModeProject InstallMode = "project"
)

type Options struct {
	Args       []string
	Mode       InstallMode
	Target     string
	WorkDir    string
	HomeDir    string
	Package    string
	AssetRoot  string
	ConfigPath string
	Input      io.Reader
	Output     io.Writer
}

type Result struct {
	Message   string
	Package   string
	Installed []InstalledFile
}

type InstalledFile struct {
	Source      string
	Destination string
	Transformed bool
}

func Run(ctx context.Context, options Options) (Result, error) {
	select {
	case <-ctx.Done():
		return Result{}, ctx.Err()
	default:
	}
	fileMap, err := config.Load(options.ConfigPath, options.WorkDir)
	if err != nil {
		return Result{}, err
	}
	source, err := assets.NewSource(ctx, fileMap.DataSource, options.AssetRoot)
	if err != nil {
		return Result{}, err
	}
	defer source.Close()
	available, err := packages.Discover(source)
	if err != nil {
		return Result{}, err
	}
	packageName, err := choosePackage(options, available)
	if err != nil {
		return Result{}, err
	}
	target, err := install.ResolveTarget(install.TargetOptions{
		FileMap: fileMap,
		Mode:    string(options.Mode),
		Target:  options.Target,
		HomeDir: options.HomeDir,
		WorkDir: options.WorkDir,
	})
	if err != nil {
		return Result{}, err
	}
	installed, err := install.Package(ctx, source, packageName, target)
	if err != nil {
		return Result{}, err
	}

	files := make([]InstalledFile, 0, len(installed))
	for _, item := range installed {
		files = append(files, InstalledFile{
			Source:      item.Source,
			Destination: item.Destination,
			Transformed: item.Transformed,
		})
	}

	return Result{
		Message:   fmt.Sprintf("%s installed package %q for %s (%d files)", ToolName, packageName, target.Agent, len(files)),
		Package:   packageName,
		Installed: files,
	}, nil
}

func choosePackage(options Options, available []packages.Package) (string, error) {
	if len(available) == 0 {
		return "", errors.New("no AISupport packages found")
	}
	if options.Package != "" {
		for _, candidate := range available {
			if candidate.Name == options.Package {
				return options.Package, nil
			}
		}
		return "", fmt.Errorf("unknown package %q", options.Package)
	}
	if len(available) == 1 {
		return available[0].Name, nil
	}
	if options.Input == nil || options.Output == nil {
		return "", fmt.Errorf("package name is required; available packages: %s", packageNames(available))
	}
	fmt.Fprintln(options.Output, "Select an AISupport package:")
	for i, candidate := range available {
		fmt.Fprintf(options.Output, "%d) %s\n", i+1, candidate.Name)
	}
	fmt.Fprint(options.Output, "Package: ")
	var selected string
	if _, err := fmt.Fscan(options.Input, &selected); err != nil {
		return "", fmt.Errorf("read package selection: %w", err)
	}
	for i, candidate := range available {
		if selected == candidate.Name || selected == fmt.Sprint(i+1) {
			return candidate.Name, nil
		}
	}
	return "", fmt.Errorf("unknown package selection %q", selected)
}

func packageNames(available []packages.Package) string {
	names := make([]string, 0, len(available))
	for _, candidate := range available {
		names = append(names, candidate.Name)
	}
	return strings.Join(names, ", ")
}
