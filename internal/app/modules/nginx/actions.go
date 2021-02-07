package nginx

import (
	"fmt"
	"os"
	"path"

	"github.com/denpolischuk/fock-cli/internal/app/config"
	"github.com/denpolischuk/fock-cli/internal/app/consts"
	"github.com/denpolischuk/fock-cli/internal/app/modules"
	"github.com/denpolischuk/fock-cli/internal/app/utils"
	"github.com/kyokomi/emoji"
	"github.com/urfave/cli/v2"
)

func getInitAction(conf *config.GlobalConfig) modules.ActionGetter {
	return func(c *cli.Context) error {
		if err := conf.Read(); err != nil {
			config.ReadErrorDefaultHandler(err)
		}

		var p string
		if c.Args().Len() > 1 {
			return emoji.Errorf("%s wrong amount of arguments was passed. You have to be in nginx folder or provide path to it as an argument.", consts.Emojis["fail"])
		} else if c.Args().First() != "" {
			p = path.Clean(c.Args().First())
		} else {
			var err error
			p, err = os.Getwd()
			if err != nil {
				return err
			}
		}

		if conf.PathToNginx != "" {
			resp := utils.PromptUserYesOrNo("You already have configured path to nginx preview, do you want to rewrite it? (N/y): ")
			if resp != "y" {
				fmt.Println("Alright... Let's pretend that never happened.")
				return nil
			}
		}

		conf.PathToNginx = p

		if err := conf.Write(); err != nil {
			return err
		}
		emoji.Printf("%s nginx initialized succesfully.\n", consts.Emojis["success"])

		return nil
	}
}
