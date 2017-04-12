package inputs

import (
	"fmt"

	"github.com/Wikia/konfigurator/config"
)

type Input interface {
	Unmarshal(b []byte, t interface{}) error
}

var registry map[config.VariableSource]Input

func Register(source config.VariableSource, input Input) error {
	has, _ := registry[source]
	if has {
		return fmt.Errorf("Input already defined: %s", source)
	}

	if registry == nil {
		registry = map[config.VariableSource]Input{}
	}

	registry[source] = input

	return nil
}

func Get(source config.VariableSource) Input {
	return registry[source]
}

func GetRegisteredNames() []config.VariableSource {
	keys := make([]config.VariableSource, len(registry))
	i := 0
	for k := range registry {
		keys[i] = k
		i += 1
	}

	return keys
}
