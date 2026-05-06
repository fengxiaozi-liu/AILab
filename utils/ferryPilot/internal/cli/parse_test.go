package cli

import (
	"testing"

	"ferrypilot/internal/app"
)

func TestParseInstallMode(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want app.InstallMode
	}{
		{name: "global short", args: []string{"-g"}, want: app.InstallModeGlobal},
		{name: "global long", args: []string{"--global"}, want: app.InstallModeGlobal},
		{name: "project short", args: []string{"-p"}, want: app.InstallModeProject},
		{name: "project long", args: []string{"--project"}, want: app.InstallModeProject},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.args)
			if err != nil {
				t.Fatalf("Parse returned error: %v", err)
			}
			if got.Mode != tt.want {
				t.Fatalf("Mode = %q, want %q", got.Mode, tt.want)
			}
		})
	}
}

func TestParseRejectsInvalidModes(t *testing.T) {
	tests := [][]string{
		{},
		{"-g", "-p"},
		{"--global", "--project"},
	}
	for _, args := range tests {
		if _, err := Parse(args); err == nil {
			t.Fatalf("Parse(%v) returned nil error", args)
		}
	}
}

func TestParseTarget(t *testing.T) {
	got, err := Parse([]string{"-g", "-t", "codex"})
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}
	if got.Target != "codex" {
		t.Fatalf("Target = %q, want codex", got.Target)
	}
}
