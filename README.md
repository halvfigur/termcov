# termcov
termcov takes a Golang coverage profile, as produced by e.g. `go test -coverprofile=coverage.out`, and writes colorized
output to the terminal.

## Installation
```bash
go get github.com/halvfigur/termcov
```

## Usage
```bash
$ go test -coverprofile=coverage.out
$ termcov coverage.out
```
