package inputs

import (
	"fmt"

	"github.com/Wikia/konfigurator/config"
	"github.com/Wikia/konfigurator/model"
)

type Type string

const (
	SIMPLE  Type = "simple"
	POD          = "pod"
	VAULT        = "valut"
	FILE         = "file"
	PODYAML      = "pod_yaml"
)

type Input interface {
	Fetch(variable config.VariableDef) ([]model.Variable, error)
}

var registry map[Type]Input

func Register(inputType Type, input Input) error {
	has, _ := registry[inputType]
	if has {
		return fmt.Errorf("Input already defined: %s", inputType)
	}

	if registry == nil {
		registry = map[Type]Input{}
	}

	registry[inputType] = input

	return nil
}

func Get(source Type) Input {
	return registry[source]
}

func GetRegisteredNames() []Type {
	keys := make([]Type, len(registry))
	i := 0
	for k := range registry {
		keys[i] = k
		i += 1
	}

	return keys
}
