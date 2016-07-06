## Security Scan

This simple Go tool originally forked from https://github.com/stefansundin/secrets-scanner scans a Git history for things that look like passwords or secrets, and outputs the results complete with filename, commit SHA, author, etc to a CSV file.

It is completely customizable in that you can input your own custom matcher if those that are given are not enough.

### Installation

To install just run:

```
go install github.com/onetwopunch
```

### Usage

Navigate to the git repo working directory you want to scan and run:

```
security-scan
```

You'll get a text report that look like this:

```
Reading Git history from /Users/ryan/Workspace/repo...
[Password or Secret Assignment] Matches Found: 37
[AWS Access Key ID] Matches Found: 85
[Redis URL with Password] Matches Found: 19
[URL Basic auth] Matches Found: 18
[Google Access Token] Matches Found: 8
[Google API] Matches Found: 0
[Slack API] Matches Found: 0
[Slack Bot] Matches Found: 0
[Gem Fury v1] Matches Found: 0
[Gem Fury v2] Matches Found: 0
```

And a CSV file will be generated in the current directory. For more details about usage, just run:

```
$ security-scan -h

Usage of security-scan:
  -git string
    	Git working directory to scan (defaults to current working directory)
  -h	Usage
  -m string
    	JSON file containing a list of matchers
	[
	  {
	    "description":string,
	    "regex":string
	  }, ...
	]
	 (default "$GOPATH/src/github.com/onetwopunch/security-scan/matchers.json")
  -o string
    	Output CSV filename (default "security-scan.csv")
```

The matchers file is pretty self-explanatory: you can add or remove matchers to fit your org's needs. If you need to make changes, just run:

```
cp $GOPATH/src/github.com/onetwopunch/security-scan/matchers.json my_matchers.json
```

Then edit it to your heart's desire and to see the results, run:

```
security-scan -m my_matchers.json -git /path/to/repo
```

### Contributing

If you feel like your matcher is generic enough to add to the default, please feel free to submit a PR. As well, PR's are always welcome.
