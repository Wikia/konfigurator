package outputs

import (
	"fmt"
	"os"
	"path/filepath"

	"strings"

	"github.com/Wikia/konfigurator/model"
)

type OutputEnvrc struct{}

func (o *OutputEnvrc) Save(name string, destination string, vars []model.Variable) error {
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
		_, err = cfgFile.WriteString(fmt.Sprintf("export %s=\"%s\"\n", strings.ToUpper(variable.Name), variable.Value))

		if err != nil {
			return err
		}
	}

	return nil
}

func init() {
	Register("envrc", &OutputEnvrc{})
}
