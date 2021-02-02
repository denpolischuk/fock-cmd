package bookmarks

import (
	"fmt"
	"os"

	"github.com/denpolischuk/fock-cli/internal/app/config"
	"github.com/denpolischuk/fock-cli/internal/app/consts"
	"github.com/denpolischuk/fock-cli/internal/app/modules"
	"github.com/denpolischuk/fock-cli/internal/app/shared/types/bookmark"
	"github.com/denpolischuk/fock-cli/internal/app/utils"
	"github.com/kyokomi/emoji"
	"github.com/urfave/cli/v2"
)

const defaultListOutputLimit int16 = 15

func getAddBookmarkAction(conf *config.GlobalConfig) modules.ActionGetter {
	return func(c *cli.Context) error {
		if err := conf.Read(); err != nil {
			config.ReadErrorDefaultHandler(err)
		}
		if c.Args().Len() != 2 {
			emoji.Printf("%s wrong amount of arguments was passed. See usage `fock bm add -h`", consts.Emojis["fail"])
			os.Exit(1)
		}

		alias, URL := c.Args().Get(0), c.Args().Get(1)

		if conf.Bookmarks.List[alias] != "" {
			resp := utils.PromptUserYesOrNo(fmt.Sprintf("Bookamark with alias `%s` already exists. Do you want to overwrite it? (N/y)", alias))
			if resp != "y" {
				emoji.Printf("%s no changes were made", consts.Emojis["fail"])
				os.Exit(0)
			}
		}

		if err := conf.Bookmarks.Add(alias, URL); err != nil {
			emoji.Println(consts.Emojis["fail"] + err.Error())
			os.Exit(1)
		}

		if err := conf.Write(); err != nil {
			config.WriteErrorDefaultHandler(err)
		}

		emoji.Printf("%s %s -> %s succesfully added.", consts.Emojis["success"], alias, URL)

		return nil
	}
}

func getListBookmarksAction(conf *config.GlobalConfig) modules.ActionGetter {
	return func(c *cli.Context) error {
		if err := conf.Read(); err != nil {
			config.ReadErrorDefaultHandler(err)
		}
		if conf.Bookmarks == nil || conf.Bookmarks.List == nil || len(conf.Bookmarks.List) == 0 {
			emoji.Printf("%s you don't have any bookmarks yet...\n", consts.Emojis["think"])
			os.Exit(0)
		}

		limit := defaultListOutputLimit

		if c.Bool("all") {
			limit = bookmark.BookmarksCap
		}

		var i int16 = 0
		for alias, URL := range conf.Bookmarks.List {
			if i == limit {
				return nil
			}

			fmt.Printf("%s -> %s\n", alias, URL)
			i++
		}

		return nil
	}
}

func getRemoveBookmarksAction(conf *config.GlobalConfig) modules.ActionGetter {
	return func(c *cli.Context) error {
		if err := conf.Read(); err != nil {
			config.ReadErrorDefaultHandler(err)
		}

		if conf.Bookmarks == nil || conf.Bookmarks.List == nil || len(conf.Bookmarks.List) == 0 {
			emoji.Printf("%s you don't have any bookmarks yet...\n", consts.Emojis["think"])
			os.Exit(0)
		}

		if c.Bool("all") {
			conf.Bookmarks.List = make(map[string]string)
			if err := conf.Write(); err != nil {
				config.WriteErrorDefaultHandler(err)
			}
			fmt.Println("Bookmarks list was successfully erased.")

			return nil
		}

		alias := c.Args().First()

		if c.Args().Len() != 1 || alias == "" {
			emoji.Printf("%s you have to specify the name of the bookmark you want to remove.", consts.Emojis["fail"])
			os.Exit(1)
		}

		if conf.Bookmarks.List[alias] == "" {
			emoji.Printf("%s bookmark with alias `%s` doesn't exist.", consts.Emojis["fail"], alias)
			os.Exit(0)
		}

		URL := conf.Bookmarks.List[alias]
		delete(conf.Bookmarks.List, alias)

		if err := conf.Write(); err != nil {
			config.WriteErrorDefaultHandler(err)
		}

		emoji.Printf("%s %s -> %s succesfully removed.", consts.Emojis["success"], alias, URL)

		return nil
	}
}
