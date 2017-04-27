package inputs

import (
	"fmt"

	"crypto/tls"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/Wikia/konfigurator/config"
	"github.com/Wikia/konfigurator/model"
	"github.com/hashicorp/consul/api"
)

type Consul struct {
	client    *api.Client
	queryOpts *api.QueryOptions
}

func (c *Consul) initClient() error {
	config := config.Get()

	tr := &http.Transport{}
	if config.Consul.TLSSkipVerify {
		tr = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}

	clientCfg := api.Config{
		Address:    config.Consul.Address,
		HttpClient: &http.Client{Transport: tr},
	}

	var err error
	c.client, err = api.NewClient(&clientCfg)

	if err != nil {
		return err
	}

	c.queryOpts = &api.QueryOptions{Datacenter: config.Consul.Datacenter}

	// authenticating
	if len(config.Consul.Token) != 0 {
		c.queryOpts.Token = config.Consul.Token
	}

	return nil
}

func (c *Consul) Fetch(variable model.VariableDef) (*model.Variable, error) {
	if variable.Source != model.CONSUL {
		return nil, fmt.Errorf("ConsulInput: Invalid variable type: %s for %s", variable.Type, variable.Name)
	}

	if c.client == nil {
		c.initClient()
	}

	consulValue, qm, err := c.client.KV().Get(variable.Value.(string), c.queryOpts)

	if err != nil {
		return nil, err
	}

	if consulValue == nil || consulValue.Value == nil {
		log.WithFields(log.Fields{
			"variable":   variable.Name,
			"path":       variable.Value,
			"query-meta": qm,
		}).Warning("ConsulInput: value not found")
		return nil, fmt.Errorf("ConsulInput: value for variable '%s' is missing: %s", variable.Name, variable.Value)
	}

	log.WithFields(log.Fields{
		"variable": variable.Name,
		"path":     variable.Value,
		"value":    string(consulValue.Value),
	}).Debug("Read variable from consul")

	ret := model.Variable{
		Name:  variable.Name,
		Type:  variable.Type,
		Value: string(consulValue.Value),
	}

	return &ret, nil
}

func init() {
	Register(model.CONSUL, &Consul{})
}
