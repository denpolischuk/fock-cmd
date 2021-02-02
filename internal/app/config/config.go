package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/denpolischuk/fock-cli/internal/app/shared/types/bookmark"
	"github.com/denpolischuk/fock-cli/internal/app/utils"
)

const packageJSONAppName = `"name": "@redteclab/fock"`

var (
	// ConfDirPath - user config dir path
	ConfDirPath, _ = os.UserConfigDir()
	// ConfigDirPath - Default config folder path
	ConfigDirPath = filepath.Join(ConfDirPath, "fock")
	// ConfFilePath - path to fock cli config file
	ConfFilePath = filepath.Join(ConfigDirPath, "conf.json")

	// ErrConfigNotFound - config not found err
	ErrConfigNotFound = fmt.Errorf("Couldn't find config file. Did you run fock init previously?")
)

// GlobalConfig - global config of the cli
type GlobalConfig struct {
	PathToFock string         `json:"pathToFock"`
	Bookmarks  *bookmark.List `json:"bookmarks"`
}

// Read - read global config from file
func (c *GlobalConfig) Read() error {
	if c.PathToFock != "" {
		return nil
	}

	confFile, err := os.Open(ConfFilePath)
	if err != nil {
		return ErrConfigNotFound
	}
	defer confFile.Close()

	d := json.NewDecoder(confFile)
	if err := d.Decode(c); err != nil {
		return err
	}

	if c.Bookmarks == nil {
		c.Bookmarks = bookmark.NewList()
	}

	return nil
}

func (c *GlobalConfig) beforeWrite() error {
	if err := c.checkFockPath(); err != nil {
		return err
	}
	return nil
}

// Write - write global config into file
func (c *GlobalConfig) Write() error {
	if err := c.beforeWrite(); err != nil {
		return err
	}

	if _, err := os.Stat(ConfigDirPath); os.IsNotExist(err) {
		os.Mkdir(ConfigDirPath, os.ModePerm)
	}

	file, err := os.OpenFile(ConfFilePath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return err
	}
	defer file.Close()
	encoder := json.NewEncoder(file)
	if err := encoder.Encode(c); err != nil {
		return err
	}

	return nil
}

func (c *GlobalConfig) checkFockPath() error {
	p := filepath.Join(path.Clean(c.PathToFock), "package.json")
	if !utils.FileExists(p) {
		return errors.New("fock path is not correct or fock folder misses 'package.json' file")
	}
	f, err := os.OpenFile(p, os.O_RDONLY, 0644)
	if err != nil {
		return err
	}

	contains, err := utils.FileContains(f, packageJSONAppName)
	if err != nil {
		return err
	} else if !contains {
		return errors.New("fock path is not correct or app name in package.json was changed")
	}

	return nil
}

// GetFockPath - returns safe fock path string
func (c *GlobalConfig) GetFockPath() (string, error) {
	if c.PathToFock == "" {
		if err := c.Read(); err != nil {
			return "", err
		}
		if c.PathToFock == "" {
			return "", fmt.Errorf("PathToFock is empty or config is not initialized")
		}
	}

	return filepath.Clean(c.PathToFock), nil
}

// GetNodeModulesBinPath - returns path to fock's node_modules/.bin dir
func (c *GlobalConfig) GetNodeModulesBinPath(bin string) (string, error) {
	fockPath, err := c.GetFockPath()
	if err != nil {
		return "", err
	}

	return filepath.Join(fockPath, "node_modules", ".bin", bin), nil
}
