package main

import (
	"os"
)

// this tool can use cloud-init config file to bootstrap a vm
// without using cloud-init. It of course will only understand a subset
// of cloud-init config as only provided by zos. Any changes to the config
// struct provided by zos will require a similar change here.
func main() {
	seed := os.Args[1]
	root := os.Args[2]

	log("seed directory: %s", seed)
	log("root directory: %s", root)

	if err := ApplyMeta(seed); err != nil {
		log("failed to apply meta: %v", err)
	}

	if err := ApplyUserData(seed, root); err != nil {
		log("failed to apply user data: %v", err)
	}
}
