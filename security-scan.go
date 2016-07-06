package main

import (
	"bytes"
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	var err error
	var outfile *os.File
	var lines []string

	dir := flag.String("git", "", "Git working directory to scan (defaults to current working directory)")
	output := flag.String("o", "security-scan.csv", "Output CSV filename")
	gopath := os.Getenv("GOPATH")
	defaultMatchersPath := filepath.Join(gopath, "src", "github.com", "onetwopunch", "security-scan", "matchers.json")
	matchers := flag.String("m", defaultMatchersPath, "JSON file containing a list of matchers\n\t[\n\t  {\n\t    \"description\":string,\n\t    \"regex\":string\n\t  }, ...\n\t]\n\t")
	scanners := NewScanList(*matchers)
	help := flag.Bool("h", false, "Usage")
	flag.Parse()
	if *help {
		os.Exit(Usage())
	}

	os.MkdirAll(filepath.Dir(*output), 0755)
	outfile, err = os.OpenFile(*output, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Reading Git history from %s...\n", *dir)
	lines, err = gitLog(*dir)
	if err != nil {
		log.Fatal(err)
		Usage()
	}
	writer := csv.NewWriter(outfile)
	header := []string{"Description", "Filename", "Commit", "Author", "Line"}
	writer.Write(header)
	writer.Flush()

	for _, scanner := range scanners {
		fmt.Printf("[%v] Matches Found: 0\r", scanner.description)
		for _, line := range lines {
			if scanner.ScanLine(line) {
				fmt.Printf("[%v] Matches Found: %v\r", scanner.description, len(scanner.matches))
			}
		}
		fmt.Printf("\n")
		writer.WriteAll(scanner.Records())
	}
}

func gitLog(dir string) ([]string, error) {
	var out bytes.Buffer
	var err bytes.Buffer

	cmd := exec.Command("git", "log", "-p")
	cmd.Stderr = &err
	cmd.Stdout = &out
	cmd.Dir = dir
	cmd.Run()
	if err.Len() > 0 {
		return []string{}, fmt.Errorf("%v", err.String())
	}
	lines := strings.Split(out.String(), "\n")
	return lines, nil
}

func Usage() int {
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
	flag.PrintDefaults()
	return 1
}
