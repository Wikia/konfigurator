package inputs

import (
	"fmt"

	"github.com/Wikia/konfigurator/model"
)

type Input interface {
	Fetch(variable model.VariableDef) (*model.Variable, error)
}

var registry map[model.InputType]Input

func Register(inputType model.InputType, input Input) error {
	_, has := registry[inputType]
	if has {
		return fmt.Errorf("Input already defined: %s", inputType)
	}

	if registry == nil {
		registry = map[model.InputType]Input{}
	}

	registry[inputType] = input

	return nil
}

func Get(source model.InputType) Input {
	return registry[source]
}

func GetRegisteredNames() []model.InputType {
	keys := make([]model.InputType, len(registry))
	i := 0
	for k := range registry {
		keys[i] = k
		i += 1
	}

	return keys
}

func Process(defs []model.VariableDef) ([]model.Variable, error) {
	ret := []model.Variable{}

	for _, definition := range defs {
		processor := Get(definition.Source)

		if processor == nil {
			return nil, fmt.Errorf("Could not find input processor (%s) for: %s", definition.Source, definition.Name)
		}

		variable, err := processor.Fetch(definition)

		if err != nil {
			return nil, err
		}

		ret = append(ret, *variable)
	}

	return ret, nil
}
