package config

import (
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

func GetConfig() error {

	err := setDefaultConfig()
	if err != nil {
		log.Printf("Error when trying to get home directory: %v", err)
	}

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found
			log.Printf("Configuration file not found, creating a new one at: ~ %v", viper.Get("config_path"))
			if err = viper.SafeWriteConfig(); err != nil {
				return err
			}
		} else {
			// Config file found but error
			return err
		}
	}

	log.Printf("Using configuration file %v", viper.Get("config_path"))
	return nil
}

// Sets Default values
func setDefaultConfig() error {

	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	// Create ~/.nubayrah if not exist.
	homeDir := filepath.Join(home, ".nubayrah")
	err = os.MkdirAll(homeDir, os.ModePerm)
	if err != nil {
		return err
	}

	viper.SetDefault("library_path", filepath.Join(homeDir, "library"))
	viper.SetDefault("config_path", filepath.Join(homeDir, "config.yaml"))
	viper.SetDefault("host", "0.0.0.0")
	viper.SetDefault("port", 5050)
	viper.SetDefault("home_directory", homeDir)
	viper.SetDefault("db_path", filepath.Join(homeDir, "nubayrah.db"))

	viper.SetConfigName("config") // name of config file (without extension)
	viper.SetConfigType("yaml")   // REQUIRED if the config file does not have the extension in the name

	viper.AddConfigPath(homeDir) // path to look for the config file in

	return nil
}
