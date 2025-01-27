package commands

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/cjaewon/dd/server"
	"github.com/spf13/cobra"
)

var (
	serverWatch bool
	serverBind  string
	serverPort  int

	serverCmd = &cobra.Command{
		Use:   "server <dirpath>",
		Short: "Start the embedded web server that is capable of rendering Markdown",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			dirpath := args[0]
			stat, err := os.Stat(dirpath)
			if err != nil {
				return err
			}

			if !stat.IsDir() {
				return fmt.Errorf("%s is not directory", dirpath)
			}

			// purpose of using filepath Abs and Base is for getting exact directory name.
			// if we just use stat.Name(), it will return "." ".." when dirpath is "." "..".
			absPath, err := filepath.Abs(dirpath)
			if err != nil {
				return err
			}

			dirname := filepath.Base(absPath)

			s := &server.Server{
				Bind:        serverBind,
				Port:        serverPort,
				Watch:       serverWatch,
				RootDirName: dirname,
				RootDirPath: absPath,
			}

			if err := s.Run(); err != nil {
				return err
			}

			return nil
		},
	}
)

func init() {
	serverCmd.Flags().BoolVarP(&serverWatch, "watch", "w", true, "watch filesystem for changes and do live realoding")
	serverCmd.Flags().StringVarP(&serverBind, "bind", "b", "localhost", "interface to which the server will bind (default \"localhost\")")
	serverCmd.Flags().IntVarP(&serverPort, "port", "p", 1234, "port for server listening")
}
