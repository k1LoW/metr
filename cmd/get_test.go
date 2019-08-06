package cmd

import (
	"bytes"
	"regexp"
	"strings"
	"testing"
)

func TestRunGet(t *testing.T) {
	tests := []struct {
		name         string
		args         []string
		wantExitCode int
		wantStdout   *regexp.Regexp
		wantStderr   string
	}{
		{"metr get cpu", []string{"cpu"}, 0, regexp.MustCompile(`^\d+\.\d+$`), ""},
		{"metr get foo", []string{"foo"}, 1, regexp.MustCompile(`^$`), "foo does not exist"},
		{"metr get all", []string{"all"}, 0, regexp.MustCompile(`user:\d+\.\d+\n`), ""},
	}

	for _, tt := range tests {
		stdout := new(bytes.Buffer)
		stderr := new(bytes.Buffer)
		exitCode := runGet(tt.args, 500, stdout, stderr)
		_ = stdout.String()
		if exitCode != tt.wantExitCode {
			t.Errorf("runGet(%v, 500, stdout, stderr) = %d, want = %d", tt.args, exitCode, tt.wantExitCode)
		}
		got := strings.TrimSuffix(stdout.String(), "\n")
		if !tt.wantStdout.MatchString(got) {
			t.Errorf("runGet(%v, 500, stdout, stderr) stdout = %s, want = %v", tt.args, got, tt.wantStdout)
		}
		got = strings.TrimSuffix(stderr.String(), "\n")
		if tt.wantStderr != got {
			t.Errorf("runGet(%v, 500, stdout, stderr) stderr = %s, want = %v", tt.args, got, tt.wantStderr)
		}
	}
}
