package config

import (
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type Configuration struct {
	LibraryRoot   string
	Host          string
	Port          int
	FilePath      string
	HomeDirectory string
}

func New() Configuration {

	err := setDefaultConfig()
	if err != nil {
		log.Printf("Error when trying to get home directory: %v", err)
	}

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found
			log.Printf("Configuration file not found, creating a new one at: ~ %v", viper.Get("file_path"))
			if err = viper.SafeWriteConfig(); err != nil {
				log.Printf("Error when writing config: %v", err)
				panic(err)
			}
		} else {
			// Config file found but error
			panic(err)
		}
	}

	log.Printf("Using configuration file %v", viper.Get("file_path"))

	return Configuration{
		LibraryRoot:   viper.GetString("library_root"),
		Host:          viper.GetString("host"),
		Port:          viper.GetInt("port"),
		FilePath:      viper.GetString("file_path"),
		HomeDirectory: viper.GetString("home_directory"),
	}
}

// Sets Default values
func setDefaultConfig() error {

	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	homeDir := filepath.Join(home, ".nubayrah")

	viper.SetDefault("library_path", filepath.Join(homeDir, "library"))
	viper.SetDefault("file_path", filepath.Join(homeDir, "config.yaml"))
	viper.SetDefault("host", "0.0.0.0")
	viper.SetDefault("port", 5050)
	viper.SetDefault("home_directory", homeDir)

	viper.SetConfigName("config") // name of config file (without extension)
	viper.SetConfigType("yaml")   // REQUIRED if the config file does not have the extension in the name

	viper.AddConfigPath(homeDir) // path to look for the config file in

	return nil
}
