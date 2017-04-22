package config

import (
	"os"

	"fmt"

	"os/user"
	"strings"

	"path"

	log "github.com/Sirupsen/logrus"
	"github.com/Wikia/konfigurator/model"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type Config struct {
	LogLevel     string
	KubeConfPath string
	Input        InputConfig
	Definitions  []model.VariableDef
}

type InputConfig struct {
	VaultAddress       string
	VaultToken         string
	VaultTokenPath     string
	VaultTLSSkipVerify bool
	ConsulAddress      string
}

var currentConfig *Config

func Setup(cmd *cobra.Command) error {
	cfg := Get()

	levels := make([]string, len(log.AllLevels))
	for i, level := range log.AllLevels {
		levels[i] = fmt.Sprintf("%s", level)
	}

	usr, err := user.Current()
	var homeDir string
	if err != nil {
		log.WithError(err).Warn("Could not get current user - using current working dir as home")
		homeDir, _ = os.Getwd()
	} else {
		homeDir = usr.HomeDir
	}

	cmd.PersistentFlags().StringVar(&cfg.KubeConfPath, "kubeConf", "", "Path to a kubeconf config file")
	cmd.PersistentFlags().StringVar(&cfg.Input.VaultAddress, "vaultAddress", "", "Address to a Vault server")
	cmd.PersistentFlags().StringVar(&cfg.Input.VaultToken, "vaultToken", "", "Token to be used when authenticating with Vault (overrides vaultTokenPath)")
	cmd.PersistentFlags().StringVar(&cfg.Input.VaultTokenPath, "vaultTokenPath", path.Join(homeDir, ".vault-token"), "Path to a file with Vault token")
	cmd.PersistentFlags().BoolVar(&cfg.Input.VaultTLSSkipVerify, "vaultTlsSkipVerify", false, "Should TLS certificate be verified")
	cmd.PersistentFlags().StringVar(&cfg.Input.ConsulAddress, "consulAddress", "consul.service.consul", "Address to a Consul server")
	cmd.PersistentFlags().StringVar(&cfg.LogLevel, "logLevel", "info", fmt.Sprintf("What type of logs should be emited (available: %s)", strings.Join(levels, ", ")))
	return nil
}

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
