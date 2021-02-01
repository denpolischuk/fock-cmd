package utils

import (
	"bufio"
	"fmt"
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
func FileContains(f *os.File, str string) bool {
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), str) {
			return true
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}

	return false
}

// FileExists - returns true if file exists, false otherwise
func FileExists(p string) bool {
	if _, err := os.Stat(p); os.IsNotExist(err) {
		return false
	}
	return true
}
