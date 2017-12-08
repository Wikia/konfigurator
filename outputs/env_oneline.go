package outputs

import (
	"fmt"

	"strings"

	"io"

	"github.com/Wikia/konfigurator/model"
)

type OutputOneline struct{}

func (o *OutputOneline) Save(name string, namespace string, writer io.Writer, vars []model.Variable) error {
	for _, variable := range vars {
		if variable.Type == model.REFERENCED {
			continue
		}

		value := escapeRegex.ReplaceAllString(variable.Value.(string), "\\$1")

		_, err := fmt.Fprintf(writer, "%s=\"%s\" ", strings.ToUpper(variable.Name), value)

		if err != nil {
			return err
		}
	}

	return nil
}

func init() {
	Register("envoneline", &OutputOneline{})
}
