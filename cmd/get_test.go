package cmd

import (
	"bytes"
	"os"
	"regexp"
	"strings"
	"testing"
)

func TestRunGet(t *testing.T) {
	tests := []struct {
		name         string
		args         []string
		pid          int32
		wantExitCode int
		wantStdout   *regexp.Regexp
		wantStderr   string
	}{
		{"metr get cpu", []string{"cpu"}, 0, 0, regexp.MustCompile(`^\d+\.\d+$`), ""},
		{"metr get foo", []string{"foo"}, 0, 1, regexp.MustCompile(`^$`), "foo does not exist"},
		{"metr get all", []string{"all"}, 0, 0, regexp.MustCompile(`user:\d+\.\d+\n`), ""},
		{"metr get proc_cpu", []string{"proc_cpu"}, 0, 1, regexp.MustCompile(`^$`), "proc_cpu does not exist"},
		{"metr get proc_cpu -p 1", []string{"proc_cpu"}, int32(os.Getpid()), 0, regexp.MustCompile(`^\d+\.\d+$`), ""},
		{"metr get all -p 1`", []string{"all"}, int32(os.Getpid()), 0, regexp.MustCompile(`proc_cpu:\d+\.\d+\n`), ""},
	}

	for _, tt := range tests {
		stdout := new(bytes.Buffer)
		stderr := new(bytes.Buffer)
		exitCode := runGet(tt.args, 100, tt.pid, stdout, stderr)
		_ = stdout.String()
		if exitCode != tt.wantExitCode {
			t.Errorf("runGet(%v, 100, stdout, stderr) = %d, want = %d", tt.args, exitCode, tt.wantExitCode)
		}
		got := strings.TrimSuffix(stdout.String(), "\n")
		if !tt.wantStdout.MatchString(got) {
			t.Errorf("runGet(%v, 100, stdout, stderr) stdout = %s, want = %v", tt.args, got, tt.wantStdout)
		}
		got = strings.TrimSuffix(stderr.String(), "\n")
		if tt.wantStderr != got {
			t.Errorf("runGet(%v, 100, stdout, stderr) stderr = %s, want = %v", tt.args, got, tt.wantStderr)
		}
	}
}
