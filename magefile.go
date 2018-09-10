// +build mage

package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/magefile/mage/mg" // mg contains helpful utility functions, like Deps
)

var samFilePath = "./sam.yml"

type function struct {
	name string
	path string
}

var funcList = []function{
	{"myfunc", "./functions/myfunc/"},
}

// A build step that requires additional params, or platform specific steps for example
func Build() error {
	fmt.Println("Building...")
	for _, f := range funcList {
		fmt.Println("* ", f.name)
		cmd := exec.Command("go", "build", "-o", "build/"+f.name, f.path)
		cmd.Env = append(os.Environ(), "GOARCH=amd64", "GOOS=linux")
		if err := cmd.Run(); err != nil {
			return err
		}
	}
	return nil
}

// Clean up after yourself
func Clean() {
	fmt.Println("Cleaning...")
	os.RemoveAll("build")
}

type config struct {
	StackName    string
	CodeS3Bucket string
	CodeS3Prefix string
	Parameters   []string
}

func loadConfig() (cfg *config, err error) {
	configFile := os.Getenv("CONFIG_FILE")
	if configFile == "" {
		log.Println("CONFIG_FILE is not available, use 'param.cfg' instead")
		configFile = "./param.cfg"
	}

	cfg = &config{}
	cfg.Parameters = []string{}

	fp, err := os.Open(configFile)
	if err != nil {
		return
	}
	defer fp.Close()

	scanner := bufio.NewScanner(fp)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if len(line) == 0 {
			continue
		}

		idx := strings.Index(line, "=")
		if idx < 0 {
			log.Printf("Warning, invalid format of cfg file: '%s'\n", line)
			continue
		}

		key := line[:idx]
		value := line[(idx + 1):]

		switch key {
		case "StackName":
			cfg.StackName = value
		case "CodeS3Bucket":
			cfg.CodeS3Bucket = value
		case "CodeS3Prefix":
			cfg.CodeS3Prefix = value
		default:
			cfg.Parameters = append(cfg.Parameters, line)
		}
	}

	return cfg, nil
}

// Package creates an archive and upload it to S3
func Package() error {
	mg.Deps(Build)

	cfg, err := loadConfig()
	if err != nil {
		return err
	}

	templateFile := "./template.yml"

	fmt.Println("Packaging...")
	pkgCmd := exec.Command("aws", "cloudformation", "package",
		"--template-file", templateFile,
		"--s3-bucket", cfg.CodeS3Bucket,
		"--s3-prefix", cfg.CodeS3Prefix,
		"--output-template-file", samFilePath)

	pkgOut, err := pkgCmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error: ", string(pkgOut))
		fmt.Println(err)
		return err
	}
	fmt.Println("Generated template file: ", samFilePath)

	return nil
}

// Deploying CloudFormation stack
func Deploy() error {
	mg.Deps(Package)

	cfg, err := loadConfig()
	if err != nil {
		return err
	}

	// fmt.Printf("Package > %s", string(pkgOut))
	fmt.Println("Deploy...")
	args := []string{
		"cloudformation", "deploy",
		"--template-file", samFilePath,
		"--stack-name", cfg.StackName,
		"--capabilities", "CAPABILITY_IAM",
		"--parameter-overrides",
	}
	args = append(args, cfg.Parameters...)
	deployCmd := exec.Command("aws", args...)

	deployOut, err := deployCmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error: ", string(deployOut))
		fmt.Println(err)
		return err
	}

	fmt.Println("Done!")
	return nil
}
