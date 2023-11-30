package main

import (
	"fmt"
	"os/exec"
)

func Miscellaneous(root string) error {
	// generate host keys
	if err := exec.Command("ssh-keygen", "-A", "-f", root).Run(); err != nil {
		return fmt.Errorf("failed to generate host keys '%w'", err)
	}

	return nil
}
