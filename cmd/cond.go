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

// condCmd represents the cond command
var condCmd = &cobra.Command{
	Use:   "cond [CONDITION]",
	Short: "returns CONDITION result using exit code",
	Long:  `returns CONDITION result using exit code.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.WithStack(errors.New("metr requires one arg"))
		}
		if pid > 0 && name != "" {
			return errors.WithStack(errors.New("target option can only be either --pid or --name"))
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		os.Exit(runCond(args, interval, pid, name, dir, os.Stdout, os.Stderr))
	},
}

var testCmd = &cobra.Command{
	Use:   "test [CONDITION]",
	Short: "alias for cond",
	Long:  `alias for cond.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.WithStack(errors.New("metr requires one arg"))
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		os.Exit(runCond(args, interval, pid, name, dir, os.Stdout, os.Stderr))
	},
}

func runCond(args []string, interval int, pid int32, name, dir string, stdout, stderr io.Writer) (exitCode int) {
	mcond := args[0]
	var (
		m   *metrics.Metrics
		err error
		lgr *zap.Logger
	)
	if dir != "" {
		lgr, err = logger.NewLogger(dir)
	} else {
		lgr, err = logger.NewNoLogger()
	}
	if err != nil {
		_, _ = fmt.Fprintf(stderr, "%s\n", err)
		return 1
	}

	switch {
	case name != "":
		m, err = metrics.GetMetricsByName(time.Duration(interval)*time.Millisecond, name)
		if err != nil {
			_, _ = fmt.Fprintf(stderr, "%s\n", err)
			return 1
		}
	default:
		m, err = metrics.GetMetrics(time.Duration(interval)*time.Millisecond, pid)
		if err != nil {
			_, _ = fmt.Fprintf(stderr, "%s\n", err)
			return 1
		}
	}
	ms := m.Raw()
	got, err := expr.Eval(fmt.Sprintf("(%s) == true", mcond), ms)
	if err != nil {
		_, _ = fmt.Fprintf(stderr, "%s\n", err)
		return 1
	}
	fields := []zapcore.Field{
		zap.String("command", "cond"),
		zap.String("condition", mcond),
	}
	for k, v := range ms {
		fields = append(fields, zap.Any(k, v))
	}
	if got.(bool) {
		fields = append(fields, zap.Int("result", 0))
		lgr.Info("execute metr cond", fields...)
		return 0
	} else {
		fields = append(fields, zap.Int("result", 1))
		lgr.Info("execute metr cond", fields...)
		return 1
	}
}

func init() {
	condCmd.Flags().IntVarP(&interval, "interval", "i", 500, "metric measurement interval (millisecond)")
	condCmd.Flags().Int32VarP(&pid, "pid", "p", 0, "PID of target process")
	condCmd.Flags().StringVarP(&name, "name", "P", "", "Name of target process")
	condCmd.Flags().StringVarP(&dir, "log-dir", "", "", "log directory")
	testCmd.Flags().IntVarP(&interval, "interval", "i", 500, "metric measurement interval (millisecond)")
	testCmd.Flags().Int32VarP(&pid, "pid", "p", 0, "PID of target process")
	testCmd.Flags().StringVarP(&name, "name", "P", "", "Name of target process")
	testCmd.Flags().StringVarP(&dir, "log-dir", "", "", "log directory")
	rootCmd.AddCommand(condCmd)
	rootCmd.AddCommand(testCmd)
}
