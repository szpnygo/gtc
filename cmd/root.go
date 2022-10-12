/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"strings"

	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"github.com/szpnygo/gtc/client"
)

var (
	api string
)

var RootCmd = &cobra.Command{
	Use:   "gtc",
	Short: "gtc",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if len(api) == 0 {
			pterm.Warning.Printfln("signaling address is empty")
			return
		}
		if !strings.HasPrefix(api, "ws://") && !strings.HasPrefix(api, "wss://") {
			pterm.Warning.Printfln("signaling address should start with ws:// or wss://")
			return
		}

		c := client.NewGTCClient(api)
		if err := c.Run(); err != nil {
			pterm.Error.Printfln("Failed to run client: %v", err)
		}
	},
}

// Execute ...
func Execute() {
	RootCmd.Flags().StringVarP(&api, "signaling", "s", "", "signaling server address")
	cobra.CheckErr(RootCmd.Execute())
}
