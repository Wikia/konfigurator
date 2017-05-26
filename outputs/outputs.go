package outputs

import (
	"fmt"

	"github.com/Wikia/konfigurator/model"
)

type Output interface {
	Save(name string, namespace string, destination string, vars []model.Variable) error
}

var registry map[string]Output

func Register(name string, output Output) error {
	_, has := registry[name]
	if has {
		return fmt.Errorf("Output already defined: %s", name)
	}

	if registry == nil {
		registry = map[string]Output{}
	}

	registry[name] = output

	return nil
}

func Get(name string) Output {
	return registry[name]
}

func GetRegisteredNames() []string {
	keys := make([]string, len(registry))
	i := 0
	for k := range registry {
		keys[i] = k
		i += 1
	}

	return keys
}
