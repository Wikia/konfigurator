package outputs

import (
	"fmt"

	"strings"

	"io"

	"github.com/Wikia/konfigurator/model"
)

type OutputEnvrc struct{}

func (o *OutputEnvrc) Save(name string, namespace string, writer io.Writer, vars []model.Variable) error {
	for _, variable := range vars {
		if variable.Type == model.REFERENCED {
			continue
		}

		value := fmt.Sprintf("%s", variable.Value)
		if strings.Contains(value, "(") {
			value = fmt.Sprintf("\"%s\"", value)
		}

		_, err := fmt.Fprintf(writer, "export %s=%s\n", strings.ToUpper(variable.Name), value)

		if err != nil {
			return err
		}
	}

	return nil
}

func init() {
	Register("envrc", &OutputEnvrc{})
}
