package config

import (
	"fmt"

	"strings"

	"path"

	log "github.com/Sirupsen/logrus"
	"github.com/Wikia/konfigurator/model"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
)

type Config struct {
	LogLevel     string
	KubeConfPath string
	Vault        VaultConfig
	Consul       ConsulConfig
	Definitions  []model.VariableDef
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

	cmd.PersistentFlags().String("kubeConf", "", "Path to a kubeconf config file")
	cmd.PersistentFlags().String("vaultAddress", "", "Address to a Vault server")
	cmd.PersistentFlags().String("vaultToken", "", "Token to be used when authenticating with Vault (overrides vaultTokenPath)")
	cmd.PersistentFlags().String("vaultTokenPath", path.Join(homedir.Dir(), ".vault-token"), "Path to a file with Vault token")
	cmd.PersistentFlags().Bool("vaultTlsSkipVerify", false, "Should TLS certificate be verified")
	cmd.PersistentFlags().String("consulAddress", "consul.service.consul", "Address to a Consul server")
	cmd.PersistentFlags().String("consulToken", "", "Token to be used when authenticating with Consul")
	cmd.PersistentFlags().String("consulDatacenter", "", "Datacenter to be used in Consul")
	cmd.PersistentFlags().Bool("consulTlsSkipVerify", false, "Should TLS certificate be verified")
	cmd.PersistentFlags().String("logLevel", "info", fmt.Sprintf("What type of logs should be emited (available: %s)", strings.Join(levels, ", ")))

	return nil
}

func Get() *Config {
	if currentConfig == nil {
		currentConfig = new(Config)
	}

	return currentConfig
}
