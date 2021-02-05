package bookmarks

import (
	"fmt"

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

func checkBookmarks(conf *config.GlobalConfig) error {
	if conf.Bookmarks == nil || conf.Bookmarks.List == nil || len(conf.Bookmarks.List) == 0 {
		return emoji.Errorf("%s you don't have any bookmarks yet...\n", consts.Emojis["think"])
	}

	return nil
}

func checkIfBookmarkExists(conf *config.GlobalConfig, alias string) error {
	if conf.Bookmarks.List[alias] == "" {
		return emoji.Errorf("%s bookmark with alias `%s` doesn't exist.", consts.Emojis["fail"], alias)
	}

	return nil
}

func getAddBookmarkAction(conf *config.GlobalConfig) modules.ActionGetter {
	return func(c *cli.Context) error {
		if err := conf.Read(); err != nil {
			config.ReadErrorDefaultHandler(err)
		}
		if c.Args().Len() != 2 {
			return emoji.Errorf("%s wrong amount of arguments was passed. See usage `fock bm add -h`", consts.Emojis["fail"])
		}

		alias, URL := c.Args().Get(0), c.Args().Get(1)

		if conf.Bookmarks.List[alias] != "" {
			resp := utils.PromptUserYesOrNo(fmt.Sprintf("Bookamark with alias `%s` already exists. Do you want to overwrite it? (N/y)", alias))
			if resp != "y" {
				return emoji.Errorf("%s no changes were made", consts.Emojis["fail"])
			}
		}

		if err := conf.Bookmarks.Add(alias, URL); err != nil {
			return emoji.Errorf(consts.Emojis["fail"] + err.Error())
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

		if err := checkBookmarks(conf); err != nil {
			return err
		}

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

		if err := checkBookmarks(conf); err != nil {
			return err
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
			return emoji.Errorf("%s you have to specify the name of the bookmark you want to remove.", consts.Emojis["fail"])
		}

		if err := checkIfBookmarkExists(conf, alias); err != nil {
			return err
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

func getOpenBookmarkAction(conf *config.GlobalConfig) modules.ActionGetter {
	return func(c *cli.Context) error {
		if err := conf.Read(); err != nil {
			config.ReadErrorDefaultHandler(err)
		}

		if err := checkBookmarks(conf); err != nil {
			return err
		}

		alias := c.Args().First()
		if c.Args().Len() != 1 || alias == "" {
			return emoji.Errorf("%s you have to specify exactly one name of the bookmark you want to open.", consts.Emojis["fail"])
		}

		if err := checkIfBookmarkExists(conf, alias); err != nil {
			return err
		}

		browser.OpenURL(conf.Bookmarks.List[alias])

		return nil
	}
}
