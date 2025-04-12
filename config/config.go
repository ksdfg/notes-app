package config

import (
	"github.com/spf13/viper"
	"log"
)

// Config holds the configuration for the application
type Config struct {
	LogLevel string `mapstructure:"LOG_LEVEL"`
}

func (c Config) validate() {
	if c.LogLevel == "" {
		log.Fatalln("LOG_LEVEL is required")
	}
}

// Unexported variable to implement singleton pattern
var config *Config

func init() {
	// Setup viper
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")

	// Set default values for config vars
	viper.SetDefault("LOG_LEVEL", "info")

	// Automatically override values in config file with those in environment
	viper.AutomaticEnv()

	// Read config file
	err := viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
		} else {
			// Config file was found but another error was produced
			log.Fatalln(err)
		}
	}

	// Set config object
	err = viper.Unmarshal(&config)
	if err != nil {
		log.Fatalln(err)
	}

	// Validate that all config vars are set
	config.validate()
}

// Get returns the config object
func Get() *Config {
	if config == nil {
		log.Fatalln("Config not set ^._.^")
	}

	return config
}
