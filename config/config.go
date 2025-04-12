package config

import (
	"errors"
	"github.com/spf13/viper"
	"log"
)

// Config holds the configuration for the application
type Config struct {
	/*
		Logging configuration
	*/

	// LogLevel defines what is the minimum level of logs that will be logged, defaults to info
	LogLevel string `mapstructure:"LOG_LEVEL"`

	/*
	   Fiber App configuration
	*/

	// Port is the port that the server will listen on, defaults to 3000
	Port int `mapstructure:"PORT"`

	/*
	   DB configuration
	*/

	// DBHost is the host of the database server
	DBHost string `mapstructure:"DB_HOST"`
	// DBPort is the port of the database server
	DBPort int `mapstructure:"DB_PORT"`
	// DBUser is the user for authenticating to the database server
	DBUser string `mapstructure:"DB_USER"`
	// DBPassword is the password for authenticating to the database server
	DBPassword string `mapstructure:"DB_PASSWORD"`
	// DBName is the name of the database to use for the app
	DBName string `mapstructure:"DB_NAME"`
	// TestDBName is the name of the database to use for testing
	TestDBName string `mapstructure:"TEST_DB_NAME"`
	// DBSSLMode is the SSL mode to use for connecting to the database, defaults to disable
	DBSSLMode string `mapstructure:"DB_SSL_MODE"`
}

// validate checks if the required configuration fields are set and logs a fatal error if any are missing.
func (c Config) validate() {
	if c.DBName == "" {
		panic("DB_NAME must be set")
	}

	if c.DBHost == "" {
		panic("DB_HOST must be set")
	}

	if c.DBUser == "" {
		panic("DB_USER must be set")
	}

	if c.DBPassword == "" {
		panic("DB_PASS must be set")
	}

	if c.DBPort == 0 {
		panic("DB_PORT must be set")
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
	viper.SetDefault("PORT", 3000)
	viper.SetDefault("DB_SSL_MODE", "disable")

	// Automatically override values in config file with those in environment
	viper.AutomaticEnv()

	// Read config file
	err := viper.ReadInConfig()
	if err != nil {
		if errors.As(err, &viper.ConfigFileNotFoundError{}) {
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
