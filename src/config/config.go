package config

import (
	"github.com/spf13/viper"
	"log"
)

type OpalConfig struct {
	Id                string
	Url               string
	Token             string
	AdvertisedAddress string
}

type ToggleConfig struct {
	Key         string
	UsersPolicy struct {
		Source  string
		Package string
		Rule    string
	}
	Spec map[string]interface{}
}

type TargetConfig struct {
	TargetType string
	TargetSpec map[string]interface{}
}

// OpTogglesConfig Values to be loaded from configuration file, keys are case-insensitive
type OpTogglesConfig struct {
	Bind    string
	Sources []OpalConfig
	Target  TargetConfig
	Toggles []ToggleConfig
}

var GlobalConfig = OpTogglesConfig{
	Bind: ":8080", // Default value
}

func init() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/etc/optoggles/")
	viper.AddConfigPath(".") // Useful for development

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalln("failed reading configuration file: " + err.Error())
	}
	// Config file found and successfully parsed

	// TODO: Add validation
	if err := viper.Unmarshal(&GlobalConfig); err != nil {
		log.Fatalln("invalid configuration file: " + err.Error())
	}
}
