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

	"github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/yieldbot/sensuplugin/sensuutil"
	//"github.com/yieldbot/sensupluginsfile/version"
)

// app the check the handles on
var app string

// thresholds to alert on
var warnThreshold int
var critThreshold int

// print debugging info
var debug bool

//determineThreshold will test if the number of file handles is over the allowable amount
func determineThreshold(limit float64, threshold float64, numFD float64) bool {
	alarm := true
	tLimit := threshold / float64(100) * limit

	if numFD > float64(tLimit) {
		alarm = true
	} else {
		alarm = false
	}
	return alarm
}

// checkFileHandlesCmd represents the checkFileHandles command
var checkFileHandlesCmd = &cobra.Command{
	Use:   "checkFileHandles",
	Short: "Check the number of open file handles on a process",
	Run: func(sensupluginsfile *cobra.Command, args []string) {

		// the pid for the app we want to check
		var appPid string

		// values used in determing the alert state
		var sLimit, hLimit, openFd float64

		// Standalone tells the check not to error if an appPid is not found.
		// This is here as several other checks and packages use this functionality.
		Standalone = true

		switch app {
		case "":
			if viper.GetString("sensupluginsfile.checkFileHandles.app") != "" {
				app = viper.GetString("sensupluginsfile.checkHandles.app")
				appPid = GetPid(app)
				sLimit, hLimit, openFd = GetFileHandles(appPid)
				if debug {
					fmt.Printf("warning threshold: %v percent, critical threshold: %v percent\n", warnThreshold, critThreshold)
					fmt.Printf("this is the number of open files at the specific point in time: %v\n", openFd)
					fmt.Printf("app pid is: %v\n", appPid)
					fmt.Printf("This is the soft limit: %v\n", sLimit)
					fmt.Printf("This is the hard limit: %v\n", hLimit)
					sensuutil.Exit("debug")
				}
				if determineThreshold(hLimit, float64(critThreshold), openFd) {
					fmt.Printf("%v is over %v percent of the the open file handles hard limit of %v\n", app, critThreshold, hLimit)
					sensuutil.Exit("critical")
				} else if determineThreshold(sLimit, float64(warnThreshold), openFd) {
					fmt.Printf("%v is over %v percent of the open file handles soft limit of %v\n", app, warnThreshold, sLimit)
					sensuutil.Exit("warning")
				} else {
					sensuutil.Exit("ok")
				}

			} else {
				syslogLog.WithFields(logrus.Fields{
					"check":   "checkFileHandles",
					"client":  host,
					//"version": version.AppVersion(),
				}).Error(`You are missing a required configuration parameter. If unsure consult the documentation for examples and requirements`)
				sensuutil.Exit("CONFIGERROR")
			}
		default:
			appPid = GetPid(app)
			sLimit, hLimit, openFd = GetFileHandles(appPid)
			if debug {
				fmt.Printf("warning threshold: %v percent, critical threshold: %v percent\n", warnThreshold, critThreshold)
				fmt.Printf("this is the number of open files at the specific point in time: %v\n", openFd)
				fmt.Printf("app pid is: %v\n", appPid)
				fmt.Printf("This is the soft limit: %v\n", sLimit)
				fmt.Printf("This is the hard limit: %v\n", hLimit)
				sensuutil.Exit("debug")
			}
			if determineThreshold(hLimit, float64(critThreshold), openFd) {
				fmt.Printf("%v is over %v percent of the the open file handles hard limit of %v\n", app, critThreshold, hLimit)
				sensuutil.Exit("critical")
			} else if determineThreshold(sLimit, float64(warnThreshold), openFd) {
				fmt.Printf("%v is over %v percent of the open file handles soft limit of %v\n", app, warnThreshold, sLimit)
				sensuutil.Exit("warning")
			} else {
				sensuutil.Exit("ok")
			}
		}
	},
}

func init() {

	RootCmd.AddCommand(checkFileHandlesCmd)

	// set commandline flags
	checkFileHandlesCmd.Flags().StringVarP(&app, "app", "", "", "the process name")
	checkFileHandlesCmd.Flags().IntVarP(&warnThreshold, "warn", "", 75, "the alert warning threshold percentage")
	checkFileHandlesCmd.Flags().IntVarP(&critThreshold, "crit", "", 75, "the alert critical threshold percentage")
	checkFileHandlesCmd.Flags().BoolVarP(&JavaApp, "java", "", false, "java apps process detection is different")
	checkFileHandlesCmd.Flags().BoolVarP(&debug, "debug", "", false, "print debugging info instead of an alert")
}
