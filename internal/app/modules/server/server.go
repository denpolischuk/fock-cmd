package server

import (
	"fmt"

	"github.com/denpolischuk/fock-cli/internal/app/config"
	"github.com/urfave/cli/v2"
)

// Server - fock node server module
type Server struct {
	Command *cli.Command
	Config  *config.GlobalConfig
}

// New - creating new Init module
func New(conf *config.GlobalConfig) (*Server, error) {
	return &Server{
		Command: &cli.Command{
			Name:  "server",
			Usage: "interacts with fock node server",
			Subcommands: []*cli.Command{
				{
					Name:   "status",
					Usage:  "checks if server is running",
					Action: getStatusAction(conf),
				},
				{
					Name:   "stop",
					Usage:  "stops server if it's running",
					Action: getStopAction(conf),
				},
				{
					Name:  "start",
					Usage: "starts fock node server in watch mode",
					Flags: []cli.Flag{
						&cli.BoolFlag{Name: "detached", Aliases: []string{"d"}, Usage: "Run fock node server in detached mode"},
					},
					Action: getStartAction(conf),
				},
				{
					Name:   "attach",
					Usage:  "attach to already running fock node server's stdout",
					Action: getAttachAction(conf),
				},
			},
		},
		Config: conf,
	}, nil
}

// GetCommand - returns command of Init module
func (s *Server) GetCommand() (*cli.Command, error) {
	if s.Command == nil {
		return nil, fmt.Errorf("Server module doesn't have command initialized")
	}

	return s.Command, nil
}
