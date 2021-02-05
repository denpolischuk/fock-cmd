package initmodule

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/denpolischuk/fock-cli/internal/app/config"
	"github.com/denpolischuk/fock-cli/internal/app/consts"
	"github.com/denpolischuk/fock-cli/internal/app/modules"
	"github.com/denpolischuk/fock-cli/internal/app/utils"
	"github.com/kyokomi/emoji"
	"github.com/urfave/cli/v2"
)

var (
	re, _ = regexp.Compile(`\W+`)
)

const autocompletionInstalledMessage = "[Autocompletion]: Successfuly installed. You will need to restart your shell session."

func setupAutocompletion(conf *config.GlobalConfig) {
	fmt.Print("[Autocompletion]: Do you want to setup autocompletion? (Y/n): ")
	var userInput string
	fmt.Scanln(&userInput)
	if strings.ToLower(re.ReplaceAllString(userInput, "")) == "n" {
		fmt.Println("[Autocompletion]: Setup skipped... Run init again if you change your mind.")
		return
	}

	shell, err := utils.GetUserShell()
	if err != nil {
		emoji.Println(defaultShellDetectErrorMessage)
		return
	}

	switch shell {
	case "zsh":
		b := []byte(ZshAutocompletionScript)
		autocompletionPath := filepath.Join(config.ConfigDirPath, "zsh_autocompletion")

		// If autocompletion script doesn't exists then create it
		if !utils.FileExists(autocompletionPath) {
			if err := ioutil.WriteFile(autocompletionPath, b[:len(b)-1], 0644); err != nil {
				fmt.Println(err)
				return
			}
		}

		zshPath := utils.PromptPathToResource("[Autocompletion]: Provide the path to .zshrc?", consts.PathToZshRc)
		zshFile, err := os.OpenFile(zshPath, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer zshFile.Close()

		autocompletionCommand := fmt.Sprintf("%s %s", ZshRcScript, autocompletionPath)

		conains, err := utils.FileContains(zshFile, autocompletionCommand)
		if err != nil {
			fmt.Println(err)
			return
		} else if !conains {
			if _, err := zshFile.WriteString(fmt.Sprintf("\n#Fock CLI autocompletion\n%s\n", autocompletionCommand)); err != nil {
				fmt.Println(err)
				return
			}
		}

		emoji.Println(autocompletionInstalledMessage)

		break
	case "bash":
		b := []byte(BashAutocompletionScript)
		bashCompletionPath := filepath.Join("etc", "bash_completion.d", consts.AppBinName)
		if err := ioutil.WriteFile(bashCompletionPath, b[:len(b)-1], 0644); err != nil {
			fmt.Println(err)
			return
		}
		emoji.Println(autocompletionInstalledMessage)

		break
	default:
		emoji.Println(defaultShellDetectErrorMessage)
	}

	return
}

func getInitAction(conf *config.GlobalConfig) modules.ActionGetter {
	return func(c *cli.Context) error {
		conf.Read()
		if utils.FileExists(config.ConfFilePath) {
			resp := utils.PromptUserYesOrNo(fmt.Sprintf("You already have configured fock CLI, do you want to rewrite configs? (N/y): "))
			if resp != "y" {
				setupAutocompletion(conf)
				return nil
			}
		}
		conf.PathToFock = c.String("path")
		if err := conf.Write(); err != nil {
			return err
		}
		fmt.Println("New config successfully created.")
		setupAutocompletion(conf)

		return nil
	}
}
