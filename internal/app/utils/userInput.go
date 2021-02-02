package utils

import (
	"fmt"
	"regexp"
	"strings"
)

var re, _ = regexp.Compile(`\W+`)

// PromptUserYesOrNo ...
func PromptUserYesOrNo(promptMessage string) string {
	fmt.Print(promptMessage)
	var userInput string
	fmt.Scanln(&userInput)
	return strings.ToLower(re.ReplaceAllString(userInput, ""))
}
