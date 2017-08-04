// Copyright © 2017 Wikia Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"

	"bufio"
	"os"

	"path"

	"path/filepath"

	"github.com/Wikia/konfigurator/config"
	"github.com/Wikia/konfigurator/inputs"
	"github.com/Wikia/konfigurator/model"
	"github.com/Wikia/konfigurator/outputs"
	"github.com/spf13/cobra"
)

var (
	SecretName     string
	ConfigMapName  string
	DestinationDir string
)

// mergeCmd represents the update command
var mergeCmd = &cobra.Command{
	Use:   "merge",
	Short: "Fetches configuration and saves it into deployment file",
	Long: `Downloads (if necessary) all configuration and applies it on given
	deployments file`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := config.Get()

		if len(cfg.Application.Name) == 0 {
			return fmt.Errorf("Missing service name")
		}

		if len(cfg.Application.Namespace) == 0 {
			return fmt.Errorf("Missing namespace")
		}

		if len(DestinationDir) == 0 {
			return fmt.Errorf("Missing destination dir")
		}

		DestinationDir, err := filepath.Abs(DestinationDir)

		if err != nil {
			return err
		}

		varDefinitions, err := config.ParseVariableDefinitions(cfg.Application.Definitions)
		if err != nil {
			return err
		}
		variables, err := inputs.Process(varDefinitions)

		if err != nil {
			return fmt.Errorf("Error processing variables: %s", err)
		}

		deployment, leftOver, err := model.ReadDeployment(DeploymentFile)

		if err != nil {
			return err
		}

		if len(SecretName) == 0 {
			SecretName = ContainerName
		}
		if len(ConfigMapName) == 0 {
			ConfigMapName = ContainerName
		}
		err = model.UpdateDeploymentInPlace(deployment, variables, ConfigMapName, SecretName, ContainerName, Overwrite)

		if err != nil {
			return fmt.Errorf("Error updating deployment: %s", err)
		}

		destDeploymentPath := path.Join(DestinationDir, fmt.Sprintf("%s_deployment.yaml", ContainerName))

		destDeploymentFile, err := os.Create(destDeploymentPath)
		if err != nil {
			return err
		}

		defer destDeploymentFile.Close()

		wDeploy := bufio.NewWriter(destDeploymentFile)
		defer wDeploy.Flush()

		err = model.WriteDeployment(deployment, leftOver, wDeploy)

		if err != nil {
			return err
		}

		destSecretsPath := path.Join(DestinationDir, fmt.Sprintf("%s_configs.yaml", ContainerName))

		destSecretsFile, err := os.Create(destSecretsPath)
		if err != nil {
			return err
		}

		defer destSecretsFile.Close()

		wSecrets := bufio.NewWriter(destSecretsFile)
		defer wSecrets.Flush()

		out := outputs.Get("k8s-yaml")
		if out == nil {
			return fmt.Errorf("Could not get output plugin for k8s")
		}

		err = out.Save(SecretName, cfg.Application.Namespace, wSecrets, variables)

		if err != nil {
			return err
		}

		return nil
	},
	SilenceErrors: true,
	SilenceUsage:  true,
}

func init() {
	RootCmd.AddCommand(mergeCmd)

	mergeCmd.Flags().StringVarP(&DeploymentFile, "deployment", "f", "", "Deployment file with configuration that should be updated")
	mergeCmd.Flags().StringVarP(&ContainerName, "containerName", "t", "", "Name of the container to modify in deployment")
	mergeCmd.PersistentFlags().StringP("namespace", "n", "dev", "Kubernetes namespace for which files should be generated for")
	mergeCmd.Flags().StringVarP(&ConfigMapName, "configMapName", "c", "", "Name of the ConfigMap to use in the deployment mappings (defaults to 'containerName')")
	mergeCmd.Flags().StringVarP(&SecretName, "secretName", "s", "", "Name of the secret to use in the deployment mappings (defaults to 'containerName')")
	mergeCmd.Flags().StringVarP(&DestinationDir, "destinationDir", "d", ".", "Destination where to write resulting filesn")
	mergeCmd.Flags().BoolVarP(&Overwrite, "overwrite", "w", true, "Should configuration definitions be completely replaced by the new one or just appended (defaults to true)")
}
