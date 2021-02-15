package nginx

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/denpolischuk/fock-cli/internal/app/config"
	"github.com/denpolischuk/fock-cli/internal/app/consts"
	"github.com/denpolischuk/fock-cli/internal/app/modules"
	"github.com/denpolischuk/fock-cli/internal/app/utils"
	"github.com/kyokomi/emoji"
	"github.com/urfave/cli/v2"
)

func beforeBuild(filepath string, vHost string, vPort string) error {
	f, err := os.OpenFile(filepath, os.O_RDWR, os.ModePerm)
	if err != nil {
		return err
	}
	defer f.Close()

	proxyTarget := vHost
	if len(vPort) > 0 {
		proxyTarget = fmt.Sprintf("%s:%s", proxyTarget, vPort)
	}
	r, err := utils.ReplaceInFile(f, `proxy_target \"[0-9\.:]+\"`, fmt.Sprintf(`proxy_target "%s"`, proxyTarget), true)
	if err != nil {
		return err
	}

	// Clear file and write new config
	f.Truncate(0)
	f.Seek(0, 0)
	w := bufio.NewWriter(f)
	_, err = w.WriteString(r)
	if err != nil {
		return err
	}
	err = w.Flush()
	if err != nil {
		return err
	}

	return nil
}

func getInitAction(conf *config.GlobalConfig) modules.ActionGetter {
	return func(c *cli.Context) error {
		if err := conf.Read(); err != nil {
			config.ReadErrorDefaultHandler(err)
		}

		if !utils.CheckIfAppInstalled("docker") {
			return emoji.Errorf("%s you need to have docker installed to make this module work.", consts.Emojis["fail"])
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

func getBuildAction(conf *config.GlobalConfig) modules.ActionGetter {
	return func(c *cli.Context) error {
		if err := conf.Read(); err != nil {
			config.ReadErrorDefaultHandler(err)
		}

		p, err := conf.GetNginxPath()
		if err != nil {
			return err
		}

		host, port := c.String("varnish-host"), c.String("varnish-port")

		if err := beforeBuild(filepath.Join(p, "common", "server-commons-include.conf"), host, port); err != nil {
			return err
		}

		return nil
	}
}
