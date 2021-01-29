package server

import (
	"fmt"
	"log"
	"strings"
	"syscall"

	"github.com/denpolischuk/fock-cmd/internal/app/config"
	"github.com/denpolischuk/fock-cmd/internal/app/modules"
	"github.com/denpolischuk/fock-cmd/internal/app/utils"
	"github.com/kyokomi/emoji"
	"github.com/urfave/cli/v2"
)

func getStatusAction(conf *config.GlobalConfig) modules.ActionGetter {
	return func(c *cli.Context) error {
		conf.Read()
		if b, p := utils.IsProcessRunning(strings.TrimRight(conf.PathToFock, " /")); b {
			ppid, _ := p.Ppid()
			emoji.Printf("\U00002705 Fock node server is running (PID %d)", ppid)
		} else {
			emoji.Println("\U0000274C  Fock server is not running")
		}
		return nil
	}
}

func getStopAction(conf *config.GlobalConfig) modules.ActionGetter {
	return func(c *cli.Context) error {
		conf.Read()
		if b, p := utils.IsProcessRunning(conf.PathToFock + "/node_modules/.bin/nodemon"); b {
			ppid, _ := p.Ppid()
			fmt.Println(p.Cmdline())
			if err := p.SendSignal(syscall.SIGINT); err != nil {
				log.Fatal(err)
			}
			emoji.Printf("\U00002705 Fock node server (PID %d) was stopped", ppid)
		} else {
			emoji.Println("\U0000274C  Fock server is not running")
		}
		return nil
	}
}
