package outputs

import (
	"fmt"
	"os"
	"path/filepath"

	"strings"

	"regexp"

	"github.com/Wikia/konfigurator/model"
)

type OutputEnvrc struct{}

var escapeRegex = regexp.MustCompile(`([$\\_\x96])`)

func (o *OutputEnvrc) Save(name string, namespace string, destination string, vars []model.Variable) error {
	destinationPath, err := filepath.Abs(destination)

	if err != nil {
		return err
	}

	cfgFile, err := os.Create(filepath.Join(destinationPath, fmt.Sprintf("%s.envrc", name)))
	if err != nil {
		return err
	}

	defer cfgFile.Close()

	for _, variable := range vars {
		if variable.Type == model.REFERENCE {
			continue
		}

		value := escapeRegex.ReplaceAllString(variable.Value.(string), "\\$1")

		_, err = cfgFile.WriteString(fmt.Sprintf("export %s=\"%s\"\n", strings.ToUpper(variable.Name), value))

		if err != nil {
			return err
		}
	}

	return nil
}

func init() {
	Register("envrc", &OutputEnvrc{})
}
