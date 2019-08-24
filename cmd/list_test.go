package cmd

import (
	"bytes"
	"os"
	"regexp"
	"strings"
	"testing"
)

func TestRunList(t *testing.T) {
	tests := []struct {
		desc       string
		pid        int32
		name       string
		wantStdout *regexp.Regexp
		wantStderr string
	}{
		{"metr list", 0, "", regexp.MustCompile(`user \(now:\d+\.\d+ %\)`), ""},
		{"metr list -p $PID", int32(os.Getpid()), "", regexp.MustCompile(`proc_cpu \(now:\d+\.\d+ %\)`), ""},
		{"metr list -P [Name of target process]", 0, "go", regexp.MustCompile(`proc_cpu \(now:\d+\.\d+ %\)`), ""},
	}

	for _, tt := range tests {
		stdout := new(bytes.Buffer)
		stderr := new(bytes.Buffer)
		exitCode := runList([]string{}, 100, tt.pid, tt.name, stdout, stderr)
		_ = stdout.String()
		if exitCode != 0 {
			t.Errorf("%s: runList([], 100, stdout, stderr) = %d, want = %d", tt.desc, exitCode, 0)
		}
		got := strings.TrimSuffix(stdout.String(), "\n")
		if !tt.wantStdout.MatchString(got) {
			t.Errorf("%s: runList([], 100, stdout, stderr) stdout = %s, want = %v", tt.desc, got, tt.wantStdout)
		}
		got = strings.TrimSuffix(stderr.String(), "\n")
		if tt.wantStderr != got {
			t.Errorf("%s: runList([], 100, stdout, stderr) stderr = %s, want = %v", tt.desc, got, tt.wantStderr)
		}
	}
}
