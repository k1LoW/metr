package cmd

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func TestRunCond(t *testing.T) {
	tests := []struct {
		name         string
		args         []string
		pid          int32
		wantExitCode int
		wantStderr   string
	}{
		{"metr cond 'cpu > 100'", []string{"cpu > 100"}, 0, 1, ""},
		{"metr cond 'cpu >= 0'", []string{"cpu >= 0"}, 0, 0, ""},
		{"metr cond 'cpu > 100 or mem < 100'", []string{"cpu > 100 or mem < 100"}, 0, 0, ""},
		{"metr cond 'foo > 10'", []string{"foo > 10"}, 0, 1, "undefined: foo"},
		{"metr cond 'proc_cpu >= 0'", []string{"proc_cpu >= 0"}, 0, 1, "undefined: proc_cpu"},
		{"metr cond 'proc_cpu >= 0' -p $PID", []string{"proc_cpu >= 0"}, int32(os.Getpid()), 0, ""},
	}

	for _, tt := range tests {
		stdout := new(bytes.Buffer)
		stderr := new(bytes.Buffer)
		exitCode := runCond(tt.args, 100, tt.pid, stdout, stderr)
		_ = stdout.String()
		if exitCode != tt.wantExitCode {
			t.Errorf("runCond(%v, 100, stdout, stderr) = %d, want = %d", tt.args, exitCode, tt.wantExitCode)
		}
		got := strings.TrimSuffix(stderr.String(), "\n")
		if tt.wantStderr != got {
			t.Errorf("runCond(%v, 100, stdout, stderr) stderr = %s, want = %v", tt.args, got, tt.wantStderr)
		}
	}
}
