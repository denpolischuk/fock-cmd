package initmodule

import (
	"fmt"
	"strings"

	"github.com/denpolischuk/fock-cmd/internal/app/config"
	"github.com/denpolischuk/fock-cmd/internal/app/modules"
	"github.com/urfave/cli/v2"
)

// Init - initialization module
type Init struct {
	Command *cli.Command
	Config  *config.GlobalConfig
}

func getInitAction(conf *config.GlobalConfig) modules.ActionGetter {
	return func(c *cli.Context) error {
		if err := conf.Read(); err == nil { // If no error was returned then config file already exists
			fmt.Print(fmt.Sprintf("You already have configured fock path (%s), do you want to rewrite the config? (N/y): ", conf.PathToFock))
			var userInput string
			fmt.Scanln(&userInput)
			if strings.Trim(strings.ToLower(userInput), " %n") != "y" {
				return nil
			}
		}

		conf.PathToFock = c.String("path")
		if err := conf.Write(); err != nil {
			return err
		}
		fmt.Println("New config successfully created.")

		return nil
	}
}

// New - creating new Init module
func New(conf *config.GlobalConfig) (*Init, error) {
	return &Init{
		Command: &cli.Command{
			Name:  "init",
			Usage: "config CLI tool to work with your Fock instance",
			Flags: []cli.Flag{
				&cli.StringFlag{Name: "path", Aliases: []string{"p"}, Usage: "Path to Fock root folder"},
			},
			Action: getInitAction(conf),
		},
		Config: conf,
	}, nil
}

// GetCommand - returns command of Init module
func (i *Init) GetCommand() (*cli.Command, error) {
	if i.Command == nil {
		return nil, fmt.Errorf("Init module doesn't have command initialized")
	}

	return i.Command, nil
}
