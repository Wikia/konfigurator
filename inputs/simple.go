package inputs

import (
	"fmt"

	"github.com/Wikia/konfigurator/model"
)

type Simple struct{}

func (i *Simple) Fetch(variable model.VariableDef) (model.Variable, error) {
	if variable.Source != SIMPLE {
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
	Register(SIMPLE, &Simple{})
}
