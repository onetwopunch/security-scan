package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type scannerEntry struct {
	Description string `json:"description"`
	Regex       string `json:"regex"`
}

type scanList []scannerEntry

func NewScanList(filepath string) []*Scanner {
	var list scanList
	var data []byte
	var err error

	if data, err = ioutil.ReadFile(filepath); err != nil {
		log.Fatal(err)
	}

	if err = json.Unmarshal(data, &list); err != nil {
		log.Fatal(err)
	}
	var scanners []*Scanner
	for _, entry := range list {
		scanner := NewScanner(entry.Description, entry.Regex)
		scanners = append(scanners, scanner)
	}
	return scanners
}
