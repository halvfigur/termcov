package main

import (
	"bufio"
	"fmt"
	"go/build"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"

	"golang.org/x/tools/cover"
)

const (
	blackOnWhite = "\033[30;47m"
	green        = "\033[32m"
	red          = "\033[31m"
	clear        = "\033[0m"
)

func findFile(filename string) (string, error) {
	dir, filename := filepath.Split(filename)
	pkg, err := build.Import(dir, ".", build.FindOnly)
	if err != nil {
		return "", fmt.Errorf("can't find %q: %v", filename, err)
	}
	return filepath.Join(pkg.Dir, filename), nil
}

func percentCovered(p *cover.Profile) float64 {
	var total, covered int64
	for _, b := range p.Blocks {
		total += int64(b.NumStmt)
		if b.Count > 0 {
			covered += int64(b.NumStmt)
		}
	}
	if total == 0 {
		return 0
	}
	return float64(covered) / float64(total) * 100
}

func renderFile(filename string) error {
	profiles, err := cover.ParseProfiles(filename)
	if err != nil {
		return err
	}

	for _, profile := range profiles {
		actualFilename, err := findFile(profile.FileName)
		if err != nil {
			return err
		}

		bytes, err := ioutil.ReadFile(actualFilename)
		if err != nil {
			return fmt.Errorf("failed to read: %q, %w", actualFilename, err)
		}

		if err := renderBoundaries(profile, bytes); err != nil {
			return err
		}
	}

	return nil
}

func renderBoundaries(profile *cover.Profile, bytes []byte) error {
	writer := bufio.NewWriter(os.Stdout)

	coveragePercent := strconv.FormatFloat(percentCovered(profile), 'f', 2, 64)

	// Builder write operations always return a nil error
	writer.WriteString("\n" + blackOnWhite + " --- " + profile.FileName + " " + coveragePercent + "%" + " --- " + clear + "\n")

	offset := 0
	for _, boundary := range profile.Boundaries(bytes) {
		writer.Write(bytes[offset:boundary.Offset])

		if boundary.Start {
			if boundary.Count > 0 {
				writer.WriteString(green)
			} else {
				writer.WriteString(red)
			}
		} else {
			writer.WriteString(clear)
		}

		offset = boundary.Offset
	}

	if offset < len(bytes) {
		writer.Write(bytes[offset:])
	}

	return writer.Flush()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage:", os.Args[0], "<file>")
		os.Exit(1)
	}

	if err := renderFile(os.Args[1]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
