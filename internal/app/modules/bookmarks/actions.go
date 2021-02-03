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
	"github.com/pkg/browser"
	"github.com/urfave/cli/v2"
)

const defaultListOutputLimit int = 15

func checkBookmarks(conf *config.GlobalConfig) {
	if conf.Bookmarks == nil || conf.Bookmarks.List == nil || len(conf.Bookmarks.List) == 0 {
		emoji.Printf("%s you don't have any bookmarks yet...\n", consts.Emojis["think"])
		os.Exit(0)
	}

	return
}

func checkIfBookmarkExists(conf *config.GlobalConfig, alias string) {
	if conf.Bookmarks.List[alias] == "" {
		emoji.Printf("%s bookmark with alias `%s` doesn't exist.", consts.Emojis["fail"], alias)
		os.Exit(0)
	}

	return
}

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
		checkBookmarks(conf)

		limit := defaultListOutputLimit

		if c.Bool("all") {
			limit = bookmark.BookmarksCap
		}

		i := 0
		fmt.Println("Here goes your bookmark list: ")
		for alias, URL := range conf.Bookmarks.List {
			fmt.Printf("%s -> %s\n", alias, URL)
			i++
			if i == limit {
				return nil
			}
		}

		if i < len(conf.Bookmarks.List)-1 {
			fmt.Printf("And %d more...", len(conf.Bookmarks.List)-i-1)
		}

		return nil
	}
}

func getRemoveBookmarksAction(conf *config.GlobalConfig) modules.ActionGetter {
	return func(c *cli.Context) error {
		if err := conf.Read(); err != nil {
			config.ReadErrorDefaultHandler(err)
		}

		checkBookmarks(conf)

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
		checkIfBookmarkExists(conf, alias)

		URL := conf.Bookmarks.List[alias]
		delete(conf.Bookmarks.List, alias)

		if err := conf.Write(); err != nil {
			config.WriteErrorDefaultHandler(err)
		}

		emoji.Printf("%s %s -> %s succesfully removed.", consts.Emojis["success"], alias, URL)

		return nil
	}
}

func getOpenBookmarkAction(conf *config.GlobalConfig) modules.ActionGetter {
	return func(c *cli.Context) error {
		if err := conf.Read(); err != nil {
			config.ReadErrorDefaultHandler(err)
		}

		checkBookmarks(conf)

		alias := c.Args().First()
		if c.Args().Len() != 1 || alias == "" {
			emoji.Printf("%s you have to specify exactly one name of the bookmark you want to open.", consts.Emojis["fail"])
			os.Exit(1)
		}

		checkIfBookmarkExists(conf, alias)

		browser.OpenURL(conf.Bookmarks.List[alias])

		return nil
	}
}
