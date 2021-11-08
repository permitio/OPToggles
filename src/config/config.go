package config

import (
	"github.com/spf13/viper"
	"log"
	"os"
)

// OpTogglesConfig Values to be loaded from configuration file, keys are case-insensitive

type OpalConfig struct {
	Id                string
	Url               string
	Token             string
	AdvertisedAddress string
}

type ToggleConfig struct {
	Key           string
	UsersDocument struct {
		Source  string
		Package string
		Rule    string
	}
	Spec map[string]interface{}
}

type TargetConfig struct {
	TargetType string
	// TODO: Replace with generic map that decodes per target type
	TargetSpec map[string]interface{}
}

type OpTogglesConfig struct {
	Sources []OpalConfig
	Target  TargetConfig
	Toggles []ToggleConfig
}

var GlobalConfig OpTogglesConfig

func init() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/etc/optoggles/") // TODO: Enable overriding it with a cmdline variable
	viper.AddConfigPath(".")               // TODO: This is for debug, leave it?

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
		} else {
			// Config file was found but another error was produced
			panic("Can't read config file: " + err.Error())
		}
	}
	// Config file found and successfully parsed

	// TODO: Add validation
	if err := viper.Unmarshal(&GlobalConfig); err != nil {
		panic("Can't unmarshal configuration: " + err.Error())
	}

	log.Printf(os.Getwd())
	log.Println("Loaded configuration file: ", GlobalConfig)
}
