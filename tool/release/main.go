package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/Sn0wo2/OpenCloudflareCDN/internal/util"
)

func runCmd(name string, args ...string) (string, error) {
	output, err := exec.Command(name, args...).CombinedOutput()
	if err != nil {
		var exitError *exec.ExitError
		if errors.As(err, &exitError) && strings.Contains(string(exitError.Stderr), "No names found") {
			return "", nil
		}

		return "", fmt.Errorf("failed to run command '%s %s': %w\n%s", name, strings.Join(args, " "), err, util.BytesToString(output))
	}

	return strings.TrimSpace(util.BytesToString(output)), nil
}

func executeStep(description string, command string, args ...string) {
	fmt.Println(description)

	cmd := exec.Command(command, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		panic(err)
	}
}

func main() {
	lastTag, err := runCmd("git", "describe", "--tags", "--abbrev=0")
	if err != nil {
		panic(err)
	}

	if lastTag == "" {
		fmt.Println("No tags found.")
	} else {
		fmt.Printf("Latest tag: %s\n", lastTag)
	}

	fmt.Print("Enter new tag: ")

	newTag, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		panic(err)
	}

	newTag = strings.TrimSpace(newTag)

	if newTag == "" {
		panic("No tag entered, aborting.")
	}

	executeStep(fmt.Sprintf("Tagging %s...", newTag), "git", "tag", newTag)
	executeStep(fmt.Sprintf("Pushing tag %s...", newTag), "git", "push", "origin", newTag)

	fmt.Printf("Successfully tagged and pushed %s.\n", newTag)
}
