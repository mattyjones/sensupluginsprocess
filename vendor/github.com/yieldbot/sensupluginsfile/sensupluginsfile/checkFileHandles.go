// Copyright © 2016 Yieldbot
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package sensupluginsfile

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/nitro"
	"github.com/yieldbot/sensuplugin/sensuutil"
)

var app string
var warnThreshold int
var critThreshold int

var timer *nitro.B
var debug bool

func determineThreshold(limit float64, threshold float64, numFD float64) bool {
	alarm := true
	tLimit := threshold / float64(100) * limit

	if numFD > float64(tLimit) {
		alarm = true
	} else {
		alarm = false
	}
	// fmt.Printf("alarm: %v\n", alarm)
	return alarm
}

// checkFileHandlesCmd represents the checkFileHandles command
var checkFileHandlesCmd = &cobra.Command{
	Use:   "checkFileHandles",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(sensupluginsfile *cobra.Command, args []string) {

		var appPid string
		var sLimit, hLimit, openFd float64

		// Standalone tells the app not to error if an appPid is not found.
		Standalone = true

		// need to find a way to get the printf stuff into the sensu exit function
		if app != "" {
			appPid = GetPid(app)
			sLimit, hLimit, openFd = GetFileHandles(appPid)
			if debug {
				fmt.Printf("warning threshold: %v percent, critical threshold: %v percent\n", warnThreshold, critThreshold)
				fmt.Printf("this is the number of open files at the specific point in time: %v\n", openFd)
				fmt.Printf("app pid is: %v\n", appPid)
				fmt.Printf("This is the soft limit: %v\n", sLimit)
				fmt.Printf("This is the hard limit: %v\n", hLimit)
				timer.Step("exit debug")
				sensuutil.Exit("debug")
			}
			if determineThreshold(hLimit, float64(critThreshold), openFd) {
				fmt.Printf("%v is over %v percent of the the open file handles hard limit of %v\n", app, critThreshold, hLimit)
				timer.Step("exit critical")
				sensuutil.Exit("critical")
			} else if determineThreshold(sLimit, float64(warnThreshold), openFd) {
				fmt.Printf("%v is over %v percent of the open file handles soft limit of %v\n", app, warnThreshold, sLimit)
				timer.Step("exit warning")
				sensuutil.Exit("warning")
			} else {
				timer.Step("exit ok")
				sensuutil.Exit("ok", "I'd far rather be happy than right any day")
			}
		} else {

			fmt.Printf("Please enter a process name to check. \n")
			fmt.Printf("If unsure consult the documentation for examples and requirements\n")
			sensuutil.Exit("configerror")
		}
	},
}

func init() {
	timer = nitro.Initalize()

	RootCmd.AddCommand(checkFileHandlesCmd)

	// set commandline flags
	checkFileHandlesCmd.Flags().StringVarP(&app, "app", "", "sbin/init", "the process name")
	checkFileHandlesCmd.Flags().IntVarP(&warnThreshold, "warn", "", 75, "the alert warning threshold percentage")
	checkFileHandlesCmd.Flags().IntVarP(&critThreshold, "crit", "", 75, "the alert critical threshold percentage")
	checkFileHandlesCmd.Flags().BoolVarP(&JavaApp, "java", "", false, "java apps process detection is different")
	checkFileHandlesCmd.Flags().BoolVarP(&debug, "debug", "", false, "print debugging info instead of an alert")
	checkFileHandlesCmd.Flags().BoolVar(&nitro.AnalysisOn, "stepAnalysis", false, "display memory and timing of different steps of the program")
}
