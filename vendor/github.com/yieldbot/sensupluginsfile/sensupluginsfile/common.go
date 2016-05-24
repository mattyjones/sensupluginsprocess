// Get the number of open files for a process and compare that against /proc/<pid>/limits and alert if
// over the given threshold.
//
//
// LICENSE:
//   Copyright 2015 Yieldbot. <devops@yieldbot.com>
//   Released under the MIT License; see LICENSE
//   for details.

package sensupluginsfile

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"github.com/yieldbot/sensuplugin/sensuutil"
)

// JavaApp  This is used to let the process -> pid function know how it will match the process name
var JavaApp bool

// Standalone should be set to true if this is not being used externally.
// When being used internally, if an appPid is not found then it will simply
// raise a CONFIG_ERROR. If the GetPid function is used in other programs the
// developer may wish to have an error message and seperate exit code.
var Standalone bool

//GetPid returns the pid for the desired process
func GetPid(app string) string {
	// RED the match is not working for non-java apps. If the string is not matched 100% it will fail.
	JavaApp = true
	pidExp := regexp.MustCompile("[0-9]+")
	termExp := regexp.MustCompile(`pts/[0-9]`)
	appPid := ""

	/// the pid for the binary
	goPid := os.Getpid()
	// if Debug {
	// 	fmt.Printf("golang binary pid: %v\n", goPid)
	// }

	psAEF := exec.Command("ps", "-aef")

	out, err := psAEF.Output()
	if err != nil {
		panic(err)
	}

	psAEF.Start()

	lines := strings.Split(string(out), "\n")

	if !JavaApp {
		for i := range lines {
			if !strings.Contains(lines[i], strconv.Itoa(goPid)) && !termExp.MatchString(lines[i]) {
				words := strings.Split(lines[i], " ")
				for j := range words {
					if app == words[j] {
						appPid = pidExp.FindString(lines[i])
					}
				}
			}
		}
	} else {
		for i := range lines {
			if strings.Contains(lines[i], app) && !strings.Contains(lines[i], strconv.Itoa(goPid)) && !termExp.MatchString(lines[i]) {
				appPid = pidExp.FindString(lines[i])

			}
		}
	}
	if appPid == "" {
		if Standalone {
			sensuutil.ConfigError()
		} else {
			return ""
		}
	}
	return appPid
}

//GetFileHandles returns the current number of open file handles for a process
func GetFileHandles(pid string) (float64, float64, float64) {
	var _s, _h string
	var s, h float64
	limitExp := regexp.MustCompile("[0-9]+")
	filename := `/proc/` + pid + `/limits`
	fdLoc := "/proc/" + pid + "/fd"
	numFD := 0.0

	limits, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	lines := strings.Split(string(limits), "\n")
	for i := range lines {
		if strings.Contains(lines[i], "open files") {
			limits := limitExp.FindAllString(lines[i], 2)
			_s = limits[0]
			_h = limits[1]

			s, err = strconv.ParseFloat(_s, 64)
			if err != nil {
				fmt.Println("warning", "I can't parse the soft limit")
				panic(err)

			}
			h, err = strconv.ParseFloat(_h, 64)
			if err != nil {
				fmt.Println("warning", "I can't parse the hard limit")
				panic(err)

			}
		}
	}

	files, _ := ioutil.ReadDir(fdLoc)
	for _ = range files {
		numFD++
	}
	if numFD == 0.0 {
		fmt.Printf("There are no open file descriptors for the process, did you use sudo?\n")
		fmt.Printf("If unsure of the use, consult the documentation for examples and requirements\n")
		sensuutil.Exit("PERMISSIONERROR")
	}
	// fmt.Printf("s: %v\n",s)
	// fmt.Printf("h: %v\n",h)
	// fmt.Printf("numFD: %v\n",numFD)
	return s, h, numFD
}
