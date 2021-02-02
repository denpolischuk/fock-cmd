package server

import (
	"log"
	"os"
	"os/exec"
	"syscall"

	"github.com/denpolischuk/fock-cli/internal/app/config"
	"github.com/denpolischuk/fock-cli/internal/app/consts"
	"github.com/denpolischuk/fock-cli/internal/app/modules"
	"github.com/denpolischuk/fock-cli/internal/app/utils"
	"github.com/kyokomi/emoji"
	"github.com/urfave/cli/v2"
)

func getStatusAction(conf *config.GlobalConfig) modules.ActionGetter {
	return func(c *cli.Context) error {
		if err := conf.Read(); err != nil {
			config.ReadErrorDefaultHandler(err)
		}
		substr, _ := conf.GetNodeModulesBinPath("nodemon")
		if b, p := utils.IsProcessRunning(substr); b {
			emoji.Printf(consts.Emojis["success"]+" Fock node server is running (PID %d)", p.Pid)
		} else {
			emoji.Println(notRunningMessage)
		}
		return nil
	}
}

func getStopAction(conf *config.GlobalConfig) modules.ActionGetter {
	return func(c *cli.Context) error {
		if err := conf.Read(); err != nil {
			config.ReadErrorDefaultHandler(err)
		}
		substr, _ := conf.GetNodeModulesBinPath("nodemon")
		if b, p := utils.IsProcessRunning(substr); b {
			if err := p.SendSignal(syscall.SIGINT); err != nil {
				log.Fatal(err)
			}
			cmd := exec.Command("rm", logOutputFilePath)
			go cmd.Start()
			emoji.Printf(consts.Emojis["success"]+" Fock node server (PID %d) was stopped", p.Pid)
		} else {
			emoji.Println(notRunningMessage)
		}
		return nil
	}
}

func getStartAction(conf *config.GlobalConfig) modules.ActionGetter {
	return func(c *cli.Context) error {
		if err := conf.Read(); err != nil {
			config.ReadErrorDefaultHandler(err)
		}
		substr, _ := conf.GetNodeModulesBinPath("nodemon")
		if b, p := utils.IsProcessRunning(substr); b {
			emoji.Printf(consts.Emojis["think"]+" Fock node server (PID %d) is already running", p.Pid)
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
			cmd.Stderr = os.Stderr
		}

		if err := cmd.Start(); err != nil {
			cmd.Stdout = nil
			cmd.Stderr = nil
			log.Fatal(err)
		}

		if !c.Bool("detached") {
			cmd.Wait()
		} else {
			emoji.Printf(consts.Emojis["success"]+" Fock node server (PID %d) was started", cmd.Process.Pid)
		}

		return nil
	}
}

func getAttachAction(conf *config.GlobalConfig) modules.ActionGetter {
	return func(c *cli.Context) error {
		if err := conf.Read(); err != nil {
			config.ReadErrorDefaultHandler(err)
		}

		substr, _ := conf.GetNodeModulesBinPath("nodemon")
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
