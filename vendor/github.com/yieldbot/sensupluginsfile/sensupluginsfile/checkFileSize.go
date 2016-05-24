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
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

func testFile(f string) string {
	file, err := os.Open(f)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// get the file size
	stat, err := file.Stat()
	if err != nil {
		panic(err)
	}

	fmt.Printf("The size of %v  is: %v\n", f, stat.Size())

	if float64(stat.Size()) > float64(1024) {
		fmt.Printf("this is critical\n")
	} else if float64(stat.Size()) > float64(512) {
		fmt.Printf("this is warning\n")
	} else {
		fmt.Printf("this is fine\n")
	}
	return "fine"

}

func visitFile(fp string, fi os.FileInfo, err error) error {
	if err != nil {
		fmt.Println(err) // can't walk here,
		return nil       // but continue walking elsewhere
	}
	if !!fi.IsDir() {
		return nil // not a file.  ignore.
	}
	matched, err := filepath.Match("*.gz", fi.Name())
	if err != nil {
		fmt.Println(err) // malformed pattern
		return err       // this is fatal.
	}
	if !matched {
		testFile(fp)
		// fmt.Println(fp)
	}
	return nil
}

// checkFileSizeCmd represents the checkFileSize command
var checkFileSizeCmd = &cobra.Command{
	Use:   "checkFileSize",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(sensupluginsfile *cobra.Command, args []string) {
		filepath.Walk("/var/log/", visitFile)
		// TODO: Work your own magic here
		fmt.Println("checkFileSize called")
	},
}

func init() {
	RootCmd.AddCommand(checkFileSizeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// checkFileSizeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// checkFileSizeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
