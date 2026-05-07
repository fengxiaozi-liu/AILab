package cli

import (
	"flag"
	"fmt"
	"os"

	"ferryaide/internal/app"
)

func Parse(args []string) (app.Options, error) {
	flags := flag.NewFlagSet(app.ToolName, flag.ContinueOnError)
	flags.SetOutput(discard{})

	var globalShort bool
	var globalLong bool
	var projectShort bool
	var projectLong bool
	var targetShort string
	var targetLong string
	var configPath string

	flags.BoolVar(&globalShort, "g", false, "install globally")
	flags.BoolVar(&globalLong, "global", false, "install globally")
	flags.BoolVar(&projectShort, "p", false, "install into the current project")
	flags.BoolVar(&projectLong, "project", false, "install into the current project")
	flags.StringVar(&targetShort, "t", "", "target agent")
	flags.StringVar(&targetLong, "target", "", "target agent")
	flags.StringVar(&configPath, "config", "", "path to file_map.json")

	if err := flags.Parse(args); err != nil {
		return app.Options{}, err
	}
	remaining := flags.Args()
	if len(remaining) > 1 {
		return app.Options{}, fmt.Errorf("expected at most one package name, got %d", len(remaining))
	}

	global := globalShort || globalLong
	project := projectShort || projectLong
	if global == project {
		return app.Options{}, fmt.Errorf("specify exactly one of -g/--global or -p/--project")
	}

	target := targetLong
	if target == "" {
		target = targetShort
	}

	mode := app.InstallModeGlobal
	if project {
		mode = app.InstallModeProject
	}

	workDir, err := os.Getwd()
	if err != nil {
		return app.Options{}, fmt.Errorf("resolve working directory: %w", err)
	}
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return app.Options{}, fmt.Errorf("resolve home directory: %w", err)
	}

	options := app.Options{
		Args:       args,
		Mode:       mode,
		Target:     target,
		WorkDir:    workDir,
		HomeDir:    homeDir,
		ConfigPath: configPath,
	}
	if len(remaining) == 1 {
		options.Package = remaining[0]
	}
	return options, nil
}

type discard struct{}

func (discard) Write(p []byte) (int, error) {
	return len(p), nil
}
