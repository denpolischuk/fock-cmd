package utils

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// PromptPathToResource - Asks user path to zsh or uses default one if user input is empty
func PromptPathToResource(promptStr string, def string) string {
	path := def
	var userInput string
	fmt.Printf("%s [%s]: ", promptStr, path)
	userInput = ""
	fmt.Scanln(&userInput)
	if userInput != "" {
		path = userInput
	}

	return filepath.Clean(path)
}

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
