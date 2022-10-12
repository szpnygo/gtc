package cmd

import (
	"github.com/spf13/cobra"
	"github.com/szpnygo/gtc/server"
)

func init() {
	RootCmd.AddCommand(serverCmd)
}

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "gtc server",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		server.Server()
	},
}
