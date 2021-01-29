package server

import (
	"fmt"
	"log"
	"strings"
	"syscall"

	"github.com/denpolischuk/fock-cmd/internal/app/config"
	"github.com/denpolischuk/fock-cmd/internal/app/modules"
	"github.com/denpolischuk/fock-cmd/internal/app/utils"
	"github.com/kyokomi/emoji"
	"github.com/urfave/cli/v2"
)

// Server - fock node server module
type Server struct {
	Command *cli.Command
	Config  *config.GlobalConfig
}

func getStatusAction(conf *config.GlobalConfig) modules.ActionGetter {
	return func(c *cli.Context) error {
		conf.Read()
		if b, p := utils.IsProcessRunning(strings.TrimRight(conf.PathToFock, " /")); b {
			ppid, _ := p.Ppid()
			emoji.Printf("\U00002705 Fock node server is running (PID %d)", ppid)
		} else {
			emoji.Println("\U0000274C  Fock server is not running")
		}
		return nil
	}
}

func getStopAction(conf *config.GlobalConfig) modules.ActionGetter {
	return func(c *cli.Context) error {
		conf.Read()
		if b, p := utils.IsProcessRunning(conf.PathToFock + "/node_modules/.bin/nodemon"); b {
			ppid, _ := p.Ppid()
			fmt.Println(p.Cmdline())
			if err := p.SendSignal(syscall.SIGINT); err != nil {
				log.Fatal(err)
			}
			emoji.Printf("\U00002705 Fock node server (PID %d) was stopped", ppid)
		} else {
			emoji.Println("\U0000274C  Fock server is not running")
		}
		return nil
	}
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
			},
		},
		Config: conf,
	}, nil
}

// GetCommand - returns command of Init module
func (s *Server) GetCommand() (*cli.Command, error) {
	if s.Command == nil {
		return nil, fmt.Errorf("Init module doesn't have command initialized")
	}

	return s.Command, nil
}
