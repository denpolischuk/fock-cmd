package consts

import "os"

var (
	// PathToZshRc - default path to zshrc
	PathToZshRc = os.Getenv("HOME") + "/.zshrc"

	// Emojis - Map of emojis for stdout
	Emojis = map[string]string{
		"success": "\U00002705",
		"fail":    "\U0000274C",
		"dead":    "\U00002620",
		"think":   "\U0001F914",
	}
)

const (
	// AppBinName - binary file name
	AppBinName = "fock"
)
