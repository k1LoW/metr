package cmd

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func TestRunCheck(t *testing.T) {
	tests := []struct {
		desc         string
		warningCond  string
		criticalCond string
		pid          int32
		name         string
		wantExitCode int
		wantStdout   string
		wantStderr   string
	}{
		{"metr check -w 'cpu > 100' -c 'cpu > 100'", "cpu > 100", "cpu > 100", 0, "", 0, "METR OK: w(cpu > 100) c(cpu > 100)", ""},
		{"metr check -w 'cpu < 100' -c 'cpu > 100'", "cpu < 100", "cpu > 100", 0, "", 1, "METR WARNING: w(cpu < 100) c(cpu > 100)", ""},
		{"metr check -w 'cpu >= 0' -c 'cpu < 100'", "cpu >= 0", "cpu < 100", 0, "", 2, "METR CRITICAL: w(cpu >= 0) c(cpu < 100)", ""},
		{"metr check -w 'cpu >= 0'", "cpu >= 0", "", 0, "", 1, "METR WARNING: w(cpu >= 0) c()", ""},
		{"metr check -w 'proc_cpu >= 0'", "proc_cpu >= 0", "", 0, "", 3, "METR UNKNOWN: undefined: proc_cpu", ""},
		{"metr check -w 'proc_cpu >= 0' -p $PID", "proc_cpu >= 0", "", int32(os.Getpid()), "", 1, "METR WARNING: w(proc_cpu >= 0) c()", ""},
		{"metr check -w 'proc_cpu >= 0' -P [Name of target process]", "proc_cpu >= 0", "", 0, "go", 1, "METR WARNING: w(proc_cpu >= 0) c()", ""},
	}

	for _, tt := range tests {
		stdout := new(bytes.Buffer)
		stderr := new(bytes.Buffer)
		exitCode := runCheck([]string{}, tt.warningCond, tt.criticalCond, 100, tt.pid, tt.name, stdout, stderr)
		_ = stdout.String()
		if exitCode != tt.wantExitCode {
			t.Errorf("%s: runCheck([], %s, %s, 100, stdout, stderr) = %d, want = %d", tt.desc, tt.warningCond, tt.criticalCond, exitCode, tt.wantExitCode)
		}
		got := strings.TrimSuffix(stdout.String(), "\n")
		if tt.wantStdout != got {
			t.Errorf("%s: runCheck([], %s, %s, 100, stdout, stderr) stdout = %s, want = %v", tt.desc, tt.warningCond, tt.criticalCond, got, tt.wantStdout)
		}
		got = strings.TrimSuffix(stderr.String(), "\n")
		if tt.wantStderr != got {
			t.Errorf("%s: runCheck([], %s, %s, 100, stdout, stderr) stderr = %s, want = %v", tt.desc, tt.warningCond, tt.criticalCond, got, tt.wantStderr)
		}
	}
}
