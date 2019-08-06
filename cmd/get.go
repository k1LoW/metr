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

	"github.com/k1LoW/metr/metrics"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get [METRIC_NAME]",
	Short: "get metric[s]",
	Long:  `get metric[s].`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.WithStack(errors.New("requires one arg"))
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		os.Exit(runGet(args, interval, os.Stdout, os.Stderr))
	},
}

func runGet(args []string, interval int, stdout, stderr io.Writer) (exitCode int) {
	key := args[0]

	m, err := metrics.GetMetrics(time.Duration(interval) * time.Millisecond)
	if err != nil {
		_, _ = fmt.Fprintf(stderr, "%s\n", err)
		return 1
	}
	if len(args) == 1 && key == "all" {
		m.Each(func(metric metrics.Metric, v interface{}) {
			_, _ = fmt.Fprintf(stdout, "%s:%s\n", metric.Name, fmt.Sprintf(metric.Format, v))
		})
		return 0
	}
	v, ok := m.Load(key)
	if !ok {
		_, _ = fmt.Fprintf(stderr, "%s does not exist\n", key)
		return 1
	}
	_, _ = fmt.Fprintf(stdout, "%s\n", fmt.Sprintf(m.Format(key), v))
	return 0
}

func init() {
	getCmd.Flags().IntVarP(&interval, "interval", "i", 500, "metric measurement interval (millisecond)")
	rootCmd.AddCommand(getCmd)
}
