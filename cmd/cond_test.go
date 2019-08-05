package cmd

import (
	"bytes"
	"strings"
	"testing"
)

func TestRunCond(t *testing.T) {
	tests := []struct {
		name         string
		args         []string
		wantExitCode int
		wantStderr   string
	}{
		{"metr cond 'cpu > 100'", []string{"cpu > 100"}, 1, ""},
		{"metr cond 'cpu > 0'", []string{"cpu > 0"}, 0, ""},
		{"metr cond 'foo > 10'", []string{"foo"}, 1, "undefined: foo"},
	}

	for _, tt := range tests {
		stdout := new(bytes.Buffer)
		stderr := new(bytes.Buffer)
		exitCode := runCond(tt.args, 500, stdout, stderr)
		_ = stdout.String()
		if exitCode != tt.wantExitCode {
			t.Errorf("runCond(%v, 500, stdout, stderr) = %d, want = %d", tt.args, exitCode, tt.wantExitCode)
		}
		got := strings.TrimSuffix(stderr.String(), "\n")
		if tt.wantStderr != got {
			t.Errorf("runCond(%v, 500, stdout, stderr) stderr = %s, want = %v", tt.args, got, tt.wantStderr)
		}
	}
}
