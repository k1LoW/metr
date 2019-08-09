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
	"github.com/k1LoW/metr/metrics"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
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
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		os.Exit(runCond(args, interval, pid, os.Stdout, os.Stderr))
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
		os.Exit(runCond(args, interval, pid, os.Stdout, os.Stderr))
	},
}

func runCond(args []string, interval int, pid int32, stdout, stderr io.Writer) (exitCode int) {
	mcond := args[0]
	m, err := metrics.GetMetrics(time.Duration(interval)*time.Millisecond, pid)
	if err != nil {
		_, _ = fmt.Fprintf(stderr, "%s\n", err)
		return 1
	}
	got, err := expr.Eval(fmt.Sprintf("(%s) == true", mcond), m.Raw())
	if err != nil {
		_, _ = fmt.Fprintf(stderr, "%s\n", err)
		return 1
	}
	if got.(bool) {
		return 0
	} else {
		return 1
	}
}

func init() {
	condCmd.Flags().IntVarP(&interval, "interval", "i", 500, "metric measurement interval (millisecond)")
	condCmd.Flags().Int32VarP(&pid, "pid", "p", 0, "PID of target process")
	testCmd.Flags().IntVarP(&interval, "interval", "i", 500, "metric measurement interval (millisecond)")
	testCmd.Flags().Int32VarP(&pid, "pid", "p", 0, "PID of target process")
	rootCmd.AddCommand(condCmd)
	rootCmd.AddCommand(testCmd)
}
