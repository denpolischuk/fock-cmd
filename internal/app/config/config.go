package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

var (
	// ConfDirPath - user config dir path
	ConfDirPath, _ = os.UserConfigDir()
	// ConfigDirPath - Default config folder path
	ConfigDirPath = filepath.Join(ConfDirPath, "fock")
	confFilePath  = filepath.Join(ConfigDirPath, "conf.json")

	// ErrConfigNotFound - config not found err
	ErrConfigNotFound = fmt.Errorf("Couldn't find config file. Did you run fock init previously?")
)

// GlobalConfig - global config of the cli
type GlobalConfig struct {
	PathToFock string `json:"pathToFock"`
}

// Read - read global config from file
func (c *GlobalConfig) Read() error {
	if c.PathToFock != "" {
		return nil
	}

	confFile, err := os.Open(confFilePath)
	if err != nil {
		return ErrConfigNotFound
	}
	defer confFile.Close()

	d := json.NewDecoder(confFile)
	if err := d.Decode(c); err != nil {
		return err
	}

	return nil
}

// Write - write global config into file
func (c *GlobalConfig) Write() error {
	if _, err := os.Stat(ConfigDirPath); os.IsNotExist(err) {
		os.Mkdir(ConfigDirPath, os.ModePerm)
	}

	file, err := os.OpenFile(confFilePath, os.O_CREATE|os.O_WRONLY, os.ModePerm)
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
