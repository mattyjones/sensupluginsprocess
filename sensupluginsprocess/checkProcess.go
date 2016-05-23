// Copyright © 2016 Yieldbot <devops@yieldbot.com>
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

package sensupluginsprocess

import (
	"os"

	"github.com/op/go-logging"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/yieldbot/sensuplugin/sensuutil"
	"github.com/yieldbot/sensupluginsfile/sensupluginsfile"
)

// JavaApp is used to let the process -> pid function know how it will match the process name
var JavaApp = sensupluginsfile.JavaApp

var app string

var syslogLog = logging.MustGetLogger("sensu-process-check")
var stderrLog = logging.MustGetLogger("sensu-process-check")

// checkProcessCmd represents the checkProcess command
var checkProcessCmd = &cobra.Command{
	Use:   "checkProcess",
	Short: "Check to see if a process is running",
	Long: ` Check to see if a process is running by using ps to determine if the
  process has a pid and is in fact running. The usual service foo status is not
  used in this case due to redirects using runnit.`,
	Run: func(sensupluginsprocess *cobra.Command, args []string) {

		syslogBackend, _ := logging.NewSyslogBackend("checkProcess")
		stderrBackend := logging.NewLogBackend(os.Stderr, "checkProcess", 0)
		syslogBackendFormatter := logging.NewBackendFormatter(syslogBackend, sensuutil.SyslogFormat)
		stderrBackendFormatter := logging.NewBackendFormatter(stderrBackend, sensuutil.StderrFormat)
		logging.SetBackend(syslogBackendFormatter)
		logging.SetBackend(stderrBackendFormatter)

		var appPid string

		switch app {
		case "":
			if viper.GetString("sensupluginsprocess.checkProcess.app") != "" {
				app = viper.GetString("sensupluginsprocess.checkProcess.app")
				appPid = sensupluginsfile.GetPid(app)
			} else {
				syslogLog.Error(`You are missing a required configuration parameter\n
          If unsure consult the documentation for examples and requirements\n`)
				os.Exit(sensuutil.MonitoringErrorCodes["CONFIG_ERROR"])
			}
		default:
			appPid = sensupluginsfile.GetPid(app)
		}

		if appPid == "" {
			sensuutil.Exit("critical", app+" is not running")
		} else {
			sensuutil.Exit("ok")
		}
	},
}

func init() {
	RootCmd.AddCommand(checkProcessCmd)
	checkProcessCmd.Flags().StringVarP(&app, "app", "", "sbin/init", "the process name")
}
