/*
Copyright 2017 Ernest Micklei.

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

package main

import (
	"bufio"
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
	oNoticeFilename = flag.String("f", "", "filename that contains the copyright notice")
	oStarStyle      = flag.Bool("s", false, "if true then use the /* ... */ method for writing the notice else use //")
	oExtension      = flag.String("e", ".go", "file extension for which the copyright notice must be added")
	oRecurseDirs    = flag.Bool("r", false, "recursively search for files")
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
	if *oRecurseDirs {
		err = filepath.Walk(flag.Args()[0], visit)
		if err != nil {
			log.Fatalln(err)
		}

	} else {
		path := flag.Args()[0]
		dir, err := os.Open(path)
		if err != nil {
			log.Fatalln(err)
		}
		list, err := dir.Readdir(0)
		if err != nil {
			log.Fatalln(err)
		}
		for _, each := range list {
			if err := visit(filepath.Join(path, each.Name()), each, nil); err != nil {
				log.Fatalln(err)
			}
		}
	}
	waiters.Wait()
}

func visit(path string, f os.FileInfo, err error) error {
	if !f.IsDir() {
		if strings.HasSuffix(f.Name(), *oExtension) {
			waiters.Add(1)
			go func(path string) {
				processSource(filepath.Join(pwd, path), f.Mode())
				waiters.Done()
			}(path)
		}
	}
	return nil
}

func processSource(filename string, mode os.FileMode) error {
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

	writeNoticeOn(out)

	_, err = out.Write(data)
	if err != nil {
		fmt.Printf("failed to write: %s,%v\n", filename, err)
		return err
	}

	fmt.Printf("writing: %s\n", filename)
	return nil
}

func writeNoticeOn(w io.Writer) error {
	if *oStarStyle {
		io.WriteString(w, "/*\n")
		io.WriteString(w, notice)
		_, err := io.WriteString(w, "\n*/\n\n")
		return err
	}
	scanner := bufio.NewScanner(strings.NewReader(notice))
	for scanner.Scan() {
		io.WriteString(w, "// ")
		io.WriteString(w, scanner.Text())
		io.WriteString(w, "\n")
	}
	io.WriteString(w, "\n")
	return scanner.Err()
}
