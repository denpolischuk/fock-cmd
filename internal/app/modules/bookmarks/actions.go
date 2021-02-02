package bookmarks

import (
	"fmt"
	"os"

	"github.com/denpolischuk/fock-cli/internal/app/config"
	"github.com/denpolischuk/fock-cli/internal/app/consts"
	"github.com/denpolischuk/fock-cli/internal/app/modules"
	"github.com/denpolischuk/fock-cli/internal/app/utils"
	"github.com/kyokomi/emoji"
	"github.com/urfave/cli/v2"
)

func getAddBookmarkAction(conf *config.GlobalConfig) modules.ActionGetter {
	return func(c *cli.Context) error {
		if err := conf.Read(); err != nil {
			config.ReadErrorDefaultHandler(err)
		}
		if c.Args().Len() != 2 {
			emoji.Println(consts.Emojis["fail"] + " wrong amount of arguments was passed. See usage `fock bm add -h`")
			os.Exit(1)
		}

		alias, URL := c.Args().Get(0), c.Args().Get(1)

		if conf.Bookmarks.List[alias] != "" {
			resp := utils.PromptUserYesOrNo(fmt.Sprintf("Bookamark with alias `%s` already exists. Do you want to overwrite it? (N/y)", alias))
			if resp != "y" {
				emoji.Println(consts.Emojis["fail"] + " no changes were made")
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
