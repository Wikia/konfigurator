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

type LayeredConsul struct {
	client    *api.Client
	queryOpts *api.QueryOptions
}

func (c *LayeredConsul) initClient() error {
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

func (c *LayeredConsul) Fetch(variable model.VariableDef) (*model.Variable, error) {
	if variable.Source != model.LAYERED_CONSUL {
		return nil, fmt.Errorf("LayeredConsulInput: Invalid variable type: %s for %s", variable.Type, variable.Name)
	}

	if c.client == nil {
		c.initClient()
	}

	appName := variable.Context["appname"]
	env := variable.Context["environment"]
	key := variable.Value.(string)
	keyPaths := []string{
		fmt.Sprintf("config/%s/%s/%s", appName, env, key),
		fmt.Sprintf("config/%s/base/%s", appName, key),
		fmt.Sprintf("config/base/%s/%s", env, key),
	}

	for _, key := range keyPaths {
		consulValue, qm, err := c.client.KV().Get(key, c.queryOpts)

		if err != nil {
			return nil, err
		}

		if consulValue != nil && consulValue.Value != nil {
			log.WithFields(log.Fields{
				"variable": variable.Name,
				"path":     variable.Value,
				"value":    string(consulValue.Value),
			}).Debug("LayeredConsulInput: Read variable from consul")

			ret := model.Variable{
				Name:  variable.Name,
				Type:  variable.Type,
				Value: string(consulValue.Value),
			}

			return &ret, nil
		}

		log.WithFields(log.Fields{
			"variable":   variable.Name,
			"path":       variable.Value,
			"query-meta": qm,
		}).Debug("LayeredConsulInput: value not found")
	}

	return nil, fmt.Errorf("LayeredConsulInput: Could not find value in Consul")
}

func init() {
	Register(model.LAYERED_CONSUL, &LayeredConsul{})
}
