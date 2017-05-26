package config

import (
	"fmt"

	"strings"

	"path"

	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/Wikia/konfigurator/model"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	viper "github.com/spf13/viper"
)

type Config struct {
	LogLevel     string
	KubeConfPath string
	Vault        VaultConfig
	Consul       ConsulConfig
	Application  model.Application
}

type VaultConfig struct {
	Address       string
	Token         string
	TokenPath     string
	TLSSkipVerify bool
}

type ConsulConfig struct {
	Address       string
	Datacenter    string
	Token         string
	TLSSkipVerify bool
}

var currentConfig *Config

func Setup(cmd *cobra.Command) error {
	levels := make([]string, len(log.AllLevels))
	for i, level := range log.AllLevels {
		levels[i] = fmt.Sprintf("%s", level)
	}

	homeDir, err := homedir.Dir()
	if err != nil {
		log.WithError(err).Warn("Error getting user home dir - using current working dir")
		homeDir, _ = os.Getwd()
	}

	tokenDir := path.Join(homeDir, ".vault-token")
	viper.SetDefault("loglevel", "info")
	viper.SetDefault("vault.tokenpath", tokenDir)

	cmd.PersistentFlags().String("logLevel", "info", fmt.Sprintf("What type of logs should be emited (available: %s)", strings.Join(levels, ", ")))
	cmd.PersistentFlags().String("kubeConf", "", "Path to a kubeconf config file")
	cmd.PersistentFlags().String("vaultAddress", "", "Address to a Vault server")
	cmd.PersistentFlags().String("vaultToken", "", "Token to be used when authenticating with Vault (overrides vaultTokenPath)")
	cmd.PersistentFlags().String("vaultTokenPath", tokenDir, "Path to a file with Vault token")
	cmd.PersistentFlags().Bool("vaultTlsSkipVerify", false, "Should TLS certificate be verified")
	cmd.PersistentFlags().String("consulAddress", "consul.service.consul", "Address to a Consul server")
	cmd.PersistentFlags().String("consulToken", "", "Token to be used when authenticating with Consul")
	cmd.PersistentFlags().String("consulDatacenter", "", "Datacenter to be used in Consul")
	cmd.PersistentFlags().Bool("consulTlsSkipVerify", false, "Should TLS certificate be verified")

	viper.BindPFlag("loglevel", cmd.PersistentFlags().Lookup("logLevel"))
	viper.BindPFlag("kubeconf", cmd.PersistentFlags().Lookup("kubeConf"))
	viper.BindPFlag("vault.address", cmd.PersistentFlags().Lookup("vaultAddress"))
	viper.BindPFlag("vault.token", cmd.PersistentFlags().Lookup("vaultToken"))
	viper.BindPFlag("vault.tokenpath", cmd.PersistentFlags().Lookup("vaultTokenPath"))
	viper.BindPFlag("vault.tlsskipverify", cmd.PersistentFlags().Lookup("vaultTlsSkipVerify"))
	viper.BindPFlag("consul.address", cmd.PersistentFlags().Lookup("consulAddress"))
	viper.BindPFlag("consul.token", cmd.PersistentFlags().Lookup("consulToken"))
	viper.BindPFlag("consul.datacenter", cmd.PersistentFlags().Lookup("consulDatacenter"))
	viper.BindPFlag("consul.tlsskipverify", cmd.PersistentFlags().Lookup("consulTlsSkipVerify"))

	return nil
}

func Get() *Config {
	if currentConfig == nil {
		currentConfig = new(Config)
	}

	return currentConfig
}
