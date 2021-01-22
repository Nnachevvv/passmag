package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/nnachevv/passmag/storage"
	"github.com/nnachevv/passmag/user"
	"github.com/spf13/cobra"
)

// NewListCmd creates a new listCmd
func NewListCmd(stdio terminal.Stdio) *cobra.Command {
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "Lists all password from your vault",
		Long:  `Ask for authorization and lists all password from your vault`,
		RunE: func(cmd *cobra.Command, args []string) error {
			u, err := user.EnterSession(stdio)
			if err != nil {
				return err
			}

			var s storage.Storage

			err = json.Unmarshal(u.VaultData, &s)
			if err != nil {
				return err
			}

			for n, p := range s.Passwords {
				fmt.Printf("%s : %s\n", n, string(p))
			}

			return nil
		},
	}
	return listCmd
}
