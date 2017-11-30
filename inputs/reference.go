package inputs

import (
	"fmt"

	"github.com/Wikia/konfigurator/model"
)

type Reference struct{}

func (i *Reference) Fetch(variable model.VariableDef) (*model.Variable, error) {
	if variable.Source != model.REFERENCE {
		return nil, fmt.Errorf("Invalid variable source: %s for %s", variable.Source, variable.Name)
	}

	var destination model.VariableDestination

	if variable.Destination != "" {
		destination = variable.Destination
	} else {
		destination = model.CONFIGMAP
	}

	ret := model.Variable{
		Name:        variable.Name,
		Type:        model.REFERENCED,
		Destination: destination,
		Value:       variable.Value,
	}

	return &ret, nil
}

func init() {
	Register(model.REFERENCE, &Reference{})
}
