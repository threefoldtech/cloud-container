package main

import (
	"path/filepath"
	"syscall"
)

const (
	MetaFile = "meta-data"
)

type Metadata struct {
	InstanceID string `yaml:"instance-id"`
	Hostname   string `yaml:"local-hostname"`
}

func ApplyMeta(seed string) error {
	var meta Metadata
	if err := load(filepath.Join(seed, MetaFile), &meta); err != nil {
		return err
	}

	if len(meta.Hostname) == 0 {
		log("host name is not set")
	}

	return syscall.Sethostname([]byte(meta.Hostname))
}
