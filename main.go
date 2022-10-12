package main

import (
	_ "embed"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/szpnygo/gtc/cmd"
)

//go:embed version
var Version string

func init() {
	cmd.RootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "show version",
	Long:  `show the gtc version`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(Version)
	},
}

func main() {
	cmd.Execute()
}
