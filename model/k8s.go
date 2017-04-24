package model

import (
	"github.com/ghodss/yaml"

	"io/ioutil"

	"bytes"
	"encoding/json"
	"os"

	"fmt"

	"regexp"

	"strings"

	v1 "k8s.io/client-go/pkg/api/v1"
	v1beta1 "k8s.io/client-go/pkg/apis/extensions/v1beta1"
)

var (
	timeStampRegex   = regexp.MustCompile(`\s+creationTimestamp: null`)
	emptyStructRegex = regexp.MustCompile(`\s+(?:status|selector|strategy): {}`)
)

func splitYamlDocument(contents []byte) [][]byte {
	return bytes.Split(contents, []byte("---\n"))
}

func ReadSecrets(filePath string) (*v1.Secret, error) {
	contents, err := ioutil.ReadFile(filePath)

	if err != nil {
		return nil, err
	}

	secret := v1.Secret{}

	for _, document := range splitYamlDocument(contents) {
		err = yaml.Unmarshal(document, &secret)

		if err != nil {
			return nil, err
		}

		if secret.Kind == "Secret" {
			break
		}
	}

	if secret.Kind != "Secret" {
		return nil, fmt.Errorf("Could not unmarshall Secrets")
	}

	return &secret, nil
}

func WriteSecrets(secret *v1.Secret, filePath string) error {
	return writeK8sYaml(secret, filePath)
}

func ReadConfigMap(filePath string) (*v1.ConfigMap, error) {
	contents, err := ioutil.ReadFile(filePath)

	if err != nil {
		return nil, err
	}

	configMap := v1.ConfigMap{}

	for _, document := range splitYamlDocument(contents) {
		err = yaml.Unmarshal(document, &configMap)

		if err != nil {
			return nil, err
		}

		if configMap.Kind == "ConfigMap" {
			break
		}
	}

	if configMap.Kind != "ConfigMap" {
		return nil, fmt.Errorf("Could not unmarshall ConfigMap")
	}

	return &configMap, nil
}

func WriteConfigMap(configMap *v1.ConfigMap, filePath string) error {
	return writeK8sYaml(configMap, filePath)
}

func ReadDeployment(filePath string) (*v1beta1.Deployment, error) {
	contents, err := ioutil.ReadFile(filePath)

	if err != nil {
		return nil, err
	}

	deployment := v1beta1.Deployment{}

	for _, document := range splitYamlDocument(contents) {
		err = yaml.Unmarshal(document, &deployment)

		if err != nil {
			return nil, err
		}

		if deployment.Kind == "Deployment" {
			break
		}
	}

	if deployment.Kind != "Deployment" {
		return nil, fmt.Errorf("Could not unmarshall Deployment")
	}

	return &deployment, nil
}

func WriteDeployment(deployment *v1beta1.Deployment, filePath string) error {
	return writeK8sYaml(deployment, filePath)
}

func writeK8sYaml(data interface{}, filePath string) error {
	secretFile, err := os.Create(filePath)
	if err != nil {
		return err
	}

	defer secretFile.Close()

	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	output, err := yaml.JSONToYAML(jsonData)
	if err != nil {
		return err
	}

	output = timeStampRegex.ReplaceAll(output, []byte(""))
	output = emptyStructRegex.ReplaceAll(output, []byte(""))

	_, err = secretFile.Write(output)

	if err != nil {
		return err
	}

	return nil
}

func UpdateDeployment(deployment *v1beta1.Deployment, configMap *v1.ConfigMap, secret *v1.Secret, containerName string, variables []VariableDef, overwriteEnv bool) error {
	var dstContainer *v1.Container

	for idx, container := range deployment.Spec.Template.Spec.Containers {
		if container.Name == containerName {
			dstContainer = &deployment.Spec.Template.Spec.Containers[idx]
			break
		}
	}

	if dstContainer == nil {
		return fmt.Errorf("Could not find container '%s' in deployment", containerName)
	}

	if overwriteEnv {
		dstContainer.Env = []v1.EnvVar{}
	}

	for _, variable := range variables {
		var envVarSource *v1.EnvVarSource

		switch variable.Type {
		case CONFIGMAP:
			envVarSource = &v1.EnvVarSource{
				ConfigMapKeyRef: &v1.ConfigMapKeySelector{
					Key:                  strings.ToLower(variable.Name),
					LocalObjectReference: v1.LocalObjectReference{Name: configMap.Name},
				},
			}
		case SECRET:
			envVarSource = &v1.EnvVarSource{
				SecretKeyRef: &v1.SecretKeySelector{
					Key:                  strings.ToLower(variable.Name),
					LocalObjectReference: v1.LocalObjectReference{Name: secret.Name},
				},
			}
		}

		for _, envs := range dstContainer.Env {
			if envs.Name == variable.Name {
				envs.Value = ""
				envs.ValueFrom = envVarSource
				envVarSource = nil
				break
			}
		}

		if envVarSource != nil {
			dstContainer.Env = append(dstContainer.Env, v1.EnvVar{Name: strings.ToUpper(variable.Name), ValueFrom: envVarSource})
			envVarSource = nil
		}
	}

	return nil
}
