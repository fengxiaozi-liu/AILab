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

	"golang.org/x/term"
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
	if input, ok := terminalFile(options.Input); ok {
		if output, ok := terminalFile(options.Output); ok && term.IsTerminal(int(input.Fd())) && term.IsTerminal(int(output.Fd())) {
			return choosePackageInteractive(input, options.Output, available)
		}
	}
	return choosePackagePrompt(options.Input, options.Output, available)
}

func choosePackagePrompt(input io.Reader, output io.Writer, available []packages.Package) (string, error) {
	fmt.Fprintln(output, "Select an AISupport package:")
	for i, candidate := range available {
		fmt.Fprintf(output, "%d) %s\n", i+1, candidate.Name)
	}
	fmt.Fprint(output, "Package: ")
	var selected string
	if _, err := fmt.Fscan(input, &selected); err != nil {
		return "", fmt.Errorf("read package selection: %w", err)
	}
	for i, candidate := range available {
		if selected == candidate.Name || selected == fmt.Sprint(i+1) {
			return candidate.Name, nil
		}
	}
	return "", fmt.Errorf("unknown package selection %q", selected)
}

func choosePackageInteractive(input fileReader, output io.Writer, available []packages.Package) (string, error) {
	oldState, err := term.MakeRaw(int(input.Fd()))
	if err != nil {
		return "", fmt.Errorf("enable interactive package selection: %w", err)
	}
	defer term.Restore(int(input.Fd()), oldState)

	selected := 0
	renderedLines := 0
	render := func() {
		if renderedLines > 0 {
			fmt.Fprintf(output, "\x1b[%dA", renderedLines)
		}
		fmt.Fprintln(output, "Select an AISupport package:")
		for i, candidate := range available {
			prefix := "  "
			if i == selected {
				prefix = "> "
			}
			fmt.Fprintf(output, "\x1b[2K\r%s%s\n", prefix, candidate.Name)
		}
		renderedLines = len(available) + 1
	}

	render()
	buffer := make([]byte, 1)
	for {
		if _, err := input.Read(buffer); err != nil {
			return "", fmt.Errorf("read package selection: %w", err)
		}
		switch buffer[0] {
		case 3:
			return "", errors.New("package selection canceled")
		case '\r', '\n':
			fmt.Fprintf(output, "\r\n")
			return available[selected].Name, nil
		case 27:
			key, err := readEscapeKey(input)
			if err != nil {
				return "", err
			}
			switch key {
			case "up":
				if selected == 0 {
					selected = len(available) - 1
				} else {
					selected--
				}
				render()
			case "down":
				selected = (selected + 1) % len(available)
				render()
			}
		case 0, 224:
			key, err := readWindowsKey(input)
			if err != nil {
				return "", err
			}
			switch key {
			case "up":
				if selected == 0 {
					selected = len(available) - 1
				} else {
					selected--
				}
				render()
			case "down":
				selected = (selected + 1) % len(available)
				render()
			}
		}
	}
}

func readEscapeKey(input io.Reader) (string, error) {
	sequence := make([]byte, 2)
	if _, err := io.ReadFull(input, sequence); err != nil {
		return "", fmt.Errorf("read package selection: %w", err)
	}
	if sequence[0] != '[' {
		return "", nil
	}
	switch sequence[1] {
	case 'A':
		return "up", nil
	case 'B':
		return "down", nil
	default:
		return "", nil
	}
}

func readWindowsKey(input io.Reader) (string, error) {
	buffer := make([]byte, 1)
	if _, err := input.Read(buffer); err != nil {
		return "", fmt.Errorf("read package selection: %w", err)
	}
	switch buffer[0] {
	case 72:
		return "up", nil
	case 80:
		return "down", nil
	default:
		return "", nil
	}
}

type fileReader interface {
	io.Reader
	Fd() uintptr
}

func terminalFile(value any) (fileReader, bool) {
	file, ok := value.(fileReader)
	return file, ok
}

func packageNames(available []packages.Package) string {
	names := make([]string, 0, len(available))
	for _, candidate := range available {
		names = append(names, candidate.Name)
	}
	return strings.Join(names, ", ")
}
