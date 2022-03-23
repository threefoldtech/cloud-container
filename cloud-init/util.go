package main

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

func load(file string, out interface{}) error {
	fd, err := os.Open(file)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %w", file, err)
	}
	defer fd.Close()

	return yaml.NewDecoder(fd).Decode(out)
}

func log(format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, format, a...)
	fmt.Fprintf(os.Stderr, "\n")
}
