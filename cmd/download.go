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

	log "github.com/Sirupsen/logrus"
	"github.com/Wikia/konfigurator/config"
)

var (
	OutputFmt   string
	Destination string
	ServiceName string
)

// downloadCmd represents the download command
var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Download configuration and stores it locally",
	Long:  `Fetches configuration for configured sources and stores it locally`,
	RunE: func(cmd *cobra.Command, args []string) error {
		out := outputs.Get(OutputFmt)

		if out == nil {
			return fmt.Errorf("Unknown output format: %s", OutputFmt)
		}

		if len(ServiceName) == 0 {
			return fmt.Errorf("Missing service name")
		}

		cfg := config.Get()
		variables, err := inputs.Process(cfg.Definitions)

		if err != nil {
			return fmt.Errorf("Error processing variables: %s", err)
		}

		err = out.Save(ServiceName, Destination, variables)

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
	workingDir, err := os.Getwd()

	if err != nil {
		log.WithError(err).Error("Error getting working directory")
		os.Exit(-6)
	}
	downloadCmd.Flags().StringVarP(&OutputFmt, "output", "o", "k8s-yaml", fmt.Sprintf("Output format (available formats: %v)", outputs.GetRegisteredNames()))
	downloadCmd.Flags().StringVarP(&Destination, "destination", "d", workingDir, "Where to store the output files")
	downloadCmd.Flags().StringVarP(&ServiceName, "serviceName", "s", "", "What is the service name which settings will be downloaded as")
}
