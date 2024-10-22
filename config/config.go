package config

import (
	"log"

	"github.com/spf13/viper"
)

func Load() error {
	err := setDefaultConfig()
	if err != nil {
		log.Printf("Error when trying to get home directory: %v", err)
		return err
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
