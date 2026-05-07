package app

import (
	"context"
	"errors"
	"fmt"
	"io"
	"sort"
	"strings"

	"ferryaide/internal/assets"
	"ferryaide/internal/config"
	"ferryaide/internal/install"
	"ferryaide/internal/packages"

	"golang.org/x/term"
)

const ToolName = "ferryaide"

var ErrNotImplemented = errors.New("ferryaide Go implementation is not complete yet")

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
	showCursorAfterRun := hideTerminalCursor(options.Output)
	defer showCursorAfterRun()

	fileMap, err := config.Load(options.ConfigPath, options.WorkDir)
	if err != nil {
		return Result{}, err
	}
	targetName, err := chooseTarget(options, fileMap)
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
		Target:  targetName,
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

func chooseTarget(options Options, fileMap config.FileMap) (string, error) {
	if options.Target != "" {
		name, _, err := config.LookupTarget(fileMap, options.Target)
		return name, err
	}
	available := targetNames(fileMap)
	if len(available) == 0 {
		return "", errors.New("no target agents found")
	}
	if len(available) == 1 {
		return available[0], nil
	}
	if options.Input == nil || options.Output == nil {
		return "", fmt.Errorf("target agent is required; available targets: %s", strings.Join(available, ", "))
	}
	return chooseNamedItem(options.Input, options.Output, "target agent", "Target", available, fileMap.DefaultTarget)
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
	return chooseNamedItem(options.Input, options.Output, "AISupport package", "Package", packageNamesList(available), "")
}

func chooseNamedItem(input io.Reader, output io.Writer, itemName string, prompt string, available []string, defaultName string) (string, error) {
	if inputFile, ok := terminalFile(input); ok {
		if outputFile, ok := terminalFile(output); ok && term.IsTerminal(int(inputFile.Fd())) && term.IsTerminal(int(outputFile.Fd())) {
			return chooseInteractive(inputFile, output, itemName, available, defaultName)
		}
	}
	return choosePrompt(input, output, itemName, prompt, available)
}

func choosePrompt(input io.Reader, output io.Writer, itemName string, prompt string, available []string) (string, error) {
	fmt.Fprintf(output, "Select a %s:\n", itemName)
	for i, candidate := range available {
		fmt.Fprintf(output, "%d) %s\n", i+1, candidate)
	}
	fmt.Fprintf(output, "%s: ", prompt)
	var selected string
	if _, err := fmt.Fscan(input, &selected); err != nil {
		return "", fmt.Errorf("read %s selection: %w", itemName, err)
	}
	for i, candidate := range available {
		if selected == candidate || selected == fmt.Sprint(i+1) {
			return candidate, nil
		}
	}
	return "", fmt.Errorf("unknown %s selection %q", itemName, selected)
}

func chooseInteractive(input fileReader, output io.Writer, itemName string, available []string, defaultName string) (string, error) {
	oldState, err := term.MakeRaw(int(input.Fd()))
	if err != nil {
		return "", fmt.Errorf("enable interactive %s selection: %w", itemName, err)
	}
	defer term.Restore(int(input.Fd()), oldState)

	selected := 0
	if defaultName != "" {
		for i, candidate := range available {
			if candidate == defaultName {
				selected = i
				break
			}
		}
	}
	renderedLines := 0
	render := func() {
		if renderedLines > 0 {
			fmt.Fprintf(output, "\x1b[%dA", renderedLines)
		}
		fmt.Fprintf(output, "Select a %s:\n", itemName)
		for i, candidate := range available {
			prefix := "  "
			if i == selected {
				prefix = "> "
			}
			fmt.Fprintf(output, "\x1b[2K\r%s%s\n", prefix, candidate)
		}
		renderedLines = len(available) + 1
	}

	render()
	buffer := make([]byte, 1)
	for {
		if _, err := input.Read(buffer); err != nil {
			return "", fmt.Errorf("read %s selection: %w", itemName, err)
		}
		switch buffer[0] {
		case 3:
			return "", fmt.Errorf("%s selection canceled", itemName)
		case '\r', '\n':
			return available[selected], nil
		case 27:
			key, err := readEscapeKey(input, itemName)
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
			key, err := readWindowsKey(input, itemName)
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

func hideTerminalCursor(output io.Writer) func() {
	if outputFile, ok := terminalFile(output); ok && term.IsTerminal(int(outputFile.Fd())) {
		fmt.Fprint(output, "\x1b[?25l")
		return func() {
			fmt.Fprint(output, "\x1b[?25h")
		}
	}
	return func() {}
}

func readEscapeKey(input io.Reader, itemName string) (string, error) {
	sequence := make([]byte, 2)
	if _, err := io.ReadFull(input, sequence); err != nil {
		return "", fmt.Errorf("read %s selection: %w", itemName, err)
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

func readWindowsKey(input io.Reader, itemName string) (string, error) {
	buffer := make([]byte, 1)
	if _, err := input.Read(buffer); err != nil {
		return "", fmt.Errorf("read %s selection: %w", itemName, err)
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
	return strings.Join(packageNamesList(available), ", ")
}

func packageNamesList(available []packages.Package) []string {
	names := make([]string, 0, len(available))
	for _, candidate := range available {
		names = append(names, candidate.Name)
	}
	return names
}

func targetNames(fileMap config.FileMap) []string {
	names := make([]string, 0, len(fileMap.Targets))
	for name := range fileMap.Targets {
		if name != fileMap.DefaultTarget {
			names = append(names, name)
		}
	}
	sort.Strings(names)
	if fileMap.DefaultTarget != "" {
		if _, ok := fileMap.Targets[fileMap.DefaultTarget]; ok {
			names = append([]string{fileMap.DefaultTarget}, names...)
		}
	}
	return names
}
