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

	"io"

	log "github.com/Sirupsen/logrus"
	"github.com/Wikia/konfigurator/helpers"
	v1 "k8s.io/client-go/pkg/api/v1"
	v1beta1 "k8s.io/client-go/pkg/apis/extensions/v1beta1"
)

var (
	timeStampRegex        = regexp.MustCompile(`\s+creationTimestamp: null`)
	emptyStructRegex      = regexp.MustCompile(`\s+(?:status|selector|strategy): {}`)
	yamlDocumentSeparator = []byte("---\n")
)

func splitYamlDocument(contents []byte) [][]byte {
	return bytes.Split(contents, yamlDocumentSeparator)
}

func ReadSecrets(filePath string) (*v1.Secret, [][]byte, error) {
	contents, err := ioutil.ReadFile(filePath)

	if err != nil {
		return nil, nil, err
	}

	var secret v1.Secret
	var document []byte
	idx := 0
	documents := splitYamlDocument(contents)

	for idx, document = range documents {
		secret = v1.Secret{}
		err = yaml.Unmarshal(document, &secret)

		if err != nil {
			log.WithError(err).Warn("Error parsing YAML document")
			continue
		}

		if secret.Kind == "Secret" {
			break
		}
	}

	if secret.Kind != "Secret" {
		return nil, nil, fmt.Errorf("Could not unmarshall Secrets")
	}

	return &secret, append(documents[0:idx], documents[idx+1:]...), nil
}

func WriteSecrets(secret *v1.Secret, leftOver [][]byte, writer io.Writer) error {
	return writeK8sYaml(secret, leftOver, writer)
}

func ReadConfigMap(filePath string) (*v1.ConfigMap, [][]byte, error) {
	contents, err := ioutil.ReadFile(filePath)

	if err != nil {
		return nil, nil, err
	}

	var configMap v1.ConfigMap
	var document []byte
	idx := 0
	documents := splitYamlDocument(contents)

	for idx, document = range documents {
		configMap = v1.ConfigMap{}
		err = yaml.Unmarshal(document, &configMap)

		if err != nil {
			log.WithError(err).Warn("Error parsing YAML document")
			continue
		}

		if configMap.Kind == "ConfigMap" {
			break
		}
	}

	if configMap.Kind != "ConfigMap" {
		return nil, nil, fmt.Errorf("Could not unmarshall ConfigMap")
	}

	return &configMap, append(documents[0:idx], documents[idx+1:]...), nil
}

func WriteConfigMap(configMap *v1.ConfigMap, leftOver [][]byte, writer io.Writer) error {
	return writeK8sYaml(configMap, leftOver, writer)
}

func ReadDeployment(filePath string) (*v1beta1.Deployment, [][]byte, error) {
	contents, err := ioutil.ReadFile(filePath)

	if err != nil {
		return nil, nil, err
	}

	deployment := v1beta1.Deployment{}
	idx := 0
	var document []byte
	documents := splitYamlDocument(contents)

	for idx, document = range splitYamlDocument(contents) {
		err = yaml.Unmarshal(document, &deployment)

		if err != nil {
			log.WithError(err).Warn("Error parsing YAML document")
			continue
		}

		if deployment.Kind == "Deployment" {
			break
		}
	}

	if deployment.Kind != "Deployment" {
		return nil, nil, fmt.Errorf("Could not unmarshall Deployment")
	}

	return &deployment, append(documents[0:idx], documents[idx+1:]...), nil
}

func WriteDeployment(deployment *v1beta1.Deployment, leftOver [][]byte, writer io.Writer) error {
	return writeK8sYaml(deployment, leftOver, writer)
}

func marshalK8sEntity(obj interface{}) ([]byte, error) {
	jsonData, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}

	output, err := yaml.JSONToYAML(jsonData)
	if err != nil {
		return nil, err
	}

	output = timeStampRegex.ReplaceAll(output, []byte(""))
	output = emptyStructRegex.ReplaceAll(output, []byte(""))

	return output, nil
}

func writeK8sYaml(obj interface{}, leftOver [][]byte, writer io.Writer) error {
	output, err := marshalK8sEntity(obj)

	if err != nil {
		return err
	}

	leftOver = append(leftOver, output)

	_, err = writer.Write(bytes.Join(leftOver, yamlDocumentSeparator))

	if err != nil {
		return err
	}

	return nil
}

func DiffDeploymets(deployment1 *v1beta1.Deployment, deployment2 *v1beta1.Deployment) error {
	deployYaml1, err := marshalK8sEntity(deployment1)

	if err != nil {
		return err
	}

	deployYaml2, err := marshalK8sEntity(deployment2)

	if err != nil {
		return err
	}

	helpers.RenderDiff(os.Stdout, string(deployYaml1), string(deployYaml2))

	return nil
}

func getDeploymentContainer(deployment *v1beta1.Deployment, containerName string) (*v1.Container, error) {
	for idx, container := range deployment.Spec.Template.Spec.Containers {
		if container.Name == containerName {
			return &deployment.Spec.Template.Spec.Containers[idx], nil
		}
	}

	return nil, fmt.Errorf("Could not find container '%s' in deployment", containerName)
}

func UpdateDeploymentInPlace(deployment *v1beta1.Deployment, variables []Variable, configMapName string, secretName string, containerName string, overwriteEnv bool) error {
	dstContainer, err := getDeploymentContainer(deployment, containerName)

	if err != nil {
		return err
	}

	previousEnv := dstContainer.Env

	if overwriteEnv {
		dstContainer.Env = []v1.EnvVar{}
	}

	for _, variable := range variables {
		var envVarSource *v1.EnvVarSource
		var envVarSimple *v1.EnvVar

		switch variable.Destination {
		case CONFIGMAP:
			switch variable.Type {
			case REFERENCED:
				envVarSource = &v1.EnvVarSource{
					FieldRef: &v1.ObjectFieldSelector{
						FieldPath: variable.Value.(string),
					},
				}
			case STANDARD:
				envVarSource = &v1.EnvVarSource{
					ConfigMapKeyRef: &v1.ConfigMapKeySelector{
						Key:                  strings.ToLower(variable.Name),
						LocalObjectReference: v1.LocalObjectReference{Name: configMapName},
					},
				}
			case INLINE:
				envVarSimple = &v1.EnvVar{
					Name:  strings.ToUpper(variable.Name),
					Value: variable.Value.(string),
				}
			}
		case SECRET:
			switch variable.Type {
			case REFERENCED:
				envVarSource = &v1.EnvVarSource{
					FieldRef: &v1.ObjectFieldSelector{
						FieldPath: variable.Value.(string),
					},
				}
			case STANDARD:
				envVarSource = &v1.EnvVarSource{
					SecretKeyRef: &v1.SecretKeySelector{
						Key:                  strings.ToLower(variable.Name),
						LocalObjectReference: v1.LocalObjectReference{Name: secretName},
					},
				}
			case INLINE:
				envVarSimple = &v1.EnvVar{
					Name:  strings.ToUpper(variable.Name),
					Value: variable.Value.(string),
				}
			}
		}

		for _, envVar := range dstContainer.Env {
			if envVar.Name == strings.ToUpper(variable.Name) {
				if envVarSource != nil {
					envVar.Value = ""
					envVar.ValueFrom = envVarSource
					envVarSource = nil
				} else if envVarSimple != nil {
					envVar.Value = envVarSimple.Value
					envVar.ValueFrom = nil
					envVarSimple = nil
				}
				break
			}
		}

		if envVarSource != nil {
			dstContainer.Env = append(dstContainer.Env, v1.EnvVar{Name: strings.ToUpper(variable.Name), ValueFrom: envVarSource})
			envVarSource = nil
		} else if envVarSimple != nil {
			dstContainer.Env = append(dstContainer.Env, *envVarSimple)
			envVarSimple = nil
		}
	}

	changed, removed, added := diffEnvs(previousEnv, dstContainer.Env)

	log.WithFields(log.Fields{"changed": changed, "removed": removed, "added": added}).Info("Deployment updated")

	return nil
}

func diffEnvs(src, dst []v1.EnvVar) (changed, removed, added []string) {
	diff := map[string]byte{}
	for _, srcEnv := range src {
		diff[srcEnv.Name] |= 1
		for _, dstEnv := range dst {
			diff[dstEnv.Name] |= 2
			if srcEnv.Name == dstEnv.Name {
				if srcEnv.ValueFrom != nil && dstEnv.ValueFrom != nil && srcEnv.ValueFrom.String() == dstEnv.ValueFrom.String() {
					diff[srcEnv.Name] |= 4
				}
				break
			} else {

			}
		}
	}

	for name, value := range diff {
		switch value {
		case 1:
			removed = append(removed, name)
		case 3:
			changed = append(changed, name)
		case 2:
			added = append(added, name)
		}
	}

	return
}

func UpdateDeployment(deployment *v1beta1.Deployment, configMap *v1.ConfigMap, secret *v1.Secret, containerName string, variables []VariableDef, overwriteEnv bool) error {
	dstContainer, err := getDeploymentContainer(deployment, containerName)

	if err != nil {
		return err
	}

	previousEnv := dstContainer.Env

	if overwriteEnv {
		dstContainer.Env = []v1.EnvVar{}
	}

	for _, variable := range variables {
		var envVarSource *v1.EnvVarSource
		var envVarSimple *v1.EnvVar

		switch variable.Destination {
		case CONFIGMAP:
			switch variable.Source {
			case REFERENCE:
				envVarSource = &v1.EnvVarSource{
					FieldRef: &v1.ObjectFieldSelector{
						FieldPath: variable.Value.(string),
					},
				}
			case SIMPLE:
				envVarSimple = &v1.EnvVar{
					Name:  strings.ToUpper(variable.Name),
					Value: variable.Value.(string),
				}
			default:
				envVarSource = &v1.EnvVarSource{
					ConfigMapKeyRef: &v1.ConfigMapKeySelector{
						Key:                  strings.ToLower(variable.Name),
						LocalObjectReference: v1.LocalObjectReference{Name: configMap.Name},
					},
				}
			}

		case SECRET:
			switch variable.Source {
			case REFERENCE:
				envVarSource = &v1.EnvVarSource{
					FieldRef: &v1.ObjectFieldSelector{
						FieldPath: variable.Value.(string),
					},
				}
			case SIMPLE:
				envVarSimple = &v1.EnvVar{
					Name:  strings.ToUpper(variable.Name),
					Value: variable.Value.(string),
				}
			default:
				envVarSource = &v1.EnvVarSource{
					SecretKeyRef: &v1.SecretKeySelector{
						Key:                  strings.ToLower(variable.Name),
						LocalObjectReference: v1.LocalObjectReference{Name: secret.Name},
					},
				}
			}
		}

		for _, envVar := range dstContainer.Env {
			if envVar.Name == strings.ToUpper(variable.Name) {
				envVar.Value = ""
				envVar.ValueFrom = envVarSource
				envVarSource = nil
				break
			}
		}

		if envVarSource != nil {
			dstContainer.Env = append(dstContainer.Env, v1.EnvVar{Name: strings.ToUpper(variable.Name), ValueFrom: envVarSource})
			envVarSource = nil
		} else if envVarSimple != nil {
			dstContainer.Env = append(dstContainer.Env, *envVarSimple)
			envVarSimple = nil
		}
	}

	changed, removed, added := diffEnvs(previousEnv, dstContainer.Env)

	log.WithFields(log.Fields{"changed": changed, "removed": removed, "added": added}).Info("Deployment updated")

	return nil
}
