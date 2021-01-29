package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
)

var (
	configFilePath = fmt.Sprintf("%s/.config/fock", os.Getenv("HOME"))
	fileName       = configFilePath + "/conf.json"
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

	confFile, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer confFile.Close()

	d := json.NewDecoder(confFile)
	if err := d.Decode(c); err != nil {
		return err
	}

	return err
}

// Write - write global config into file
func (c *GlobalConfig) Write() error {
	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		os.Mkdir(configFilePath, os.ModePerm)
	}

	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY, os.ModePerm)
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
		return "", fmt.Errorf("PathToFock is empty or config is not initialized")
	}

	return strings.TrimRight(c.PathToFock, " /"), nil
}
