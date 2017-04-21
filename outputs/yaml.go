package outputs

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	"github.com/ghodss/yaml"

	"encoding/base64"
	"encoding/json"

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

	cfgFile, err := os.Create(filepath.Join(destinationPath, fmt.Sprintf("%s_configMap.yaml", name)))
	if err != nil {
		return err
	}

	defer cfgFile.Close()

	secretFile, err := os.Create(filepath.Join(destinationPath, fmt.Sprintf("%s_secrets.yaml", name)))
	if err != nil {
		return err
	}

	defer secretFile.Close()

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
			secretStr := fmt.Sprintf("%s", variable.Value)
			encodedText := make([]byte, base64.StdEncoding.EncodedLen(len(secretStr)))
			base64.StdEncoding.Encode(encodedText, []byte(secretStr))
			secrets.Data[variable.Name] = encodedText
			break

		case model.CONFIGMAP:
			cfgMap.Data[variable.Name] = fmt.Sprintf("%s", variable.Value)
			break
		}
	}

	jsonData, err := json.Marshal(&cfgMap)
	if err != nil {
		return err
	}

	output, err := yaml.JSONToYAML(jsonData)
	if err != nil {
		return err
	}

	output = bytes.Replace(output, []byte("  creationTimestamp: null\n"), []byte(""), 1)

	_, err = cfgFile.Write(output)

	if err != nil {
		return err
	}

	jsonData, err = json.Marshal(&secrets)
	if err != nil {
		return err
	}

	output, err = yaml.JSONToYAML(jsonData)
	if err != nil {
		return err
	}

	output = bytes.Replace(output, []byte("  creationTimestamp: null\n"), []byte(""), 1)

	_, err = secretFile.Write(output)

	if err != nil {
		return err
	}

	return nil
}

func init() {
	Register("k8s-yaml", &OutputK8SYaml{})
}
