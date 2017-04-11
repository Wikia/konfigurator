package inputs

import "fmt"

type Input interface {
	Unmarshal(b []byte, t interface{}) error
}

var registry map[string]Input

func Register(name string, input Input) error {
	has, _ := registry[name]
	if has {
		return fmt.Errorf("Input already defined: %s", name)
	}

	if registry == nil {
		registry = map[string]Input{}
	}

	registry[name] = input

	return nil
}

func Get(name string) Input {
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
