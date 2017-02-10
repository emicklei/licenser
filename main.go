package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var (
	dryRun    = flag.Bool("d", false, "dry run, see what would change")
	overwrite = flag.Bool("f", false, "overwrite any existing copyright notice")
	goFiles   = flag.String("go", "", "file that contains copyright notice for Go files")
	pwd       string
)

func main() {
	pwd, _ = os.Getwd()
	flag.Parse()
	if len(flag.Args()) == 0 {
		flag.PrintDefaults()
		return
	}
	err := filepath.Walk(flag.Args()[0], visit)
	if err != nil {
		log.Fatalln(err)
	}
}

func visit(path string, f os.FileInfo, err error) error {
	if !f.IsDir() {
		if strings.HasSuffix(f.Name(), ".go") {
			processGoSource(filepath.Join(pwd, path))
		}
	}
	return nil
}

func processGoSource(filename string) error {

	fmt.Printf("Visited: %s\n", filename)

	return nil
}
