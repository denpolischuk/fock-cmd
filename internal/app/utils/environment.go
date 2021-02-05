package utils

import (
	"fmt"
	"os/exec"
	"regexp"
)

// GetUserShell - returns user shell
func GetUserShell() (string, error) {
	outp, err := exec.Command("bash", "-c", `echo $SHELL | grep -o -P "[A-z]+$"`).Output()
	shell := string(outp)
	if err != nil || shell == "" {
		fmt.Println(err)
		return "", err
	}

	re, rErr := regexp.Compile(`\W+`)
	if rErr != nil {
		return shell, nil
	}

	return re.ReplaceAllString(shell, ""), nil
}
