package initmodule

import (
	"fmt"
	"io/ioutil"
	"os/exec"
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
		if err := ioutil.WriteFile(config.ConfigFilePath+"/zsh_autocompletion", b[:len(b)-1], 0644); err != nil {
			fmt.Println(err)
			return
		}
		zshPath := utils.PromptPathToResource("[Autocompletion]: Provide the path to .zshrc?", consts.PathToZshRc)

		// TODO write to file using Go instead of bash and check if this script is already in file
		err := exec.Command("bash", "-c", fmt.Sprintf(`printf "\n%s %s/zsh_autocompletion\n" >> %s`, ZshRcScript, config.ConfigFilePath, zshPath)).Run()
		if err != nil {
			fmt.Println(err)
		}
		emoji.Println(autocompletionInstalledMessage)
		break
	case "bash":
		b := []byte(BashAutocompletionScript)
		if err := ioutil.WriteFile("/etc/bash_completion.d/"+consts.AppBinName, b[:len(b)-1], 0644); err != nil {
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
		if err := conf.Read(); err == nil { // If no error was returned then config file already exists
			fmt.Print(fmt.Sprintf("You already have configured fock path (%s), do you want to rewrite the config? (N/y): ", conf.PathToFock))
			var userInput string
			fmt.Scanln(&userInput)
			if strings.Trim(strings.ToLower(userInput), " %n") != "y" {
				setupAutocompletion(conf)
				return nil
			}
		} else if err != config.ErrConfigNotFound {
			return err
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
