package utils

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
)

var re, _ = regexp.Compile(`\W+`)

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

// PromptUserYesOrNo ...
func PromptUserYesOrNo(promptMessage string) string {
	fmt.Print(promptMessage)
	var userInput string
	fmt.Scanln(&userInput)
	return strings.ToLower(re.ReplaceAllString(userInput, ""))
}
