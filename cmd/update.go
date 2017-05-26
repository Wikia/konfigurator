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

	"os"

	"bufio"

	"github.com/Wikia/konfigurator/config"
	"github.com/Wikia/konfigurator/helpers"
	"github.com/Wikia/konfigurator/model"
	"github.com/spf13/cobra"
	"k8s.io/client-go/pkg/api"
	"k8s.io/client-go/pkg/apis/extensions/v1beta1"
)

var (
	DeploymentFile  string
	Overwrite       bool
	ConfigFile      string
	SecretsFile     string
	ContainerName   string
	DestinationFile string
	NoConfirm       bool
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Creates updated configuration definition file",
	Long: `Updates configuration definition file according
to defined variables and saves it to a specified destination file`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(DeploymentFile) == 0 {
			return fmt.Errorf("Missing deployment file")
		}

		if len(ConfigFile) == 0 {
			return fmt.Errorf("Missing ConfigMap file")
		}

		if len(SecretsFile) == 0 {
			return fmt.Errorf("Missing secrets file")
		}

		if len(ContainerName) == 0 {
			return fmt.Errorf("Missing container name")
		}

		if len(DestinationFile) == 0 {
			return fmt.Errorf("Missing destination file")
		}

		secret, _, err := model.ReadSecrets(SecretsFile)

		if err != nil {
			return err
		}

		configMap, _, err := model.ReadConfigMap(ConfigFile)

		if err != nil {
			return err
		}

		deployment, leftOver, err := model.ReadDeployment(DeploymentFile)

		if err != nil {
			return err
		}

		// keeping old copy for diff
		oldDeployment, err := api.Scheme.Copy(deployment)

		if err != nil {
			return err
		}

		cfg := config.Get()
		err = model.UpdateDeployment(deployment, configMap, secret, ContainerName, cfg.Application.Definitions, Overwrite)

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
	RootCmd.AddCommand(updateCmd)

	updateCmd.Flags().StringVarP(&DeploymentFile, "deployment", "f", "", "Deployment file with configuration that should be updated")
	updateCmd.Flags().StringVarP(&ContainerName, "containerName", "t", "", "Name of the container to modify in deployment")
	updateCmd.Flags().StringVarP(&ConfigFile, "configMap", "m", "", "File where ConfigMap definitions are stored")
	updateCmd.Flags().StringVarP(&SecretsFile, "secrets", "s", "", "File where Secrets are stored")
	updateCmd.Flags().StringVarP(&DestinationFile, "destinationFile", "d", "", "Destination file where to write updated deployment configuration")
	updateCmd.Flags().BoolVarP(&NoConfirm, "yes", "y", false, "Answer all questions 'yes' - no confirmations and interaction")
	updateCmd.Flags().BoolVarP(&Overwrite, "overwrite", "w", false, "Should configuration definitions be completely replaced by the new one or just appended")
}
