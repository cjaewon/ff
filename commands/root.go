package commands

import (
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use: "ff",
	}
)

func init() {
	rootCmd.AddCommand(newCmd)
	rootCmd.AddCommand(serverCmd)
}

func Execute() error {
	return rootCmd.Execute()
}
