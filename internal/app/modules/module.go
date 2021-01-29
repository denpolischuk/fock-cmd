package modules

import "github.com/urfave/cli/v2"

// Module - describes module behavior
type Module interface {
	GetCommand() (*cli.Command, error)
}

// ActionGetter - describes the function the returnes action
type ActionGetter = func(c *cli.Context) error
