package main

import (
	"encoding/json"
	"errors"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

type Configuration struct {
	LibraryRoot string
	Host        string
	Port        int
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
	defer applyConfigEnvars()
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

func applyConfigEnvars() {
	port := os.Getenv("PORT")
	if port != "" {
		p, err := strconv.Atoi(port)
		if err == nil {
			userConfig.Port = p
		}
	}

	host := os.Getenv("HOST")
	if host != "" {
		userConfig.Host = host
	}
}

// Creates default config and writes to file
func createNewDefaultConfig() error {
	config := &Configuration{
		LibraryRoot: filepath.Join(nubayrahDirectory, "library"),
		Host:        "localhost",
		Port:        5050,
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
