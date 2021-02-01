package consts

import (
	"os"
	"path/filepath"
)

var (
	//HomeDir - user's home dir
	HomeDir, _ = os.UserHomeDir()
	// PathToZshRc - default path to zshrc
	PathToZshRc = filepath.Join(HomeDir, ".zshrc")

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
