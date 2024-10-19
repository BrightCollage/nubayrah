//go:build !docker
// +build !docker

/* For building/running directly on a host OS

Creates a folder ~/.nubayrah
*/

package config

import (
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// Sets default values
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

	libraryRoot := filepath.Join(homeDir, "library")

	viper.SetDefault("library_path", libraryRoot)
	viper.SetDefault("config_path", filepath.Join(homeDir, "config.yaml"))
	viper.SetDefault("host", "0.0.0.0")
	viper.SetDefault("port", 5050)
	viper.SetDefault("db_path", filepath.Join(libraryRoot, "nubayrah.db"))

	// tells Viper to look for `dataRoot/config.yaml``
	viper.AddConfigPath(homeDir)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	return nil
}
