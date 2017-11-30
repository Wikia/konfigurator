package outputs

import (
	"fmt"

	"io"

	"github.com/Wikia/konfigurator/model"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/client-go/pkg/api/v1"
)

type OutputK8SYaml struct{}

func (o *OutputK8SYaml) Save(name string, namespace string, writer io.Writer, vars []model.Variable) error {
	cfgMap := v1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ConfigMap",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
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
			Namespace: namespace,
		},
		Data: map[string][]byte{},
		Type: v1.SecretTypeOpaque,
	}

	for _, variable := range vars {
		if variable.Type == model.REFERENCED {
			continue
		}
		switch variable.Destination {
		case model.SECRET:
			secrets.Data[variable.Name] = []byte(variable.Value.(string))
		case model.CONFIGMAP:
			cfgMap.Data[variable.Name] = variable.Value.(string)
		}
	}

	if len(cfgMap.Data) > 0 {

		err := model.WriteConfigMap(&cfgMap, [][]byte{}, writer)

		if err != nil {
			return err
		}
	}

	if len(secrets.Data) > 0 {
		fmt.Fprintln(writer, "---")

		err := model.WriteSecrets(&secrets, [][]byte{}, writer)

		if err != nil {
			return err
		}
	}

	return nil
}

func init() {
	Register("k8s-yaml", &OutputK8SYaml{})
}
