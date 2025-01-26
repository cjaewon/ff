package commands

import (
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"
)

var (
	newCmd = &cobra.Command{
		Use:  "new <filename>",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			filename := args[0]

			file, err := os.OpenFile(filename, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0644)
			if err == os.ErrExist {
				return err
			} else if err != nil {
				return err
			}

			defer file.Close()

			title := strings.TrimSuffix(filepath.Base(filename), filepath.Ext(filename))
			now := time.Now().Format(time.RFC3339)

			head := heredoc.Docf(`
                ---
                title: "%s"
                date: %s
                ---
            `, title, now)

			if _, err := file.WriteString(head); err != nil {
				return err
			}

			return nil
		},
	}
)
