package outputs

import (
	"fmt"

	"strings"

	"regexp"

	"io"

	"github.com/Wikia/konfigurator/model"
)

type OutputEnvrc struct{}

var escapeRegex = regexp.MustCompile(`([$\\_\x96])`)

func (o *OutputEnvrc) Save(name string, namespace string, writer io.Writer, vars []model.Variable) error {
	for _, variable := range vars {
		if variable.Type == model.REFERENCED {
			continue
		}

		value := escapeRegex.ReplaceAllString(variable.Value.(string), "\\$1")

		_, err := fmt.Fprintf(writer, "export %s=\"%s\"\n", strings.ToUpper(variable.Name), value)

		if err != nil {
			return err
		}
	}

	return nil
}

func init() {
	Register("envrc", &OutputEnvrc{})
}
