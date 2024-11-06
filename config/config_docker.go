//go:build docker
// +build docker

/*
For building/running with docker with `go build -tags docker`

Sets library_root to `/library` and config dir to `/data` for easier
mounting in docker.
*/

package config

import (
	"path/filepath"

	"github.com/spf13/viper"
)

func setDefaultConfig() error {
	const dataRoot = "/data"
	const libraryRoot = "/library"

	viper.SetDefault("library_path", libraryRoot)
	viper.SetDefault("config_path", filepath.Join(dataRoot, "config.yaml"))
	viper.SetDefault("host", "0.0.0.0")
	viper.SetDefault("port", 5050)
	viper.SetDefault("db_path", filepath.Join(libraryRoot, "nubayrah.db"))

	// tells Viper to look for `dataRoot/config.yaml``
	viper.AddConfigPath(dataRoot)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	return nil
}
