package outputs

import (
	"fmt"
	"path/filepath"

	"github.com/Wikia/konfigurator/model"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/client-go/pkg/api/v1"
)

type OutputK8SYaml struct{}

func (o *OutputK8SYaml) Save(name string, destination string, vars []model.Variable) error {
	destinationPath, err := filepath.Abs(destination)

	if err != nil {
		return err
	}

	cfgMap := v1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ConfigMap",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: "dev",
		},
		Data: map[string]string{},
	}
	secrets := v1.Secret{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Secret",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: "dev",
		},
		Data: map[string][]byte{},
		Type: v1.SecretTypeOpaque,
	}

	for _, variable := range vars {
		switch variable.Type {
		case model.SECRET:
			secrets.Data[variable.Name] = []byte(variable.Value.(string))
			break

		case model.CONFIGMAP:
			cfgMap.Data[variable.Name] = variable.Value.(string)
			break
		}
	}

	err = model.WriteConfigMap(&cfgMap, filepath.Join(destinationPath, fmt.Sprintf("%s_configMap.yaml", name)))

	if err != nil {
		return err
	}

	err = model.WriteSecrets(&secrets, filepath.Join(destinationPath, fmt.Sprintf("%s_secrets.yaml", name)))

	if err != nil {
		return err
	}

	return nil
}

func init() {
	Register("k8s-yaml", &OutputK8SYaml{})
}
