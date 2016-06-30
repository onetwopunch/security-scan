package main

import (
	"regexp"
	"strings"
)

type Match struct {
	line     string
	filename string
	commit   string
	author   string
}

func (match *Match) ArrayWithProvider(provider string) []string {
	return []string{
		provider,
		match.filename,
		match.commit,
		match.author,
		match.line,
	}
}

type Scanner struct {
	provider     string
	re           *regexp.Regexp
	matches      []Match
	currAuthor   string
	currFilename string
	currCommit   string
}

var reFilename *regexp.Regexp = regexp.MustCompile("[+-]{3} (a/(.*?$)|b/(.*?$))")
var reCommit *regexp.Regexp = regexp.MustCompile("^commit ([a-f0-9]{40})")
var reAuthor *regexp.Regexp = regexp.MustCompile("Author: (.*?) <")

func NewScanner(provider string, pattern string) *Scanner {
	scanner := Scanner{provider: provider}
	scanner.re = regexp.MustCompile(pattern)
	return &scanner
}

func (me *Scanner) ScanLine(line string) bool {
	var ok bool
	var filename, commit, author string

	if ok, filename = getFilename(line); ok {
		me.currFilename = filename
		return false
	}

	if ok, commit = getCommit(line); ok {
		me.currCommit = commit
		return false
	}

	if ok, author = getAuthor(line); ok {
		me.currAuthor = author
		return false
	}

	match := me.re.FindString(line)
	if len(match) > 0 {
		trimmed := strings.Trim(line, "\n")
		match := Match{
			line:     trimmed,
			author:   me.currAuthor,
			commit:   me.currCommit,
			filename: me.currFilename,
		}
		me.matches = append(me.matches, match)
		return true
	}
	return false
}

func (me *Scanner) Records() [][]string {
	res := [][]string{}
	for _, match := range me.matches {
		res = append(res, match.ArrayWithProvider(me.provider))
	}
	return res
}

func getFilename(line string) (bool, string) {
	// We can do it this way because if the filename changes
	// we will care what it changed to more than what it
	// changed from so the b/filename will overwrite unless
	// the file was deleted in which case it will be /dev/null
	// Otherwise, the filenames will be the same

	matches := reFilename.FindStringSubmatch(line)
	if len(matches) < 4 {
		return false, ""
	}
	if len(matches[2]) > 0 {
		return true, matches[2]
	}
	return true, matches[3]
}

func getCommit(line string) (bool, string) {
	matches := reCommit.FindStringSubmatch(line)
	if len(matches) > 1 {
		return true, matches[1]
	}
	return false, ""
}

func getAuthor(line string) (bool, string) {
	matches := reAuthor.FindStringSubmatch(line)
	if len(matches) > 1 {
		return true, matches[1]
	}
	return false, ""
}
