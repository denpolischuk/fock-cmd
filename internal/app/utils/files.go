package utils

import (
	"bufio"
	"io"
	"os"
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

// FileExists - returns true if file exists, false otherwise
func FileExists(p string) bool {
	if _, err := os.Stat(p); os.IsNotExist(err) {
		return false
	}
	return true
}
