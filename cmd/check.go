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
	"os"
	"time"

	"github.com/antonmedv/expr"
	"github.com/k1LoW/metr/metrics"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
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

var (
	warningCond  string
	criticalCond string
)

// checkCmd represents the check command
var checkCmd = &cobra.Command{
	Use:   "check -w [WARNING_CONDITION] -c [CRITICAL_CONDITION]",
	Short: "check metrics condition and output result with exit status code",
	Long:  `check metrics condition and output result with exit status code.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			exitWithStdout(UNKNOWN, warningCond, criticalCond, errors.New("metr requires no args"))
		}
		if warningCond == "" && criticalCond == "" {
			exitWithStdout(UNKNOWN, warningCond, criticalCond, errors.New("metr requires -w or -c option"))
		}
		m, err := metrics.Get(time.Duration(interval) * time.Millisecond)
		if err != nil {
			exitWithStdout(UNKNOWN, warningCond, criticalCond, err)
		}
		got, err := expr.Eval(fmt.Sprintf("(%s) == true", criticalCond), m.Raw())
		if err != nil {
			exitWithStdout(UNKNOWN, warningCond, criticalCond, err)
		}
		if got.(bool) {
			exitWithStdout(CRITICAL, warningCond, criticalCond, nil)
		}

		got, err = expr.Eval(fmt.Sprintf("(%s) == true", warningCond), m.Raw())
		if err != nil {
			exitWithStdout(UNKNOWN, warningCond, criticalCond, err)
		}
		if got.(bool) {
			exitWithStdout(WARNING, warningCond, criticalCond, nil)
		}

		exitWithStdout(OK, warningCond, criticalCond, nil)
	},
}

func init() {
	rootCmd.AddCommand(checkCmd)
	checkCmd.Flags().StringVarP(&warningCond, "warning", "w", "", "WARNING condition")
	checkCmd.Flags().StringVarP(&criticalCond, "critical", "c", "", "CRITICAL condition")
	checkCmd.Flags().IntVarP(&interval, "interval", "i", 500, "metric measurement interval (millisecond)")
}

func exitWithStdout(status Status, warningCond string, criticalCond string, err error) {
	// SERVICE STATUS: First line of output | First part of performance data http://nagios-plugins.org/doc/guidelines.html#PLUGOUTPUT
	if err != nil {
		fmt.Printf("%s %s: %s\n", "METR", status, err)
		os.Exit(int(status))
	}
	fmt.Printf("%s %s: w:(%s) c:(%s)\n", "METR", status, warningCond, criticalCond)
	os.Exit(int(status))
}
