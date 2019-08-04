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

	"github.com/k1LoW/metr/metrics"
	"github.com/labstack/gommon/color"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "show available metrics",
	Long:  `show available metrics.`,
	Run: func(cmd *cobra.Command, args []string) {
		if interval == 0 {
			interval = 500
		}
		m, err := metrics.Get(time.Duration(interval) * time.Millisecond)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}
		for _, metric := range metrics.NewMetrics().List() {
			if value, ok := m.Load(metric.Name); ok {
				fmt.Printf("%s (now:%v): %v\n", color.White(metric.Name, color.B), color.Cyan(fmt.Sprintf(metric.Format, value)), metric.Description)
			}
		}
		fmt.Printf("(metric measurement interval: %v ms)\n", color.Cyan(interval))
	},
}

func init() {
	listCmd.Flags().IntVarP(&interval, "interval", "i", 0, "metric measurement interval (millisecond)")
	rootCmd.AddCommand(listCmd)
}
