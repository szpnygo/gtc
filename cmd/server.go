package cmd

import (
	"github.com/spf13/cobra"
	"github.com/szpnygo/gtc/server"
)

var (
	rooms string
)

func init() {
	RootCmd.AddCommand(serverCmd)
	serverCmd.Flags().StringVarP(&rooms, "rooms", "r", "", "room list split by comma")
}

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "gtc server",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		server.Server(rooms)
	},
}
