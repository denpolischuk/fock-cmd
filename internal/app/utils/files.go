package utils

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
)

// FileContains searches string (str) in file (f). Returns true if str was found in f.
func FileContains(f io.Reader, str string) (bool, error) {
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), str) {
			return true, nil
		}
	}
	if err := scanner.Err(); err != nil {
		return false, err
	}

	return false, nil
}

// ReplaceInFile - finds and replaces string in file once
func ReplaceInFile(f io.ReadWriter, target string, replacement string, regexpMode bool) (string, error) {
	scanner := bufio.NewScanner(f)
	var str string
	re, err := regexp.Compile(target)
	if err != nil {
		return "", err
	}
	for scanner.Scan() {
		tStr := scanner.Text()
		if regexpMode {
			if re.MatchString(tStr) {
				tStr = re.ReplaceAllString(tStr, replacement)
			}
		} else {
			if strings.Contains(tStr, target) {
				tStr = strings.Replace(tStr, target, replacement, -1)
			}
		}
		str += fmt.Sprintf("%s\n", tStr)
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	return str, nil
}

// FileExists - returns true if file exists, false otherwise
func FileExists(p string) bool {
	if _, err := os.Stat(p); os.IsNotExist(err) {
		return false
	}
	return true
}
