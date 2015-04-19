package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

var envs = map[string]interface{}{
	"dev": map[string]string{
		"CNAME":  "******.elasticbeanstalk.com",
		"REGION": "ap-northeast-1",
	},
	"prd": map[string]string{
		"CNAME":  "******.elasticbeanstalk.com",
		"REGION": "us-east-1",
	},
}

func main() {
	keys := make([]string, len(envs))
	i := 0
	for k := range envs {
		keys[i] = k
		i++
	}
	options := strings.Join(keys, " | ")
	if len(os.Args) < 2 {
		printUsage(options)
		return
	}
	env := os.Args[1]
	property := envs[env]
	if property == nil {
		printUsage(options)
		return
	}
	CNAME := property.(map[string]string)["CNAME"]
	REGION := property.(map[string]string)["REGION"]
	if CNAME == "" || REGION == "" {
		fmt.Println("CNAME or REGION is Empty")
		return
	}
	environments := getEnvironments(REGION)
	currentEnv := findCurrentEnv(environments, CNAME, REGION)
	if currentEnv == "" {
		fmt.Println("Not Exists Current ENV")
		return
	}

	targetEnvs := getTargetEnvs(environments, currentEnv)
	if len(targetEnvs) == 0 {
		targetEnvs = append(targetEnvs, currentEnv)
	}
	deploy(targetEnvs, REGION)
	if len(environments) > 1 {
		swapEnvironment(REGION, currentEnv, targetEnvs[0])
	}
	fmt.Println("Finished!!!")
}

func getEnvironments(region string) []string {
	out, err := exec.Command("eb", "list", "--region="+region).CombinedOutput()
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

func findCurrentEnv(environments []string, cname string, region string) string {
	for _, env := range environments {
		out, err := exec.Command("eb", "status", env, "--region="+region).CombinedOutput()
		if err != nil {
			printError(out, err)
		}
		values := strings.Split(string(out), "\n")
		for _, v := range values {
			ebCNAME := strings.TrimSpace(strings.Replace(v, "CNAME: ", "", -1))
			if ebCNAME == cname {
				return env
			}
		}
	}
	return ""
}

func getTargetEnvs(environments []string, currentEnv string) []string {
	targetEnvs := []string{}
	for _, env := range environments {
		if currentEnv != env {
			targetEnvs = append(targetEnvs, env)
		}
	}
	return targetEnvs
}

func deploy(targetEnvs []string, region string) {
	for _, env := range targetEnvs {
		fmt.Println("Deploy to " + env)
		out, err := exec.Command("eb", "deploy", env, "--region="+region).CombinedOutput()
		if err != nil {
			printError(out, err)
		}
		fmt.Println(string(out))
	}
}

func swapEnvironment(region string, currentName string, destinationName string) {
	out, err := exec.Command("eb", "swap", currentName, "--destination_name", destinationName, "--region="+region).CombinedOutput()
	if err != nil {
		printError(out, err)
	}
	fmt.Println(string(out))
}

func printUsage(options string) {
	fmt.Println("usage: go run deploy.go [" + options + "]")
}

func printError(out []byte, err error) {
	fmt.Println(string(out))
	fmt.Println(err)
	os.Exit(1)
}
