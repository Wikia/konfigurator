package inputs

import (
	"fmt"

	"io/ioutil"

	"crypto/tls"
	"net/http"
	"path/filepath"

	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/Wikia/konfigurator/config"
	"github.com/Wikia/konfigurator/model"
	"github.com/hashicorp/vault/api"
)

type Vault struct {
	client *api.Client
}

func (v *Vault) intiClient() error {
	config := config.Get()

	tr := &http.Transport{}
	if config.Vault.TLSSkipVerify {
		tr = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}

	clientCfg := api.Config{
		Address:    config.Vault.Address,
		HttpClient: &http.Client{Transport: tr},
	}

	var err error
	v.client, err = api.NewClient(&clientCfg)

	if err != nil {
		return err
	}

	// authenticating
	token := config.Vault.Token
	if len(token) == 0 {
		log.Debug("VaultInput: initial token empty - trying to read from file")
		tokenFile, err := filepath.Abs(config.Vault.TokenPath)
		if err != nil {
			return fmt.Errorf("VaultInput: could not get absolute path for token file: %s", err)
		}
		contents, err := ioutil.ReadFile(tokenFile)
		if err != nil {
			log.WithError(err).WithField("path", tokenFile).Warn("VaultInput: could not read contents of a vault token file")
		} else {
			token = string(contents)
		}
	}

	v.client.SetToken(token)
	v.client.Auth()
}

func (v *Vault) Fetch(variable model.VariableDef) (*model.Variable, error) {
	if variable.Source != model.VAULT {
		return nil, fmt.Errorf("Invalid variable type: %s for %s", variable.Type, variable.Name)
	}

	if v.client == nil {
		v.intiClient()
	}

	source := strings.SplitN(variable.Value.(string), ":", 2)

	if len(source) != 2 {
		return nil, fmt.Errorf("VaultInput: variable has incorect value syntax: %s", variable.Name)
	}

	secret, err := v.client.Logical().Read(source[0])

	if err != nil {
		return nil, err
	}

	log.WithFields(log.Fields{
		"variable": variable.Name,
		"path":     variable.Value,
		"secret":   secret.Data[source[1]],
	}).Debug("Read variable from vault")

	ret := model.Variable{
		Name:  variable.Name,
		Type:  variable.Type,
		Value: secret.Data[source[1]],
	}

	return &ret, nil
}

func init() {
	Register(model.VAULT, &Vault{})
}
