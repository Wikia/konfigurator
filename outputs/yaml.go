package outputs

import "gopkg.in/yaml.v2"

type OutputYaml struct{}

func (o *OutputYaml) Marshal(t interface{}) ([]byte, error) {
	return yaml.Marshal(t)
}

func init() {
	Register("yaml", &OutputYaml{})
}
