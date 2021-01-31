package initmodule

import (
	"fmt"

	"github.com/denpolischuk/fock-cli/internal/app/config"
	"github.com/urfave/cli/v2"
)

// Init - initialization module
type Init struct {
	Command *cli.Command
	Config  *config.GlobalConfig
}

// New - creating new Init module
func New(conf *config.GlobalConfig) (*Init, error) {
	return &Init{
		Command: &cli.Command{
			Name:  "init",
			Usage: "config CLI tool to work with your Fock instance",
			Flags: []cli.Flag{
				&cli.StringFlag{Name: "path", Aliases: []string{"p"}, Usage: "Path to Fock root folder", Required: true},
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
