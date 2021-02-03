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
					Usage:   "list bookmarks.",
					Flags: []cli.Flag{
						&cli.BoolFlag{
							Name:    "all",
							Aliases: []string{"a"},
							Usage:   "shows all bookmarks at once",
						},
					},
					ArgsUsage: "<alias> - alias of the bookmark URL.",
					Action:    getListBookmarksAction(conf),
				},
				{
					Name:      "open",
					Aliases:   []string{"o"},
					Usage:     "open bookmark.",
					ArgsUsage: "<alias> - alias of the bookmark URL.",
					Action:    getOpenBookmarkAction(conf),
					BashComplete: func(c *cli.Context) {
						if err := conf.Read(); err != nil {
							return
						}
						if conf.Bookmarks == nil || len(conf.Bookmarks.List) == 0 {
							return
						}
						for _, bm := range conf.Bookmarks.List {
							fmt.Println(bm)
						}
					},
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
