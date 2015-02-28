package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Production CNAME (ex: ***.elasticbeanstalk.com)
const CNAME = ""

func main() {
	if CNAME == "" {
		fmt.Println("CNAME is Empty")
		os.Exit(1)
	}
	productionEnv := findProductionEnv()
	if productionEnv == "" {
		fmt.Println("Not Exists CNAME")
		os.Exit(1)
	}
	bumpupVersion()
	deploy(productionEnv)
	swapEnvironment()
	fmt.Println("Finished ♬")
}

func getEnvironments() []string {
	out, err := exec.Command("eb", "list").Output()
	if err != nil {
		printError(out, err)
	}
	envs := strings.Split(string(out), "\n")
	environments := []string{}
	for _, env := range envs {
		env = strings.Replace(env, "*", "", -1)
		env = strings.TrimSpace(env)
		if env != "" {
			environments = append(environments, env)
		}
	}
	return environments
}

func findProductionEnv() string {
	productionEnv := ""
	for _, env := range getEnvironments() {
		out, err := exec.Command("eb", "status", env).Output()
		if err != nil {
			printError(out, err)
		}
		values := strings.Split(string(out), "\n")
		for _, v := range values {
			cname := strings.TrimSpace(strings.Replace(v, "CNAME: ", "", -1))
			if cname == CNAME {
				return env
			}
		}
	}
	return productionEnv
}

func bumpupVersion() {
	out, err := exec.Command("npm", "version", "patch", "-m", "Upgrade to %s").Output()
	if err != nil {
		printError(out, err)
	}
	fmt.Println(string(out))
	_, _ = exec.Command("git", "push", "origin", "master").Output()
}

func deploy(productionEnv string) {
	for _, env := range getEnvironments() {
		if productionEnv == env {
			continue
		}
		fmt.Println("Deploy to " + env)
		out, err := exec.Command("eb", "deploy", env).Output()
		if err != nil {
			printError(out, err)
		}
		fmt.Println(string(out))
	}
}

func swapEnvironment() {
	out, err := exec.Command("eb", "swap").Output()
	if err != nil {
		printError(out, err)
	}
	fmt.Println(string(out))
}

func printError(out []byte, err error) {
	fmt.Println(string(out))
	fmt.Println(err)
	os.Exit(1)
}
