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

	"github.com/Wikia/konfigurator/inputs"
	"github.com/Wikia/konfigurator/outputs"
	"github.com/spf13/cobra"

	"os"

	"github.com/Wikia/konfigurator/config"
	"github.com/spf13/viper"
)

var (
	OutputFmt string
)

// downloadCmd represents the download command
var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Downloads configuration and print out its contents",
	Long:  `Fetches configuration for configured sources and outputs it on stdout`,
	RunE: func(cmd *cobra.Command, args []string) error {
		out := outputs.Get(OutputFmt)

		if out == nil {
			return fmt.Errorf("Unknown output format: %s", OutputFmt)
		}

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

		err = out.Save(cfg.Application.Name, cfg.Application.Namespace, os.Stdout, variables)

		if err != nil {
			return fmt.Errorf("Error saving variables: %s", err)
		}

		return nil
	},
	SilenceErrors: true,
	SilenceUsage:  true,
}

func init() {
	RootCmd.AddCommand(downloadCmd)

	downloadCmd.Flags().StringVarP(&OutputFmt, "output", "o", "k8s-yaml", fmt.Sprintf("Output format (available formats: %v)", outputs.GetRegisteredNames()))

	downloadCmd.PersistentFlags().StringP("namespace", "n", "dev", "Kubernetes namespace for which files should be generated for")
	downloadCmd.PersistentFlags().String("name", "", "Name of the service to download variables for")

	viper.BindPFlag("application.namespace", downloadCmd.PersistentFlags().Lookup("namespace"))
	viper.BindPFlag("application.name", downloadCmd.PersistentFlags().Lookup("name"))
}
