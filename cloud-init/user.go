package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const (
	UserDataFile = "user-data"
)

type User struct {
	Name string   `yaml:"name"`
	Keys []string `yaml:"ssh_authorized_keys"`
}

type UserData struct {
	Mounts [][]string `yaml:"mounts"`
	Users  []User     `yaml:"users"`
}

func applyMounts(root string, mounts [][]string) error {
	for _, mount := range mounts {
		if len(mount) != 6 {
			log("mount is not valid: %v", mount)
			continue
		}

		source := mount[0]
		target := mount[1]
		fstype := mount[2]

		if len(target) == 0 || target == "/" {
			log("invalid mount target '%s'", target)
			continue
		}

		if !filepath.IsAbs(target) {
			return fmt.Errorf("invalid target mount must be absolute path: '%s'", target)
		}

		log("mounting %s (%s) on %s", source, fstype, target)

		path := filepath.Join(root, target)
		if err := os.MkdirAll(path, 0755); err != nil {
			return fmt.Errorf("failed to prepare mountpoint '%s'", target)
		}

		// we need to detect the type of the filesystem to do a
		// syscall.Mount because it does not support `auto` type.
		// which means we still need to exec something like `blkid` to detect
		// the filesystem type of the given device (if auto)
		// then do syscall.Mount, so instead we can just exec `mount`
		// command since it know how to do this anyway

		if err := exec.Command("mount", "-t", fstype, source, path).Run(); err != nil {
			return fmt.Errorf("failed to mount device '%s' (%s) on '%s': %w", source, fstype, path, err)
		}
	}

	return nil
}

func applyUsers(root string, users []User) error {
	// currently this code only sets up the user authorized keys.
	// it does not `useradd`

	for _, user := range users {
		// we only support root user for
		// cloud-containers
		if user.Name != "root" {
			continue
		}

		path := filepath.Join(root, "/root")

		if err := os.MkdirAll(path, 0750); err != nil {
			return fmt.Errorf("failed to create root home directory: %w", err)
		}

		path = filepath.Join(root, "root", ".ssh")
		if err := os.MkdirAll(path, 0700); err != nil {
			return fmt.Errorf("failed to create root .ssh directory: %w", err)
		}

		err := ioutil.WriteFile(
			filepath.Join(path, "authorized_keys"),
			[]byte(strings.Join(user.Keys, "\n")),
			0664,
		)

		if err != nil {
			return fmt.Errorf("failed to write authorized_keys file: %w", err)
		}
	}

	return nil
}
func ApplyUserData(seed, root string) error {
	var data UserData

	if err := load(filepath.Join(seed, UserDataFile), &data); err != nil {
		return fmt.Errorf("failed to load user-data file: %w", err)
	}

	if err := applyMounts(root, data.Mounts); err != nil {
		return err
	}

	return applyUsers(root, data.Users)
}
