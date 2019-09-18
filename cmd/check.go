/*
Copyright © 2019 Ken'ichiro Oyama <k1lowxb@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/antonmedv/expr"
	"github.com/k1LoW/metr/logger"
	"github.com/k1LoW/metr/metrics"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Status int

const (
	OK Status = iota
	WARNING
	CRITICAL
	UNKNOWN
)

func (s Status) String() string {
	switch {
	case s == OK:
		return "OK"
	case s == WARNING:
		return "WARNING"
	case s == CRITICAL:
		return "CRITICAL"
	default:
		return "UNKNOWN"
	}
}

type Result struct {
	stdout io.Writer
	stderr io.Writer
	logger *zap.Logger
}

func (r *Result) exitWithStdout(status Status, warningCond string, criticalCond string, ms map[string]interface{}, err error) int {
	// SERVICE STATUS: First line of output | First part of performance data http://nagios-plugins.org/doc/guidelines.html#PLUGOUTPUT
	if err != nil {
		_, _ = fmt.Fprintf(r.stdout, "%s %s: %s\n", "METR", status, err)
	} else {
		_, _ = fmt.Fprintf(r.stdout, "%s %s: w(%s) c(%s)\n", "METR", status, warningCond, criticalCond)
	}
	if r.logger != nil {
		fields := []zapcore.Field{
			zap.String("command", "check"),
			zap.String("warning_condition", warningCond),
			zap.String("critical_condition", criticalCond),
		}
		for k, v := range ms {
			fields = append(fields, zap.Any(k, v))
		}
		fields = append(fields, zap.String("result", status.String()))
		r.logger.Info("execute metr check", fields...)
	}
	return int(status)
}

var (
	warningCond  string
	criticalCond string
)

// checkCmd represents the check command
var checkCmd = &cobra.Command{
	Use:   "check -w [WARNING_CONDITION] -c [CRITICAL_CONDITION]",
	Short: "check metrics condition and output result with exit status code",
	Long:  `check metrics condition and output result with exit status code.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if pid > 0 && name != "" {
			return errors.WithStack(errors.New("target option can only be either --pid or --name"))
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		os.Exit(runCheck(args, warningCond, criticalCond, interval, pid, name, dir, os.Stdout, os.Stderr))
	},
}

func runCheck(args []string, warningCond, criticalCond string, interval int, pid int32, name, dir string, stdout, stderr io.Writer) (exitCode int) {
	var (
		m   *metrics.Metrics
		err error
		lgr *zap.Logger
	)
	ms := make(map[string]interface{})
	if dir != "" {
		lgr, err = logger.NewLogger(dir)
	} else {
		lgr, err = logger.NewNoLogger()
	}
	r := &Result{
		stdout: stdout,
		stderr: stderr,
		logger: lgr,
	}
	if err != nil {
		return r.exitWithStdout(UNKNOWN, warningCond, criticalCond, ms, err)
	}
	if len(args) > 0 {
		return r.exitWithStdout(UNKNOWN, warningCond, criticalCond, ms, errors.New("metr requires no args"))
	}
	if warningCond == "" && criticalCond == "" {
		return r.exitWithStdout(UNKNOWN, warningCond, criticalCond, ms, errors.New("metr requires -w or -c option"))
	}
	switch {
	case name != "":
		m, err = metrics.GetMetricsByName(time.Duration(interval)*time.Millisecond, name)
		if err != nil {
			return r.exitWithStdout(UNKNOWN, warningCond, criticalCond, ms, err)
		}
	default:
		m, err = metrics.GetMetrics(time.Duration(interval)*time.Millisecond, pid)
		if err != nil {
			return r.exitWithStdout(UNKNOWN, warningCond, criticalCond, ms, err)
		}
	}
	if criticalCond != "" {
		ms = m.Raw()
		got, err := expr.Eval(fmt.Sprintf("(%s) == true", criticalCond), ms)
		if err != nil {
			return r.exitWithStdout(UNKNOWN, warningCond, criticalCond, ms, err)
		}
		if got.(bool) {
			return r.exitWithStdout(CRITICAL, warningCond, criticalCond, ms, nil)
		}
	}
	if warningCond != "" {
		ms = m.Raw()
		got, err := expr.Eval(fmt.Sprintf("(%s) == true", warningCond), ms)
		if err != nil {
			return r.exitWithStdout(UNKNOWN, warningCond, criticalCond, ms, err)
		}
		if got.(bool) {
			return r.exitWithStdout(WARNING, warningCond, criticalCond, ms, nil)
		}
	}

	return r.exitWithStdout(OK, warningCond, criticalCond, ms, nil)
}

func init() {
	rootCmd.AddCommand(checkCmd)
	checkCmd.Flags().StringVarP(&warningCond, "warning", "w", "", "WARNING condition")
	checkCmd.Flags().StringVarP(&criticalCond, "critical", "c", "", "CRITICAL condition")
	checkCmd.Flags().IntVarP(&interval, "interval", "i", 500, "metric measurement interval (millisecond)")
	checkCmd.Flags().Int32VarP(&pid, "pid", "p", 0, "PID of target process")
	checkCmd.Flags().StringVarP(&name, "name", "P", "", "Name of target process")
	checkCmd.Flags().StringVarP(&dir, "log-dir", "", "", "log directory")
}
