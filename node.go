package main

import (
	"fmt"
	"os/exec"
)

func nodejs() {
	bumpupVersion()
}

func bumpupVersion() {
	out, err := exec.Command("npm", "version", "patch", "-m", "Upgrade to %s").Output()
	if err != nil {
		printError(out, err)
	}
	fmt.Println(string(out))
	_, _ = exec.Command("git", "push", "origin", "master").Output()
}
