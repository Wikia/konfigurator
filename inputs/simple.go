package inputs

import (
	"fmt"

	"github.com/Wikia/konfigurator/config"
	"github.com/Wikia/konfigurator/model"
)

type Simple struct{}

func (i *Simple) Fetch(variable config.VariableDef) (model.Variable, error) {
	if variable.Source != config.SIMPLE {
		return nil, fmt.Errorf("Invalid variable type: %s for %s", variable.Type, variable.Name)
	}

	ret := model.Variable{
		Name:  variable.Name,
		Type:  variable.Type,
		Value: variable.Value,
	}

	return ret, nil
}

func init() {
	Register(config.SIMPLE, &Simple{})
}
