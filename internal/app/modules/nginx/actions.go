package nginx

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

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
	r, err := utils.ReplaceInFile(f, `proxy_target\s+\"[localhost0-9\.:]+\"`, fmt.Sprintf(`proxy_target "%s"`, proxyTarget), true)
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

func buildImage(p string) error {
	dockerCommand := fmt.Sprintf("docker build -t %s %s", DefaultImageName, p)

	cmd := exec.Command("bash", "-c", dockerCommand)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	cmd.Stdout = nil
	cmd.Stderr = nil
	if err != nil {
		return err
	}

	return nil
}

func runImage(portMap string, detached bool) error {
	d := ""
	if detached {
		d = "-d"
	}
	dockerCommand := fmt.Sprintf("docker run -p %s %s %s", portMap, d, DefaultImageName)

	cmd := exec.Command("bash", "-c", dockerCommand)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	cmd.Stdout = nil
	cmd.Stderr = nil
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
		return buildImage(p)
	}
}

func getRunAction(conf *config.GlobalConfig) modules.ActionGetter {
	return func(c *cli.Context) error {
		if err := conf.Read(); err != nil {
			config.ReadErrorDefaultHandler(err)
		}

		_, err := conf.GetNginxPath()
		if err != nil {
			return err
		}

		d := c.Bool("detached")
		p := c.String("port")

		if err := runImage(p, d); err != nil {
			return err
		}

		return nil
	}
}

func getStopAction(conf *config.GlobalConfig) modules.ActionGetter {
	return func(c *cli.Context) error {
		if err := conf.Read(); err != nil {
			config.ReadErrorDefaultHandler(err)
		}

		_, err := conf.GetNginxPath()
		if err != nil {
			return err
		}

		runningContainers, err := exec.Command("bash", "-c", fmt.Sprintf(`docker ps -q --filter="ancestor=%s"`, DefaultImageName)).Output()
		if err != nil {
			return err
		}
		if len(runningContainers) > 0 {
			stopCommand := fmt.Sprintf(`docker stop %s`, runningContainers)
			err = exec.Command("bash", "-c", stopCommand).Run()
			if err != nil {
				return err
			}
			emoji.Printf("%s %s container stopped\n", consts.Emojis["success"], strings.Trim(strings.Join(strings.Split(string(runningContainers), "\n"), ", "), ", \n"))
		} else {
			emoji.Printf("%s there are no running containers\n", consts.Emojis["think"])
		}

		return nil
	}
}

func getStartAction(conf *config.GlobalConfig) modules.ActionGetter {
	return func(c *cli.Context) error {
		if err := conf.Read(); err != nil {
			config.ReadErrorDefaultHandler(err)
		}

		p, err := conf.GetNginxPath()
		if err != nil {
			return err
		}

		host, port, pMap, d := c.String("varnish-host"), c.String("varnish-port"), c.String("port"), c.Bool("detached")

		if err := beforeBuild(filepath.Join(p, "common", "server-commons-include.conf"), host, port); err != nil {
			return err
		}
		if err := buildImage(p); err != nil {
			return err
		}

		if err := runImage(pMap, d); err != nil {
			return err
		}

		return nil
	}
}
