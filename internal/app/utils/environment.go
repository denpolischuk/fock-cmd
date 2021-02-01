package utils

import (
	"fmt"
	"os/exec"
	"regexp"
)

// GetUserShell - returns user shell
func GetUserShell() (string, error) {
	outp, err := exec.Command("bash", "-c", `echo $SHELL | grep -o -P "[A-z]+$"`).Output()
	if err != nil {
		fmt.Println(err)
		return "", nil
	}

	shell := string(outp)

	re, rErr := regexp.Compile(`\W+`)
	if rErr != nil {
		return shell, nil
	}

	return re.ReplaceAllString(shell, ""), nil
}
