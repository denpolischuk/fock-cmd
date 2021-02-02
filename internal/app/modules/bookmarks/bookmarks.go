package bookmarks

import (
	"fmt"

	"github.com/denpolischuk/fock-cli/internal/app/config"
	"github.com/urfave/cli/v2"
)

// Bookmarks - module
type Bookmarks struct {
	Command *cli.Command
	Config  *config.GlobalConfig
}

// New - creating new Bookmarks module
func New(conf *config.GlobalConfig) (*Bookmarks, error) {
	return &Bookmarks{
		Command: &cli.Command{
			Name:    "bookmarks",
			Usage:   "opens stored URL in your browser by its alias.",
			Aliases: []string{"bm", "bms"},
			Subcommands: []*cli.Command{
				{
					Name:      "add",
					Aliases:   []string{"create"},
					Usage:     "adds new bookmark.",
					UsageText: BookmarksAddUsage,
					ArgsUsage: "<alias> - url alias for easy access, <URL> - url that will be stored under alias",
					Action:    getAddBookmarkAction(conf),
				},
				{
					Name:        "remove",
					Aliases:     []string{"rm"},
					Description: "removes bookmark by it's alias.",
					Usage:       "fock bm rm [alias].",
					ArgsUsage:   "[alias] - alias of the bookmark you want to remove.",
					Action:      getRemoveBookmarksAction(conf),
					Flags: []cli.Flag{
						&cli.BoolFlag{
							Name:    "all",
							Aliases: []string{"a"},
							Usage:   "removes all bookmarks.",
						},
					},
				},
				{
					Name:    "list",
					Aliases: []string{"l"},
					Usage:   "list bookmark.",
					Flags: []cli.Flag{
						&cli.BoolFlag{
							Name:    "all",
							Aliases: []string{"a"},
							Usage:   "shows all flags at once",
						},
					},
					ArgsUsage: "<alias> - url alias for easy access, <URL> - url that will be stored under alias",
					Action:    getListBookmarksAction(conf),
				},
			},
		},
		Config: conf,
	}, nil
}

// GetCommand - returns command of Bookmarks module
func (i *Bookmarks) GetCommand() (*cli.Command, error) {
	if i.Command == nil {
		return nil, fmt.Errorf("Init module doesn't have command initialized")
	}

	return i.Command, nil
}
