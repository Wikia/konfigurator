package config

import (
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/Wikia/konfigurator/model"
	"github.com/spf13/viper"
)

type Config struct {
	KubeConfPath string
	Definitions  []VariableDef
}

type VariableDef struct {
	Name   string
	Source VariableSource
	Type   model.VariableType
	Value  interface{}
}

var currentConfig *Config

func Get() *Config {
	if currentConfig == nil {
		currentConfig = new(Config)

		if err := viper.Unmarshal(&currentConfig); err != nil {
			log.WithError(err).Error("Error parsing config file")
			os.Exit(-5)
		}
	}

	return currentConfig
}
