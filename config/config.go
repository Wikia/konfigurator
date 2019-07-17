package config

import (
	"fmt"

	"strings"

	"path"

	"os"

	"regexp"

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
var variableRegex = regexp.MustCompile(`^(?P<type>\w+)\((?P<value>[^)]+)?\)(?:\s*->\s*(?P<destination>\w+))?$`)
var layeredConsulRegex = regexp.MustCompile(`^(?P<key>[^#]+)(?:#(?P<appname>[^@]+)@(?P<environment>\w+))?$`)

func Setup(cmd *cobra.Command) error {
	levels := make([]string, len(log.AllLevels))
	for i, level := range log.AllLevels {
		levels[i] = string(level)
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

	_ = viper.BindPFlag("loglevel", cmd.PersistentFlags().Lookup("logLevel"))
	_ = viper.BindPFlag("kubeconf", cmd.PersistentFlags().Lookup("kubeConf"))
	_ = viper.BindPFlag("vault.address", cmd.PersistentFlags().Lookup("vaultAddress"))
	_ = viper.BindPFlag("vault.token", cmd.PersistentFlags().Lookup("vaultToken"))
	_ = viper.BindPFlag("vault.tokenpath", cmd.PersistentFlags().Lookup("vaultTokenPath"))
	_ = viper.BindPFlag("vault.tlsskipverify", cmd.PersistentFlags().Lookup("vaultTlsSkipVerify"))
	_ = viper.BindPFlag("consul.address", cmd.PersistentFlags().Lookup("consulAddress"))
	_ = viper.BindPFlag("consul.token", cmd.PersistentFlags().Lookup("consulToken"))
	_ = viper.BindPFlag("consul.datacenter", cmd.PersistentFlags().Lookup("consulDatacenter"))
	_ = viper.BindPFlag("consul.tlsskipverify", cmd.PersistentFlags().Lookup("consulTlsSkipVerify"))

	return nil
}

func Get() *Config {
	if currentConfig == nil {
		currentConfig = new(Config)
	}

	return currentConfig
}

func ParseVariableDefinitions(values map[string]string) ([]model.VariableDef, error) {
	var ret []model.VariableDef

	for name, value := range values {
		def := model.NewVariableDef()
		def.Destination = model.CONFIGMAP
		def.Name = strings.TrimSpace(name)
		matches := variableRegex.FindStringSubmatch(strings.TrimSpace(value))

		if matches == nil || len(matches) != 4 {
			return nil, fmt.Errorf("Error parsing variable definition (%s): %s", name, value)
		}

		varType := model.InputType(matches[1])
		switch varType {
		case model.SIMPLE:
			def.Source = varType
			def.Value = matches[2]
		case model.REFERENCE:
			def.Source = varType
			def.Value = matches[2]
		case model.VAULT:
			def.Source = varType
			def.Value = matches[2]
			def.Destination = model.SECRET
		case model.CONSUL:
			def.Source = varType
			def.Value = matches[2]
		case model.LAYERED_CONSUL:
			def.Source = varType
			valueMatches := layeredConsulRegex.FindStringSubmatch(matches[2])
			if len(valueMatches) != 4 {
				return nil, fmt.Errorf("Error parsing layered consul value (%s): %s", name, value)
			}
			def.Value = valueMatches[1]
			def.Context["appname"] = valueMatches[2]
			def.Context["environment"] = valueMatches[3]
		default:
			return nil, fmt.Errorf("Unknown variable source (%s): %s", varType, value)
		}

		if len(matches[3]) != 0 {
			varDestination := model.VariableDestination(matches[3])
			switch varDestination {
			case model.CONFIGMAP:
				def.Destination = varDestination
			case model.SECRET:
				def.Destination = varDestination
			default:
				return nil, fmt.Errorf("Unknown variable type (%s): %s", name, value)
			}
		}

		ret = append(ret, def)
	}

	return ret, nil
}
