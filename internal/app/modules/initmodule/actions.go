package initmodule

import (
	"fmt"
	"strings"

	"github.com/denpolischuk/fock-cmd/internal/app/config"
	"github.com/denpolischuk/fock-cmd/internal/app/modules"
	"github.com/urfave/cli/v2"
)

func getInitAction(conf *config.GlobalConfig) modules.ActionGetter {
	return func(c *cli.Context) error {
		if err := conf.Read(); err == nil { // If no error was returned then config file already exists
			fmt.Print(fmt.Sprintf("You already have configured fock path (%s), do you want to rewrite the config? (N/y): ", conf.PathToFock))
			var userInput string
			fmt.Scanln(&userInput)
			if strings.Trim(strings.ToLower(userInput), " %n") != "y" {
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

		return nil
	}
}
