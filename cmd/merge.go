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

	"github.com/Wikia/konfigurator/config"
	"github.com/Wikia/konfigurator/helpers"
	"github.com/Wikia/konfigurator/inputs"
	"github.com/Wikia/konfigurator/model"
	"github.com/spf13/cobra"
	"k8s.io/client-go/pkg/api"
	"k8s.io/client-go/pkg/apis/extensions/v1beta1"
)

var (
	SecretName string
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

		// keeping old copy for diff
		oldDeployment, err := api.Scheme.Copy(deployment)

		err = model.UpdateDeploymentInPlace(deployment, variables, SecretName, ContainerName, Overwrite)

		if err != nil {
			return fmt.Errorf("Error updating deployment: %s", err)
		}

		if !NoConfirm {
			model.DiffDeploymets(oldDeployment.(*v1beta1.Deployment), deployment)

			confirm, err := helpers.AskConfirm(os.Stdout, os.Stdin, "Apply changes?")

			if err != nil {
				return err
			}

			if !confirm {
				return nil
			}
		}

		destFile, err := os.Create(DestinationFile)
		if err != nil {
			return err
		}

		defer destFile.Close()

		w := bufio.NewWriter(destFile)
		defer w.Flush()

		err = model.WriteDeployment(deployment, leftOver, w)

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
	mergeCmd.Flags().StringVarP(&SecretName, "secretName", "s", "", "Name of the secret to use in the deployment mappings (defaults to 'containerName')")
	mergeCmd.Flags().StringVarP(&DestinationFile, "destinationFile", "d", "", "Destination file where to write updated deployment configuration")
	mergeCmd.Flags().BoolVarP(&NoConfirm, "yes", "y", false, "Answer all questions 'yes' - no confirmations and interaction")
	mergeCmd.Flags().BoolVarP(&Overwrite, "overwrite", "w", false, "Should configuration definitions be completely replaced by the new one or just appended")
}
