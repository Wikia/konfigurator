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
Application:
  Name: my_app
  Namespace: staging
  Definitions:
    # This value will be inserted directly into configuration (config map destination is default for simple type)
    simple_var: simple(some value)
    # This simple secret will be inserted into secrets as it is
    simple_secret: simple(abracadabra)->secret
    # This value will be fetched from the configured Vault server under path "/secret/app/temp" under key "test" (secret is also default for vault type)
    vault_secret: vault(/secret/app/temp:test)
    # This value will be fetched from configured Consul server from the KV path "config/base/dev/DATACENTER"
    consul_var: consul(config/base/dev/DATACENTER)
    # This value references internal k8s variables available inside POD
    reference_var: simple(spec.nodeName)->reference
    # This is Wikia-specific hierarchical configuration in Consul - it will try to fetch values from three different localtions in Consul (ordered):
    # config/sample_app/development/some_key, config/sample_app/base/some_key. config/base/development/some_key
    layered_var: layered_consul(some_key#smaple_app@development)
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

### get
This command will fetch all the variables from the defined sources and output them on STDOUT as multi document YAML file.
 
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

When outputting in `envrc` format, values with type `reference` will be omitted.

#### options
```
  -d, --destinationFolder string   Where to store the output files (default "/Users/harnas/_Projects_/_golang_/src/github.com/Wikia/konfigurator")
  -h, --help                       help for download
      --name string                Name of the service to download variables for
  -n, --namespace string           Kubernetes namespace for which files should be generated for (default "dev")
  -o, --output string              Output format (available formats: [envrc k8s-yaml]) (default "k8s-yaml")
```

### set
This command will update k8s POD definition with the configured variables and secrets inserting references to proper ConfigMap and Secret.
All variables will be injected as environment variables with names the same as variable name.

#### options

```
  -m, --configMap string         File where ConfigMap definitions are stored
  -t, --containerName string     Name of the container to modify in deployment
  -f, --deployment string        Deployment file with configuration that should be updated
  -d, --destinationFile string   Destination file where to write updated deployment configuration
  -h, --help                     help for update
  -w, --overwrite                Should configuration definitions be completely replaced by the new one or just appended
  -s, --secrets string           File where Secrets are stored
  -y, --yes                      Answer all questions 'yes' - no confirmations and interaction
```

### merge
This command will take deployment file and insert in-line values for config (non-secrets) as env variables and reference
any secrets specified in the configuration. It will put in specified folder resulting files (updated deployment and secrets file).

#### options

```
  -t, --containerName string    Name of the container to modify in deployment
  -f, --deployment string       Deployment file with configuration that should be updated
  -d, --destinationDir string   Destination where to write resulting filesn (default ".")
  -h, --help                    help for merge
  -n, --namespace string        Kubernetes namespace for which files should be generated for (default "dev")
  -w, --overwrite               Should configuration definitions be completely replaced by the new one or just appended
  -s, --secretName string       Name of the secret to use in the deployment mappings (defaults to 'containerName')
```
