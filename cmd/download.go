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

	"github.com/Wikia/konfigurator/outputs"
	"github.com/spf13/cobra"

	"os"

	"github.com/Wikia/konfigurator/model"
	log "github.com/sirupsen/logrus"
)

var (
	OutputFmt   string
	Destination string
)

// downloadCmd represents the download command
var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Download configuration and stores it locally",
	Long:  `Fetches configuration for configured sources and stores it locally`,
	Run: func(cmd *cobra.Command, args []string) {
		out := outputs.Get(OutputFmt)

		if out == nil {
			log.WithField("output", OutputFmt).Error("Unknown output format")
			return
		}

		yamlOut := outputs.Get("yaml")
		vars := []model.Variable{
			{
				Name:  "test1",
				Type:  model.CONFIGMAP,
				Value: "foo1",
			},
			{
				Name:  "secret1",
				Type:  model.SECRET,
				Value: 123,
			},
		}
		yamlOut.Save("helios", Destination, vars)
	},
}

func init() {
	RootCmd.AddCommand(downloadCmd)
	workingDir, err := os.Getwd()

	if err != nil {
		log.WithError(err).Error("Error getting working directory")
		os.Exit(-6)
	}
	downloadCmd.Flags().StringVarP(&OutputFmt, "output", "o", "yaml", fmt.Sprintf("Output format (available formats: %v)", outputs.GetRegisteredNames()))
	downloadCmd.Flags().StringVarP(&Destination, "destination", "d", workingDir, "Where to store the output files")
}
