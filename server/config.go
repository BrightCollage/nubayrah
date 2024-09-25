package main

import (
	"encoding/json"
	"errors"
	"io/fs"
	"log"
	"os"
	"path/filepath"
)

type Configuration struct {
	LibraryRoot string
}

var userConfig *Configuration = nil
var nubayrahDirectory string = ""
var configFilePath string = ""

func init() {
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	nubayrahDirectory = filepath.Join(home, "nubayrah")
	configFilePath = filepath.Join(nubayrahDirectory, "config.json")
}

// Reads user config from ~/.nubayrah
// If this file does not exist a new default config is created there
func OpenConfig() error {
	log.Printf("Opening config at %s", configFilePath)
	configBytes, err := os.ReadFile(configFilePath)
	var perr *fs.PathError
	if err == os.ErrNotExist || errors.As(err, &perr) {
		log.Printf("Config file not found, generating new config")
		return createNewDefaultConfig()
	}

	config := &Configuration{}
	err = json.Unmarshal(configBytes, config)
	if err != nil {
		return err
	}
	userConfig = config
	return err
}

// Creates default config and writes to file
func createNewDefaultConfig() error {
	config := &Configuration{
		LibraryRoot: filepath.Join(nubayrahDirectory, "library"),
	}

	configBytes, err := json.Marshal(config)
	if err != nil {
		return err
	}

	err = os.MkdirAll(nubayrahDirectory, os.ModePerm)
	if err != nil {
		return err
	}
	err = os.WriteFile(configFilePath, configBytes, os.ModePerm)
	if err != nil {
		return err
	}

	userConfig = config
	return nil

}
