/*
Package cmd is a part of txt2automator
Copyright © 2019 Kevin Bridges kevin@cyberswat.com

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

// Min is used when determing delay between keystrokes
var Min int

// Max is used when determing delay between keystrokes
var Max int

// Dir is used to read all files from instead of converting individual files.
var Dir string

// convertCmd represents the convert command
var convertCmd = &cobra.Command{
	Use:   "convert <file>",
	Short: "convert a file into applescript",
	Long:  `Takes the given file and converts into an applescript file.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 && Dir == "" {
			return errors.New("no file or directory given to convert")
		}
		if len(args) != 0 && Dir != "" {
			return errors.New("please specify a file or a directory to convert instead of both")
		}

		if len(args) != 0 {
			for i := 0; i < len(args); i++ {
				err := processPath(args[i])
				if err != nil {
					return err
				}
			}
		}

		if Dir != "" {
			fmt.Println("Need to convert " + Dir)
			files, err := ioutil.ReadDir(Dir)
			if err != nil {
				log.Fatal(err)
			}

			for _, f := range files {
				if strings.HasPrefix(f.Name(), ".") == false {
					absPath := filepath.Join(Dir, f.Name())
					err := processPath(absPath)
					if err != nil {
						return err
					}
				}
			}
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(convertCmd)
	convertCmd.Flags().IntVar(&Min, "min", 0, "Minimum value when determing delay between keystrokes")
	convertCmd.Flags().IntVar(&Max, "max", 3, "Maximum value when determing delay between keystrokes")
	convertCmd.Flags().StringVarP(&Dir, "dir", "d", "", "Source directory to read all files from")
}

func convert(path string) (string, string, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return "", "", err
	}

	file := filepath.Base(path)
	re := regexp.MustCompile("\\.")
	fileName := re.ReplaceAllString(file, "-") + ".scpt"

	s := string(data)
	r := bufio.NewReader(strings.NewReader(s))
	newScript := "on run {input, parameters}\n\n"
	newScript += "\ttell application \"iTerm\" to activate\n\n"
	newScript += "\ttell application \"System Events\"\n"
	for {
		if c, _, err := r.ReadRune(); err != nil {
			if err == io.EOF {
				break
			} else {
				log.Fatal(err)
			}
		} else {
			var character = string(c)
			if character == "\n" {
				newScript += "\t\tkey code 36\n"
			} else if character == "\"" {
				newScript += "\t\tkeystroke \"\\" + character + "\"\n"
			} else {
				newScript += "\t\tkeystroke \"" + character + "\"\n"
			}
			rand.Seed(time.Now().UnixNano())
			randNum := rand.Intn(Max - Min)
			newScript += fmt.Sprintf("\t\tdelay 0.%d\n", randNum)
		}
	}
	newScript += "\tend tell\n\n"
	newScript += "\treturn input\n"
	newScript += "end run\n"

	return newScript, fileName, nil
}

func writeFile(contents string, fileName string) {
	f, err := os.Create(fileName)
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = f.WriteString(contents)
	if err != nil {
		fmt.Println(err)
		f.Close()
		return
	}
	err = f.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
}

func processPath(path string) error {
	newScript, fileName, err := convert(path)
	if err != nil {
		return err
	}
	writeFile(newScript, fileName)
	fmt.Printf("%s successfully converted to %s\n", path, fileName)
	return nil
}
