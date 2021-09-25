package config

import (
	"github.com/spf13/viper"
	"log"
	"optoggles/types"
	"os"
)

// OpTogglesConfig Values to be loaded from configuration file, keys are case-insensitive
type OpTogglesConfig struct {
	OPA struct {
		Address string
		JWT string
	}
	OPAL struct {
		Address string
		JWT string
	}
	TogglesTarget struct {
		TargetType string
		// TODO: Replace with generic map that decodes per target type
		LdAddress string
		LdJWT string
	}
	Toggles []types.Toggle
}


var GlobalConfig OpTogglesConfig

func init() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/etc/optoggles/") // TODO: Enable overriding it with a cmdline variable
	viper.AddConfigPath(".") // TODO: This is for debug, leave it?

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
		} else {
			// Config file was found but another error was produced
			panic("Can't read config file: " + err.Error())
		}
	}
	// Config file found and successfully parsed

	if err := viper.Unmarshal(&GlobalConfig); err != nil {
		panic("Can't unmarshal configuration: " + err.Error())
	}

	log.Printf(os.Getwd())
	log.Println("Loaded configuration file: ", GlobalConfig)
}

