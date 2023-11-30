package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func Miscellaneous(root string) error {
	// generate host keys
	log("generating host ssh keys")
	if err := os.MkdirAll(filepath.Join(root, "etc", "ssh"), 0755); err != nil {
		log("failed to create /etc/ssh: %s", err)
	}

	// for some reason, ssh-keygen does not work UNLESS the root user exists in the /etc/passwod file
	// this runs inside the initramfs so it's safe to just create a default passwrd file
	if err := os.WriteFile("/etc/passwd", []byte("root:x:0:0:root:/:\n"), 0644); err != nil {
		return fmt.Errorf("failed to prepare for host generation")
	}

	if err := exec.Command("/usr/bin/ssh-keygen", "-A", "-f", root).Run(); err != nil {
		return fmt.Errorf("failed to generate host keys '%w'", err)
	}

	return nil
}
