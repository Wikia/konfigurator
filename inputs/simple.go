package inputs

import (
	"fmt"

	"github.com/Wikia/konfigurator/model"
)

type Simple struct{}

func (i *Simple) Fetch(variable model.VariableDef) (*model.Variable, error) {
	if variable.Source != model.SIMPLE {
		return nil, fmt.Errorf("Invalid variable source: %s for %s", variable.Source, variable.Name)
	}

	ret := model.Variable{
		Name:        variable.Name,
		Source:      model.SECRET,
		Destination: variable.Destination,
		Value:       variable.Value,
	}

	return &ret, nil
}

func init() {
	Register(model.SIMPLE, &Simple{})
}
