package server

import (
	"log"
	"os"
	"os/exec"
	"syscall"

	"github.com/denpolischuk/fock-cmd/internal/app/config"
	"github.com/denpolischuk/fock-cmd/internal/app/modules"
	"github.com/denpolischuk/fock-cmd/internal/app/utils"
	"github.com/kyokomi/emoji"
	"github.com/urfave/cli/v2"
)

const (
	logOutputFilePath = "/tmp/fock_server_output"
	notRunningMessage = "\U0000274C Fock server is not running"
)

func getNodeModulesBinPath(conf *config.GlobalConfig, bin string) (string, error) {
	if fockPath, err := conf.GetFockPath(); err != nil {
		return "", err
	} else {
		return fockPath + "/node_modules/.bin/" + bin, nil
	}
}

func getStatusAction(conf *config.GlobalConfig) modules.ActionGetter {
	return func(c *cli.Context) error {
		if err := conf.Read(); err != nil {
			return err
		}
		substr, _ := getNodeModulesBinPath(conf, "nodemon")
		if b, p := utils.IsProcessRunning(substr); b {
			ppid, _ := p.Ppid()
			emoji.Printf("\U00002705 Fock node server is running (PID %d)", ppid)
		} else {
			emoji.Println(notRunningMessage)
		}
		return nil
	}
}

func getStopAction(conf *config.GlobalConfig) modules.ActionGetter {
	return func(c *cli.Context) error {
		if err := conf.Read(); err != nil {
			return err
		}
		substr, _ := getNodeModulesBinPath(conf, "nodemon")
		if b, p := utils.IsProcessRunning(substr); b {
			ppid, _ := p.Ppid()
			if err := p.SendSignal(syscall.SIGINT); err != nil {
				log.Fatal(err)
			}
			cmd := exec.Command("rm", logOutputFilePath)
			go cmd.Start()
			emoji.Printf("\U00002705 Fock node server (PID %d) was stopped", ppid)
		} else {
			emoji.Println(notRunningMessage)
		}
		return nil
	}
}

func getStartAction(conf *config.GlobalConfig) modules.ActionGetter {
	return func(c *cli.Context) error {
		if err := conf.Read(); err != nil {
			return err
		}
		substr, _ := getNodeModulesBinPath(conf, "nodemon")
		if b, p := utils.IsProcessRunning(substr); b {
			ppid, _ := p.Ppid()
			emoji.Printf("\U00002806 Fock node server (PID %d) is already running", ppid)
			return nil
		}

		fockPath, _ := conf.GetFockPath()
		cmdOut1, cmdOut2 := "", ""
		if c.Bool("detached") {
			cmdOut1, cmdOut2 = ">", logOutputFilePath
		}

		cmd := exec.Command("yarn", "--cwd", fockPath, "dev-server", cmdOut1, cmdOut2)

		if !c.Bool("detached") {
			cmd.Stdout = os.Stdout
		}

		if err := cmd.Start(); err != nil {
			cmd.Stdout = nil
			log.Fatal(err)
		}

		if !c.Bool("detached") {
			cmd.Wait()
		} else {
			emoji.Printf("\U00002705 Fock node server (PID %d) was started", cmd.Process.Pid)
		}

		return nil
	}
}

func getAttachAction(conf *config.GlobalConfig) modules.ActionGetter {
	return func(c *cli.Context) error {
		if err := conf.Read(); err != nil {
			return err
		}

		substr, _ := getNodeModulesBinPath(conf, "nodemon")
		if b, _ := utils.IsProcessRunning(substr); !b {
			emoji.Println(notRunningMessage)
			return nil
		}

		cmd := exec.Command("tail", "-f", logOutputFilePath)
		cmd.Stdout = os.Stdout

		if err := cmd.Start(); err != nil {
			cmd.Stdout = nil
			log.Fatal(err)
		}

		cmd.Wait()

		return nil
	}
}
