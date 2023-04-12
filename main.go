package main

import (
	"fmt"
	"go/build"
	"io/ioutil"
	"path/filepath"

	"golang.org/x/tools/cover"
)

func renderFile(filename string) error {
	profiles, err := cover.ParseProfiles(filename)
	if err != nil {
		return err
	}

	for _, profile := range profiles {
		actualFilename, err := findFile(profile.FileName)
		bytes, err := ioutil.ReadFile(actualFilename)
		if err != nil {
			return fmt.Errorf("failed to read file %s: %v", profile.FileName, err)
		}

		if err := renderBoundaries(profile, bytes); err != nil {
			return err
		}
	}

	return nil
}

func findFile(filename string) (string, error) {
	dir, filename := filepath.Split(filename)
	pkg, err := build.Import(dir, ".", build.FindOnly)
	if err != nil {
		return "", fmt.Errorf("can't find %q: %v", filename, err)
	}
	return filepath.Join(pkg.Dir, filename), nil
}

func renderBoundaries(profile *cover.Profile, bytes []byte) error {
	fmt.Println("File:", profile.FileName)

	offset := 0
	for _, boundary := range profile.Boundaries(bytes) {
		if boundary.Start {
			fmt.Print(string(bytes[offset:boundary.Offset]))
			if boundary.Count > 0 {
				fmt.Print("\033[32m")
			} else {
				fmt.Print("\033[31m")
			}
			offset = boundary.Offset
		} else {
			fmt.Print(string(bytes[offset:boundary.Offset]))
			fmt.Print("\033[0m")
			offset = boundary.Offset
		}
	}

	if offset < len(bytes) {
		fmt.Print(string(bytes[offset:]))
	}

	return nil
}

func main() {
	err := renderFile("coverage.out")
	if err != nil {
		panic(err)
	}
}
