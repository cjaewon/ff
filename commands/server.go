package commands

import "github.com/spf13/cobra"

var (
	serverWatch bool
	serverBind  string
	serverPort  int

	serverCmd = &cobra.Command{
		Use:   "server",
		Short: "Start the embedded web server that is capable of rendering Markdown",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}
)

func init() {
	serverCmd.Flags().BoolVarP(&serverWatch, "watch", "w", true, "watch filesystem for changes and do live realoding")
	serverCmd.Flags().StringVarP(&serverBind, "bind", "b", "127.0.0.1", "interface to which the server will bind (default \"127.0.0.1\")")
	serverCmd.Flags().IntVarP(&serverPort, "port", "p", 1234, "port for server listening")
}
