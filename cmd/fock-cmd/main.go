package main

import (
	"log"
	"os"

	"github.com/denpolischuk/fock-cmd/internal/app/config"
	"github.com/denpolischuk/fock-cmd/internal/app/modules"
	"github.com/denpolischuk/fock-cmd/internal/app/modules/initmodule"
	"github.com/denpolischuk/fock-cmd/internal/app/modules/server"
	"github.com/urfave/cli/v2"
)

// LoadModules - loads all cli modules in project
func LoadModules(c *config.GlobalConfig) (*[]modules.Module, error) {
	initModule, err := initmodule.New(c)
	if err != nil {
		return nil, err
	}

	serverModule, err := server.New(c)
	if err != nil {
		return nil, err
	}

	modules := []modules.Module{
		initModule,
		serverModule,
	}

	return &modules, nil
}

func main() {
	conf := config.GlobalConfig{}

	modules, err := LoadModules(&conf)
	if err != nil {
		log.Fatal(err)
	}

	commands := make([]*cli.Command, len(*modules))

	for i, m := range *modules {
		commands[i], err = m.GetCommand()
		if err != nil {
			log.Fatal(err)
		}
	}

	app := &cli.App{
		Name:     "Fock CLI",
		Usage:    "This is the Fock CLI tool designed to help you working with fock instace.",
		Commands: commands,
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
