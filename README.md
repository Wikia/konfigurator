# konfigurator 
[![Build Status](https://travis-ci.org/Wikia/konfigurator.svg?branch=master)](https://travis-ci.org/Wikia/konfigurator)
[![Coverage Status](https://coveralls.io/repos/github/Wikia/konfigurator/badge.svg?branch=master)](https://coveralls.io/github/Wikia/konfigurator?branch=master)

## Sample configuration

```yaml
LogLevel: debug
Consul:
  Address: consul.service.poz-dev.consul:8500
  Datacenter: poz-dev
  TlsSkipVerify: true
Vault:
  Address: https://active.vault.service.poz-dev.consul:8200
  TlsSkipVerify: true
Definitions:
  # This value will be inserted directly into configuration
  - name: Simple Variable
    type: config
    source: simple
    value: some value
  # This simple secret will be inserted into secrets as it is
  - name: Simple Secret
    type: secret
    source: simple
    value: abracadabra
  # This value will be fetched from the configured Vault server under path "/sercret/app/temp" under key "test"
  - name: SecretVault
    type: secret
    source: vault
    value: /secret/app/temp:test
  # This value will be fetched from configured Consul server from the KV path "config/base/dev/DATACENTER"
  - name: ConsulValue
    type: config
    source: consul
    value: config/base/dev/DATACENTER
  # This value refences internal k8s variables available inside POD
  - name: ReferencedValue
    type: reference
    source: simple
    value: spec.nodeName
```

## Global configuration flags
```
      --config string             config file (default is $HOME/.konfigurator.yaml)
      --consulAddress string      Address to a Consul server (default "consul.service.consul")
      --consulDatacenter string   Datacenter to be used in Consul
      --consulTlsSkipVerify       Should TLS certificate be verified
      --consulToken string        Token to be used when authenticating with Consul
  -h, --help                      help for konfigurator
      --kubeConf string           Path to a kubeconf config file
      --logLevel string           What type of logs should be emited (available: panic, fatal, error, warning, info, debug) (default "info")
      --vaultAddress string       Address to a Vault server
      --vaultTlsSkipVerify        Should TLS certificate be verified
      --vaultToken string         Token to be used when authenticating with Vault (overrides vaultTokenPath)
      --vaultTokenPath string     Path to a file with Vault token (default "$HOME/.vault-token")
```

## Available commands

### download
This command will fetch all the variables from the defined sources and put them in proper file(s).
 
#### Available types:
* **config** - values will be put into ConfigMaps
* **secret** - values will be encoded and put into Secrets
* **reference** - values will be put into Deployment as reference to other POD variables

#### Available sources:
* **simple** - values are stored statically in the configuration file
* **vault** - values are fetched from the Vault server (you will need proper token to authorize with the server)
* **consul** - values are fetched from the Consul's KV store

#### Available output formats:
* **k8s-yaml** - will save configuration into Secret and ConfigMap YAMLs for use with kubectl
* **envrc** - will save all configuration into shell compatible file for use in local development or testing

When outputing `envrc` values with type `reference` will be omitted.

#### options
```
  -d, --destinationFolder string   Where to store the output files (default "YOUR WORKING DIR")
  -h, --help                       help for download
  -o, --output string              Output format (available formats: [envrc k8s-yaml]) (default "k8s-yaml")
  -s, --serviceName string         What is the service name which settings will be downloaded as
```
### update
This command will update k8s POD definition with the configured variables and secrets inserting references to proper ConfigMap and Secret.
All variables will be injected as environment variables with names the same as variable name.

#### options

```
  -m, --configMap string         File where ConfigMap definitions are stored
  -t, --containerName string     Name of the container to modify in deployment
  -f, --deployment string        Deployment file where configuration should be updated
  -d, --destinationFile string   Destination file where to write deployment
  -h, --help                     help for update
  -w, --overwrite                Should configuration definitions be completely replaced by the new one or just appended
  -s, --secrets string           File where Secrets are stored
  -y, --yes                      Answer all questions 'yes' - no confirmations and interaction
```
