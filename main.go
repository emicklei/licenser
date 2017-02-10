// Copyright 2017 Ernest Micklei.
// 
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// 
//     http://www.apache.org/licenses/LICENSE-2.0
// 
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var (
	oDryRun         = flag.Bool("d", false, "dry run, see what would change")
	oExtension      = flag.String("e", ".go", "file extension for which the copyright notice must be added")
	oSeparatorToken = flag.String("t", "package", "source token that indicates where the actual source will start")
	oNoticeFilename = flag.String("f", "", "source token that indicates where the actual source will start")
)

var (
	pwd     string
	waiters *sync.WaitGroup
	notice  string
)

func main() {
	pwd, _ = os.Getwd()
	flag.Parse()
	if len(flag.Args()) == 0 {
		fmt.Println("Usage: licenser [flags] [path...]", "\n")
		flag.PrintDefaults()
		return
	}

	// reading copyright notice
	data, err := ioutil.ReadFile(*oNoticeFilename)
	if err != nil {
		log.Fatalf("failed to read copyright notice: %s because: %v\n", *oNoticeFilename, err)
	}
	notice = string(data)

	// find and update all files
	waiters = new(sync.WaitGroup)
	err = filepath.Walk(flag.Args()[0], visit)
	if err != nil {
		log.Fatalln(err)
	}
	waiters.Wait()
}

func visit(path string, f os.FileInfo, err error) error {
	if !f.IsDir() {
		if strings.HasSuffix(f.Name(), *oExtension) {
			waiters.Add(1)
			go func(path string) {
				processSource(filepath.Join(pwd, path), f.Mode(), *oSeparatorToken)
				waiters.Done()
			}(path)
		}
	}
	return nil
}

func processSource(filename string, mode os.FileMode, token string) error {
	if *oDryRun {
		fmt.Printf("visiting: %s\n", filename)
		return nil
	}
	fmt.Printf("reading: %s\n", filename)
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Printf("failed to read: %s,%v\n", filename, err)
		return nil
	}
	out, err := os.Create(filename)
	if err != nil {
		fmt.Printf("failed to write: %s,%v\n", filename, err)
		return err
	}
	defer out.Close()

	scanner := bufio.NewScanner(bytes.NewReader(data))
	seenPackage := false
	for scanner.Scan() {
		line := scanner.Text()
		if seenPackage {
			io.WriteString(out, line)
			io.WriteString(out, "\n")
		} else {
			seenPackage = strings.HasPrefix(strings.TrimLeft(line, " \t"), token) // allow whitespace before pkg
			if seenPackage {
				writeNoticeOn(out)
				io.WriteString(out, line)
				io.WriteString(out, "\n")
			} else {
				// ignore lines before package
			}
		}
	}
	if scanner.Err() != nil {
		fmt.Printf("failed to write: %s,%v\n", filename, scanner.Err())
		return scanner.Err()
	}
	fmt.Printf("writing: %s\n", filename)
	return nil
}

func writeNoticeOn(w io.Writer) error {
	scanner := bufio.NewScanner(strings.NewReader(notice))
	for scanner.Scan() {
		io.WriteString(w, "// ")
		io.WriteString(w, scanner.Text())
		io.WriteString(w, "\n")
	}
	io.WriteString(w, "\n")
	return scanner.Err()
}
