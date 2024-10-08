package config

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
	LibraryRoot   string
	Host          string
	Port          int
	FilePath      string
	HomeDirectory string
}

func New() (*Configuration, error) {
	return OpenConfig()
}

// Reads user config from ~/.nubayrah
// If this file does not exist a new default config is created there
func OpenConfig() (*Configuration, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	config := &Configuration{}

	config.HomeDirectory = filepath.Join(home, ".nubayrah")
	config.FilePath = filepath.Join(config.HomeDirectory, "config.json")

	defer applyConfigEnvars(config)

	log.Printf("Opening config at %s", config.FilePath)
	configBytes, err := os.ReadFile(config.FilePath)
	var perr *fs.PathError
	if err == os.ErrNotExist || errors.As(err, &perr) {
		log.Printf("Config file not found, generating new config")
		return createNewDefaultConfig(config)
	}

	err = json.Unmarshal(configBytes, config)
	if err != nil {
		return config, err
	}

	return config, err
}

func applyConfigEnvars(config *Configuration) {
	port := os.Getenv("PORT")
	if port != "" {
		p, err := strconv.Atoi(port)
		if err == nil {
			config.Port = p
		}
	}

	host := os.Getenv("HOST")
	if host != "" {
		config.Host = host
	}
}

// Creates default config and writes to file
func createNewDefaultConfig(config *Configuration) (*Configuration, error) {

	config.LibraryRoot = filepath.Join(config.HomeDirectory, "library")
	config.Host = "0.0.0.0"
	config.Port = 5050

	configBytes, err := json.Marshal(config)
	if err != nil {
		return nil, err
	}

	err = os.MkdirAll(config.LibraryRoot, os.ModePerm)
	if err != nil {
		return nil, err
	}
	err = os.WriteFile(config.FilePath, configBytes, os.ModePerm)
	if err != nil {
		return nil, err
	}

	return config, err
}
