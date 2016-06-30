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

var scanners []*Scanner

func init() {
	scanners = []*Scanner{
		NewScanner("Toml-Syntax", "(password|secret) = ['\"].*?"),
		NewScanner("AWS Access Key ID", "AKIA[0-9A-Z]{16}"),
		NewScanner("Redis URL with Password", "redis://[0-9a-zA-Z:@.\\-]+"),
		NewScanner("Google Access Token", "ya29.[0-9a-zA-Z_\\-]{68}"),
		NewScanner("Google API", "AIzaSy[0-9a-zA-Z_\\-]{33}"),
	}
}

func main() {
	var err error
	var outfile *os.File
	var lines []string

	dir := flag.String("git", "", "Git working directory to scan (defaults to current working directory)")
	output := flag.String("o", "security_scan.csv", "Output CSV filename")
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

	lines, err = gitLog(*dir)
	if err != nil {
		log.Fatal(err)
		Usage()
	}
	writer := csv.NewWriter(outfile)
	header := []string{"Provider", "Filename", "Commit", "Author", "Line"}
	writer.Write(header)
	writer.Flush()

	for _, scanner := range scanners {
		for _, line := range lines {
			if scanner.ScanLine(line) {
				fmt.Printf("[%v] Matches Found: %v\r", scanner.provider, len(scanner.matches))
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
