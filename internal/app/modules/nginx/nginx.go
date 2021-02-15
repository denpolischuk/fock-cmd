package nginx

import (
	"fmt"

	"github.com/denpolischuk/fock-cli/internal/app/config"
	"github.com/denpolischuk/fock-cli/internal/app/consts"
	"github.com/urfave/cli/v2"
)

// Nginx - module to work with preview nginx
type Nginx struct {
	Command *cli.Command
	Config  *config.GlobalConfig
}

// New ...
func New(conf *config.GlobalConfig) (*Nginx, error) {
	return &Nginx{
		Command: &cli.Command{
			Name:  "nginx",
			Usage: "tools to work with preview nginx",
			Subcommands: []*cli.Command{
				{
					Name:      "init",
					Usage:     "setup fock CLI to work with nginx. This will require having nginx preview cloned to your machine.",
					UsageText: "fock nginx init [path] - setup fock CLI to work with nginx. Path is optional unless you are not in nginx folder.",
					ArgsUsage: "<path> - path to preview nginx folder",
					Action:    getInitAction(conf),
				},
				{
					Name:      "build",
					Usage:     "builds docker image of nginx preview with correct varnish host and port.",
					UsageText: "fock nginx build - builds nginx preview image.",
					Action:    getBuildAction(conf),
					Flags: []cli.Flag{
						&cli.StringFlag{Name: "varnish-host", Usage: "Varnish host", Required: false, Value: consts.DockerHost},
						&cli.StringFlag{Name: "varnish-port", Usage: "Varnish port", Required: false, Value: consts.VarnishPort},
					},
				},
			},
		},
		Config: conf,
	}, nil
}

// GetCommand - returns command of Bookmarks module
func (ng *Nginx) GetCommand() (*cli.Command, error) {
	if ng.Command == nil {
		return nil, fmt.Errorf("Nginx module doesn't have command initialized")
	}

	return ng.Command, nil
}
