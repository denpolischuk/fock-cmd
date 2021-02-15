package nginx

import (
	"fmt"

	"github.com/denpolischuk/fock-cli/internal/app/config"
	"github.com/denpolischuk/fock-cli/internal/app/consts"
	"github.com/urfave/cli/v2"
)

var (
	varnishHostFlag = &cli.StringFlag{Name: "varnish-host", Usage: "Varnish host", Required: false, Value: consts.DockerHost}
	varnishPortFlag = &cli.StringFlag{Name: "varnish-port", Usage: "Varnish port", Required: false, Value: consts.VarnishPort}
	portMapFlag     = &cli.StringFlag{Name: "port", Aliases: []string{"p"}, Usage: "Port mapping", Required: false, Value: "80:80"}
	detachedFlag    = &cli.BoolFlag{Name: "detached", Aliases: []string{"d"}, Usage: "Detached mode", Value: false}
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
						varnishHostFlag,
						varnishPortFlag,
					},
				},
				{
					Name:      "run",
					Usage:     "runs previously built docker image of nginx preview.",
					UsageText: "fock nginx run - runs nginx preview container.",
					Action:    getRunAction(conf),
					Flags: []cli.Flag{
						portMapFlag,
						detachedFlag,
					},
				},
				{
					Name:      "stop",
					Usage:     "stops running nginx preview container.",
					UsageText: "fock nginx stop - stops nginx preview container.",
					Action:    getStopAction(conf),
				},
				{
					Name:      "start",
					Usage:     "builds docker image of nginx preview and runs it immidiately away.",
					UsageText: "fock nginx start - builds and runs nginx preview image.",
					Action:    getStartAction(conf),
					Flags: []cli.Flag{
						varnishHostFlag,
						varnishPortFlag,
						portMapFlag,
						detachedFlag,
					},
				},
				{
					Name:      "status",
					Usage:     "checks if fock preview nginx container is running",
					UsageText: "fock nginx status - checks if fock preview nginx container is running.",
					Action:    getStatusAction(conf),
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
