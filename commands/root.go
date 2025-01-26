package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use: "dd",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Do you want help? Use -h, --help flags")
		},
	}
)

func init() {
	rootCmd.AddCommand(newCmd)
}

func Execute() error {
	return rootCmd.Execute()
}
